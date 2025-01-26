package lexer

type TokenType uint

const (
	TokenPlus TokenType = iota
	TokenMinus
	TokenStar
	TokenSlash
	TokenEqual
	TokenEof
	TokenNumber
	TokenIdentifier

	TokenUnknown
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
