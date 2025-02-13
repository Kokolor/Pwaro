package lexer

type TokenType uint

const (
	TokenPlus TokenType = iota
	TokenMinus
	TokenStar
	TokenSlash
	TokenEqual
	TokenSemi
	TokenLeftParent
	TokenRightParen
	TokenEof
	TokenNumber
	TokenIdentifier

	TokenVar
	TokenFn
	TokenPrint
	TokenPrototype
	TokenCall


	TokenIntType

	TokenUnknown
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
