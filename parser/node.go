package parser

import "Pwaro/lexer"

type Node interface {
	Token() lexer.Token
	Left() Node
	Right() Node
	Statements() []Node
}

type BaseNode struct {
	token      lexer.Token
	left       Node
	right      Node
	statements []Node
}

func (n *BaseNode) Token() lexer.Token {
	return n.token
}

func (n *BaseNode) Left() Node {
	return n.left
}

func (n *BaseNode) Right() Node {
	return n.right
}

func (n *BaseNode) Statements() []Node {
	return n.statements
}

type NodeVar struct {
	BaseNode
	Name  string
	Type  string
	Value Node
}

type NodeFunc struct {
	BaseNode
	Name string
	Body Node
}

type NodeFuncCall struct {
	BaseNode
	FuncName string
}

type NodeBlock struct {
	BaseNode
}

type NodeExpr struct {
	BaseNode
	Operator lexer.TokenType
	Literal  string
	Type     string
}

type NodeIdent struct {
	BaseNode
	Name string
}

type NodePrint struct {
	BaseNode
	Expr Node
}