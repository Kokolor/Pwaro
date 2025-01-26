package main

import (
	"Pwaro/lexer"
	"fmt"
	"strings"
)

func main() {
	source := "_14hello = 7 \n+ 2"
	lex := lexer.Lexer{}
	lex.InitLexer(strings.NewReader(source))

	for {
		token := lex.Lex()
		fmt.Println("[Line:", token.Line, "]", token.ToString(), token.Value)

		if token.Type == lexer.TokenEof {
			break
		}
	}
}
