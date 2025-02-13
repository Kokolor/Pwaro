package parser

import (
	"Pwaro/lexer"
	"fmt"
)

type Parser struct {
	lex   *lexer.Lexer
	token lexer.Token
}

func (parser *Parser) InitParser(lex *lexer.Lexer) {
	parser.lex = lex
	parser.token = parser.lex.Lex()
}

func (parser *Parser) Advance() {
	parser.token = parser.lex.Lex()
}

func (parser *Parser) Expect(tokenType lexer.TokenType) {
	if parser.token.Type != tokenType {
		panic(fmt.Sprintln("Syntax Error, expected", tokenType.ToString(), "got", parser.token.Type.ToString()))
	}
	parser.Advance()
}

func (parser *Parser) ParseNumber(typ string) Node {
	token := parser.token
	parser.Expect(lexer.TokenNumber)
	return &NodeExpr{
		BaseNode: BaseNode{token: token},
		Literal:  token.Value,
		Type:     typ,
	}
}

func (parser *Parser) ParseIdent() Node {
	token := parser.token
	parser.Expect(lexer.TokenIdentifier)
	return &NodeIdent{
		BaseNode: BaseNode{token: token},
		Name:     token.Value,
	}
}

func (parser *Parser) ParseFactor(typ string) Node {
	var left Node

	switch parser.token.Type {
	case lexer.TokenNumber:
		left = parser.ParseNumber(typ)
	case lexer.TokenIdentifier:
		left = parser.ParseIdent()
		if parser.token.Type == lexer.TokenLeftParent {
			parser.Advance()
			parser.Expect(lexer.TokenRightParen)
			left = &NodeFuncCall{
				BaseNode: BaseNode{
					token: parser.token,
				},
				FuncName: left.(*NodeIdent).Name,
			}
		}
	default:
		panic(fmt.Sprintln("Syntax Error, expected number or identifier, got", parser.token.Type.ToString()))
	}

	for parser.token.Type == lexer.TokenStar || parser.token.Type == lexer.TokenSlash {
		operator := parser.token
		parser.Advance()
		right := parser.ParseFactor(typ)
		left = &NodeExpr{
			BaseNode: BaseNode{
				token: operator,
				left:  left,
				right: right,
			},
			Operator: operator.Type,
		}
	}

	return left
}

func (parser *Parser) ParseExpr(typ string) Node {
	left := parser.ParseFactor(typ)
	for parser.token.Type == lexer.TokenPlus || parser.token.Type == lexer.TokenMinus {
		operator := parser.token
		parser.Advance()
		right := parser.ParseFactor(typ)
		left = &NodeExpr{
			BaseNode: BaseNode{
				token: operator,
				left:  left,
				right: right,
			},
			Operator: operator.Type,
			Type:     typ,
		}
	}
	return left
}

func (parser *Parser) ParseBlock() Node {
	var statements []Node
	for parser.token.Type != lexer.TokenRightParen {
		statement := parser.ParseStatement()
		statements = append(statements, statement)
	}
	return &NodeBlock{
		BaseNode: BaseNode{
			statements: statements,
		},
	}
}

func (parser *Parser) ParseStatement() Node {
	switch parser.token.Type {
	case lexer.TokenPrint:
		parser.Advance()
		expression := parser.ParseExpr("i32")
		parser.Expect(lexer.TokenSemi)
		return &NodePrint{
			BaseNode: BaseNode{
				token: lexer.Token{Type: lexer.TokenPrint, Value: "print"},
				left:  expression,
			},
			Expr: expression,
		}

	case lexer.TokenVar:
		parser.Advance()
		name := parser.token
		parser.Expect(lexer.TokenIdentifier)
		typ := parser.token
		parser.Expect(lexer.TokenIntType)
		parser.Expect(lexer.TokenEqual)
		expression := parser.ParseExpr(typ.Value)
		parser.Expect(lexer.TokenSemi)
		fmt.Println("Type of", name.Value, "is", typ.Value)
		return &NodeVar{
			BaseNode: BaseNode{
				token: lexer.Token{Type: lexer.TokenVar, Value: name.Value},
				left:  expression,
			},
			Name:  name.Value,
			Value: expression,
			Type:  typ.Value,
		}

	case lexer.TokenFn:
		parser.Advance()
		name := parser.token
		parser.Expect(lexer.TokenIdentifier)
		parser.Expect(lexer.TokenLeftParent)
		body := parser.ParseBlock()
		parser.Expect(lexer.TokenRightParen)
		parser.Expect(lexer.TokenSemi)
		return &NodeFunc{
			BaseNode: BaseNode{
				token: lexer.Token{Type: lexer.TokenFn, Value: name.Value},
				left:  body,
			},
			Name: name.Value,
			Body: body,
		}
	default:
		expr := parser.ParseExpr("i32")
		parser.Expect(lexer.TokenSemi)
		return expr
	}
}

func (parser *Parser) Parse() []Node {
	var statements []Node
	for parser.token.Type != lexer.TokenEof {
		statement := parser.ParseStatement()
		statements = append(statements, statement)
	}
	return statements
}