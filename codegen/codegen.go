package codegen

import (
    "Pwaro/lexer"
    "Pwaro/parser"
    "fmt"
    "strconv"
    "tinygo.org/x/go-llvm"
)

type CodeGen struct {
    ctx        llvm.Context
    module     llvm.Module
    builder    llvm.Builder
    printf     llvm.Value
    formatStr  llvm.Value
    variables  map[string]llvm.Value
    printfType llvm.Type
}

func (codegen *CodeGen) InitCodeGen(ctx llvm.Context, module llvm.Module, builder llvm.Builder) {
    codegen.ctx = ctx
    codegen.module = module
    codegen.builder = builder
    codegen.variables = make(map[string]llvm.Value)

    int8PtrTy := llvm.PointerType(codegen.ctx.Int8Type(), 0)
    codegen.printfType = llvm.FunctionType(codegen.ctx.Int32Type(), []llvm.Type{int8PtrTy}, true)

    codegen.printf = llvm.AddFunction(codegen.module, "printf", codegen.printfType)
    codegen.formatStr = codegen.builder.CreateGlobalStringPtr("%d\n", "formatStr")
}

func (codegen *CodeGen) generateNumber(value string) llvm.Value {
    num, err := strconv.Atoi(value)
    
    if err != nil {
        panic(fmt.Sprintf("Invalid number: %s", value))
    }

    return llvm.ConstInt(codegen.ctx.Int32Type(), uint64(num), false)
}

func (codegen *CodeGen) generateVariable(name string) llvm.Value {
    if alloca, exists := codegen.variables[name]; exists {
        return codegen.builder.CreateLoad(codegen.ctx.Int32Type(), alloca, "load_"+name)
    }

    panic(fmt.Sprintf("Undefined variable: %s", name))
}

func (codegen *CodeGen) generateBinaryOp(left, right llvm.Value, op lexer.TokenType) llvm.Value {
    switch op {
    case lexer.TokenPlus:
        return codegen.builder.CreateAdd(left, right, "add")
    case lexer.TokenMinus:
        return codegen.builder.CreateSub(left, right, "sub")
    case lexer.TokenStar:
        return codegen.builder.CreateMul(left, right, "mul")
    case lexer.TokenSlash:
        return codegen.builder.CreateSDiv(left, right, "div")
    default:
        panic(fmt.Sprintf("Unknown operator: %s", op.ToString()))
    }
}

func (codegen *CodeGen) generatePrint(exprValue llvm.Value) llvm.Value {
    args := []llvm.Value{codegen.formatStr, exprValue}
    codegen.builder.CreateCall(codegen.printfType, codegen.printf, args, "printfCall")

    return llvm.ConstInt(codegen.ctx.Int32Type(), 0, false)
}

func (codegen *CodeGen) generateVarDecl(name string, value llvm.Value) llvm.Value {
    alloca := codegen.builder.CreateAlloca(codegen.ctx.Int32Type(), name)
    codegen.builder.CreateStore(value, alloca)
    codegen.variables[name] = alloca

    return llvm.ConstInt(codegen.ctx.Int32Type(), 0, false)
}

func (codegen *CodeGen) GenerateIR(node *parser.Node) llvm.Value {
    if node == nil {
        panic("Null node encountered")
    }

    switch node.Token.Type {
    case lexer.TokenNumber:
        return codegen.generateNumber(node.Token.Value)
    case lexer.TokenPrint:
        if node.Left == nil {
            panic("Print statement requires an expression")
        }

        exprValue := codegen.GenerateIR(node.Left)

        return codegen.generatePrint(exprValue)
    case lexer.TokenVar:
        if node.Left == nil {
            panic("Variable declaration requires an initializer")
        }

        varValue := codegen.GenerateIR(node.Left)

        return codegen.generateVarDecl(node.Token.Value, varValue)
    case lexer.TokenIdentifier:
        return codegen.generateVariable(node.Token.Value)
    default:
        if node.Left == nil || node.Right == nil {
            panic("Binary operation requires two operands")
        }

        leftValue := codegen.GenerateIR(node.Left)
        rightValue := codegen.GenerateIR(node.Right)

        return codegen.generateBinaryOp(leftValue, rightValue, node.Token.Type)
    }
}
