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
	left := parser.ParseNumber()

	for parser.token.Type == lexer.TokenStar || parser.token.Type == lexer.TokenSlash {
		operator := parser.token
		parser.Advance()
		right := parser.ParseNumber()
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

func (parser *Parser) ParseStatement() *Node {
	if parser.token.Type == lexer.TokenPrint {
		parser.Advance()
		expression := parser.ParseExpression()
		parser.Expect(lexer.TokenSemi)

		return &Node{Token: lexer.Token{Type: lexer.TokenPrint, Value: "print"}, Left: expression}
	}

	panic(fmt.Sprintln("Syntax Error, expected 'print' statement, got", parser.token.Type.ToString()))
}

func (parser *Parser) Parse() *Node {
	return parser.ParseStatement()
}

func Print(node *Node) string {
	if node.Token.Type == lexer.TokenNumber {
		return node.Token.Value
	}

	if node.Token.Type == lexer.TokenPrint {
		expression := Print(node.Left)
		return fmt.Sprintf("print %s;", expression)
	}

	left := Print(node.Left)
	right := Print(node.Right)
	return fmt.Sprintf("(%s %s %s)", left, node.Token.Value, right)
}
