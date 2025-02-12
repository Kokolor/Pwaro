package main

import (
	"Pwaro/codegen"
	"Pwaro/lexer"
	"Pwaro/parser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"tinygo.org/x/go-llvm"
)

func run(filename string) {
	cmd := exec.Command("lli", filename)

	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("Runtime error: %v", err))
	}

	fmt.Println("Done:\n" + string(out))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <source_file>")
		os.Exit(1)
	}

	sourceFilename := os.Args[1]
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v", err))
	}

	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {

		}
	}(sourceFile)

	content, err := io.ReadAll(sourceFile)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %v", err))
	}

	lex := lexer.Lexer{}
	lex.InitLexer(strings.NewReader(string(content)))

	p := parser.Parser{}
	p.InitParser(&lex)
	trees := p.Parse()

	if err := llvm.InitializeNativeTarget(); err != nil {
		panic("Error initializing LLVM native target")
	}

	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()
	module := ctx.NewModule("module")

	mainFuncType := llvm.FunctionType(ctx.Int32Type(), nil, false)
	mainFunction := llvm.AddFunction(module, "main", mainFuncType)
	entryBlock := llvm.AddBasicBlock(mainFunction, "entry")
	builder := ctx.NewBuilder()
	builder.SetInsertPointAtEnd(entryBlock)

	codeGen := codegen.CodeGen{}
	codeGen.InitCodeGen(ctx, module, builder)
	for _, tree := range trees {
		_ = codeGen.GenerateIR(tree)
	}

	builder.CreateRet(llvm.ConstInt(ctx.Int32Type(), 0, false))

	fmt.Println("Generated IR:")
	fmt.Println(module.String())

	if err := llvm.VerifyModule(module, llvm.ReturnStatusAction); err != nil {
		panic(fmt.Sprintf("Invalid module: %v", err))
	}

	irFilename := fmt.Sprintf("%s.ll", os.Args[1])
	err = os.WriteFile(irFilename, []byte(module.String()), 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing IR to file: %v", err))
	}

	fmt.Printf("IR written to %s\n", irFilename)

	run(irFilename)

	err = os.Remove(irFilename)
	if err != nil {
		return 
	}
}
