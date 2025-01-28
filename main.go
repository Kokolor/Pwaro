package main

import (
	"Pwaro/lexer"
	"Pwaro/parser"
	"Pwaro/codegen"
	"tinygo.org/x/go-llvm"
	"fmt"
	"strings"
)

func main() {
	source := "print 4 + 2;"
	lex := lexer.Lexer{}
	lex.InitLexer(strings.NewReader(source))

	p := parser.Parser{}
	p.InitParser(&lex)

	tree := p.Parse()

	fmt.Println(parser.Print(tree))

	err := llvm.InitializeNativeTarget()
	if err != nil {
		return
	}

	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()
	module := ctx.NewModule("module")

	mainFunction := llvm.AddFunction(module, "main", llvm.FunctionType(ctx.Int32Type(), []llvm.Type{}, false))
	entryBlock := llvm.AddBasicBlock(mainFunction, "entry")
	builder := ctx.NewBuilder()
	builder.SetInsertPointAtEnd(entryBlock)

	result := codegen.GenerateIR(tree, ctx, module, builder)

	builder.CreateRet(result)

	fmt.Println(module.String())

	if err := llvm.VerifyModule(module, llvm.ReturnStatusAction); err != nil {
		panic(err)
	}
}
