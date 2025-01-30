package parser

import "Pwaro/lexer"

type Node struct {
	Left  *Node
	Right *Node
	Token lexer.Token
	Statements []*Node
}
