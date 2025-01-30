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
    source := "fn hello(\n\tvar test = 4;\n\tprint test;\n);\nhello();\nvar neuil = 96;\nprint neuil;"
    fmt.Println(source)
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

    var _ llvm.Value
    codeGen := codegen.CodeGen{}
    codeGen.InitCodeGen(ctx, module, builder)

    for _, tree := range trees {
        _ = codeGen.GenerateIR(tree)
    }

    builder.CreateRet(llvm.ConstInt(ctx.Int32Type(), 0, false))

    fmt.Println(module.String())

    codeGen.DisplayVariables()

    if err := llvm.VerifyModule(module, llvm.ReturnStatusAction); err != nil {
        panic(err)
    }
}
