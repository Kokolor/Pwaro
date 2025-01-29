package main

import (
	"Pwaro/codegen"
	"Pwaro/lexer"
	"Pwaro/parser"
	"fmt"
	"strings"
	"tinygo.org/x/go-llvm"
)

func main() {
	source := "var hello = 4 + 2;\nprint hello + 4;\nvar test = hello + 7;\nprint hello + test;"
	lex := lexer.Lexer{}
	lex.InitLexer(strings.NewReader(source))

	p := parser.Parser{}
	p.InitParser(&lex)

	trees := p.Parse()

	for _, tree := range trees {
		fmt.Println(parser.Print(tree))
	}

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

	var lastResult llvm.Value
	codeGen := codegen.CodeGen{}
	codeGen.InitCodeGen(ctx, module, builder)

	for _, tree := range trees {
		lastResult = codeGen.GenerateIR(tree)
	}

	builder.CreateRet(lastResult)

	fmt.Println(module.String())

	if err := llvm.VerifyModule(module, llvm.ReturnStatusAction); err != nil {
		panic(err)
	}
}
