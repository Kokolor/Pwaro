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

func (parser *Parser) ParseNumber() *Node {
	token := parser.token
	parser.Expect(lexer.TokenNumber)
	return &Node{Token: token}
}

func (parser *Parser) ParseFactor() *Node {
	var left *Node

	switch parser.token.Type {
	case lexer.TokenNumber:
		left = parser.ParseNumber()
	case lexer.TokenIdentifier:
		left = &Node{Token: parser.token}
		parser.Advance()

		if parser.token.Type == lexer.TokenLeftParent {
			parser.Advance()
			parser.Expect(lexer.TokenRightParen)

			left = &Node{
				Token: lexer.Token{Type: lexer.TokenCall, Value: left.Token.Value},
				Left:  left,
			}
		}
	default:
		panic(fmt.Sprintln("Syntax Error, expected number or identifier, got", parser.token.Type.ToString()))
	}

	for parser.token.Type == lexer.TokenStar || parser.token.Type == lexer.TokenSlash {
		operator := parser.token
		parser.Advance()
		right := parser.ParseFactor()
		left = &Node{Left: left, Right: right, Token: operator}
	}

	return left
}

func (parser *Parser) ParseExpression() *Node {
	left := parser.ParseFactor()

	for parser.token.Type == lexer.TokenPlus || parser.token.Type == lexer.TokenMinus {
		operator := parser.token
		parser.Advance()
		right := parser.ParseFactor()
		left = &Node{Left: left, Right: right, Token: operator}
	}

	return left
}

func (parser *Parser) ParseBlock() []*Node {
	var statements []*Node

	for parser.token.Type != lexer.TokenRightParen {
		statement := parser.ParseStatement()
		statements = append(statements, statement)
	}

	return statements
}

func (parser *Parser) ParseStatement() *Node {
	switch parser.token.Type {
	case lexer.TokenPrint:
		parser.Advance()
		expression := parser.ParseExpression()
		parser.Expect(lexer.TokenSemi)

		return &Node{Token: lexer.Token{Type: lexer.TokenPrint, Value: "print"}, Left: expression}
	case lexer.TokenVar:
		parser.Advance()
		name := parser.token
		parser.Expect(lexer.TokenIdentifier)
		parser.Expect(lexer.TokenEqual)
		expression := parser.ParseExpression()
		parser.Expect(lexer.TokenSemi)

		return &Node{Token: lexer.Token{Type: lexer.TokenVar, Value: name.Value}, Left: expression}
	case lexer.TokenFn:
		parser.Advance()
		name := parser.token
		parser.Expect(lexer.TokenIdentifier)
		parser.Expect(lexer.TokenLeftParent)

		var body *Node
		statements := parser.ParseBlock()

		if len(statements) == 1 {
			body = statements[0]
		} else {
			body = &Node{
				Token:      lexer.Token{Type: lexer.TokenBlock, Value: "block"},
				Statements: statements,
			}
		}

		parser.Expect(lexer.TokenRightParen)
		parser.Expect(lexer.TokenSemi)

		return &Node{Token: lexer.Token{Type: lexer.TokenFn, Value: name.Value}, Left: body}
	case lexer.TokenPrototype:
		parser.Advance()
		name := parser.token
		parser.Expect(lexer.TokenIdentifier)
		parser.Expect(lexer.TokenSemi)

		return &Node{Token: lexer.Token{Type: lexer.TokenPrototype, Value: name.Value}}
	default:
		expr := parser.ParseExpression()
		parser.Expect(lexer.TokenSemi)
		return expr
	}
}

func (parser *Parser) Parse() []*Node {
	var statements []*Node

	for parser.token.Type != lexer.TokenEof {
		statement := parser.ParseStatement()
		statements = append(statements, statement)
	}

	return statements
}

func Print(node *Node) string {
	if node == nil {
		return ""
	}

	switch node.Token.Type {
	case lexer.TokenNumber:
		return node.Token.Value
	case lexer.TokenIdentifier:
		return node.Token.Value
	case lexer.TokenPrint:
		if node.Left == nil {
			return "print;"
		}

		expression := Print(node.Left)

		return fmt.Sprintf("print %s;", expression)
	case lexer.TokenVar:
		if node.Left == nil {
			return fmt.Sprintf("var %s;", node.Token.Value)
		}

		expression := Print(node.Left)

		return fmt.Sprintf("var %s = %s;", node.Token.Value, expression)
	case lexer.TokenFn:
		body := ""

		if node.Left != nil {
			body = Print(node.Left)
		}

		return fmt.Sprintf("fn %s: (\n\t%s\n);", node.Token.Value, body)
	case lexer.TokenPrototype:
		return fmt.Sprintf("prototype %s;", node.Token.Value)
	case lexer.TokenCall:
		return fmt.Sprintf("%s();", node.Left.Token.Value)
	default:
		if node.Left == nil || node.Right == nil {
			return node.Token.Value
		}

		left := Print(node.Left)
		right := Print(node.Right)

		return fmt.Sprintf("(%s %s %s)", left, node.Token.Value, right)
	}
}
