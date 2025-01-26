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

	TokenVar
	TokenFn

	TokenUnknown
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
