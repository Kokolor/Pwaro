package main

import (
	"Pwaro/lexer"
	"Pwaro/parser"
	"fmt"
	"strings"
)

func main() {
	source := "4 + 7 * 14 / 12"
	lex := lexer.Lexer{}
	lex.InitLexer(strings.NewReader(source))

	p := parser.Parser{}
	p.InitParser(&lex)

	tree := p.Parse()

	fmt.Println(parser.Print(tree))
}
