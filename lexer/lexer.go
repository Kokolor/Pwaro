package lexer

import (
	"io"
	"text/scanner"
)

type Lexer struct {
	Scanner scanner.Scanner
}

func (token TokenType) ToString() string {
	switch token {
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
	case TokenSemi:
		return "TokenSemi"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenVar:
		return "TokenVar"
	case TokenFn:
		return "TokenFn"
	case TokenPrint:
		return "TokenPrint"
	case TokenPrototype:
		return "TokenPrototype"
	case TokenCall:
		return "TokenCall"
	case TokenIntType:
		return "TokenTypeInt"
	case TokenEof:
		return "TokenEof"
	default:
		return "TokenUnknown"
	}
}

func (lexer *Lexer) InitLexer(source io.Reader) {
	lexer.Scanner.Init(source)
}

func (lexer *Lexer) Lex() Token {
	char := lexer.Scanner.Scan()
	line := lexer.Scanner.Pos().Line

	switch char {
	case '+':
		return Token{Type: TokenPlus, Value: "+", Line: line}
	case '-':
		return Token{Type: TokenMinus, Value: "-", Line: line}
	case '*':
		return Token{Type: TokenStar, Value: "*", Line: line}
	case '/':
		return Token{Type: TokenSlash, Value: "/", Line: line}
	case '=':
		return Token{Type: TokenEqual, Value: "=", Line: line}
	case ';':
		return Token{Type: TokenSemi, Value: ";", Line: line}
	case '(':
		return Token{Type: TokenLeftParent, Value: "(", Line: line}
	case ')':
		return Token{Type: TokenRightParen, Value: ")", Line: line}
	case scanner.Int:
		tokenText := lexer.Scanner.TokenText()
		return Token{Type: TokenNumber, Value: tokenText, Line: line}
	case scanner.Ident:
		tokenText := lexer.Scanner.TokenText()
		switch tokenText {
		case "var":
			return Token{Type: TokenVar, Value: "var", Line: line}
		case "fn":
			return Token{Type: TokenFn, Value: "fn", Line: line}
		case "print":
			return Token{Type: TokenPrint, Value: "print", Line: line}
		case "prototype":
			return Token{Type: TokenPrototype, Value: "prototype", Line: line}
		case "i8", "i16", "i32", "i64":
			return Token{Type: TokenIntType, Value: tokenText, Line: line}
		default:
			return Token{Type: TokenIdentifier, Value: tokenText, Line: line}
		}
	case scanner.EOF:
		return Token{Type: TokenEof, Line: line}
	default:
		tokenText := lexer.Scanner.TokenText()
		return Token{Type: TokenUnknown, Value: tokenText, Line: line}
	}
}
