package lexer

type TokenType uint

const (
	TokenPlus TokenType = iota
	TokenMinus
	TokenStar
	TokenSlash
	TokenEqual
	TokenSemi
	TokenEof
	TokenNumber
	TokenIdentifier

	TokenVar
	TokenFn
	TokenPrint

	TokenUnknown
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
