package lexer

import (
	"io"
	"text/scanner"
)

type Lexer struct {
	scanner scanner.Scanner
}

func (token *Token) ToString() string {
	switch token.Type {
	case TokenPlus:
		return "TokenPlus"
	case TokenMinus:
		return "TokenMinus"
	case TokenStar:
		return "TokenStar"
	case TokenSlash:
		return "TokenSlash"
	case TokenNumber:
		return "TokenNumber"
	case TokenEqual:
		return "TokenEqual"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenEof:
		return "TokenEof"
	default:
		return "TokenUnknown"
	}
}

func (lexer *Lexer) InitLexer(source io.Reader) {
	lexer.scanner.Init(source)
}

func (lexer *Lexer) Lex() Token {
	char := lexer.scanner.Scan()
	line := lexer.scanner.Pos().Line

	switch char {
	case '+':
		return Token{Type: TokenPlus, Line: line}
	case '-':
		return Token{Type: TokenMinus, Line: line}
	case '*':
		return Token{Type: TokenStar, Line: line}
	case '/':
		return Token{Type: TokenSlash, Line: line}
	case '=':
		return Token{Type: TokenEqual, Line: line}
	case scanner.Int:
		tokenText := lexer.scanner.TokenText()
		return Token{Type: TokenNumber, Value: tokenText, Line: line}
	case scanner.Ident:
		tokenText := lexer.scanner.TokenText()
		return Token{Type: TokenIdentifier, Value: tokenText, Line: line}
	case scanner.EOF:
		return Token{Type: TokenEof, Line: line}
	default:
		tokenText := lexer.scanner.TokenText()
		return Token{Type: TokenUnknown, Value: tokenText, Line: line}
	}
}
