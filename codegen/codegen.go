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
    functions  map[string]llvm.Value
    printfType llvm.Type
}

func (codegen *CodeGen) InitCodeGen(ctx llvm.Context, module llvm.Module, builder llvm.Builder) {
    codegen.ctx = ctx
    codegen.module = module
    codegen.builder = builder
    codegen.variables = make(map[string]llvm.Value)
    codegen.functions = make(map[string]llvm.Value)

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

func (codegen *CodeGen) generateBlock(node *parser.Node) llvm.Value {
    var lastValue llvm.Value

    for _, statement := range node.Statements {
        lastValue = codegen.GenerateIR(statement)
    }

    if lastValue.IsNil() {
        return llvm.ConstInt(codegen.ctx.Int32Type(), 0, false)
    }
    return lastValue
}

func (codegen *CodeGen) generateFunctionDecl(node *parser.Node) llvm.Value {
    name := node.Token.Value

    funcType := llvm.FunctionType(codegen.ctx.Int32Type(), []llvm.Type{}, false)
    function := llvm.AddFunction(codegen.module, name, funcType)

    entryBlock := llvm.AddBasicBlock(function, "entry")
    oldInsertPoint := codegen.builder.GetInsertBlock()

    codegen.builder.SetInsertPointAtEnd(entryBlock)

    oldVariables := codegen.variables
    codegen.variables = make(map[string]llvm.Value)

    var returnValue llvm.Value
    if node.Left != nil {
        returnValue = codegen.GenerateIR(node.Left)
    } else {
        returnValue = llvm.ConstInt(codegen.ctx.Int32Type(), 0, false)
    }

    codegen.builder.CreateRet(returnValue)

    codegen.variables = oldVariables

    if oldInsertPoint.IsNil() {
        codegen.builder.ClearInsertionPoint()
    } else {
        codegen.builder.SetInsertPointAtEnd(oldInsertPoint)
    }

    codegen.functions[name] = function
    return function
}

func (codegen *CodeGen) generateFunctionCall(node *parser.Node) llvm.Value {
    functionName := node.Left.Token.Value
    function, exists := codegen.functions[functionName]

    if !exists {
        panic(fmt.Sprintf("Function '%s' not declared", functionName))
    }

    if function.IsNil() {
        panic(fmt.Sprintf("Function '%s' is nil", functionName))
    }

    functionType := llvm.FunctionType(codegen.ctx.Int32Type(), []llvm.Type{}, false)

    return codegen.builder.CreateCall(functionType, function, []llvm.Value{}, "calltmp")
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
    case lexer.TokenFn:
        return codegen.generateFunctionDecl(node)
    case lexer.TokenIdentifier:
        return codegen.generateVariable(node.Token.Value)
    case lexer.TokenCall:
        return codegen.generateFunctionCall(node)
    case lexer.TokenBlock:
        return codegen.generateBlock(node)
    default:
        if node.Left == nil || node.Right == nil {
            panic("Binary operation requires two operands")
        }

        leftValue := codegen.GenerateIR(node.Left)
        rightValue := codegen.GenerateIR(node.Right)

        return codegen.generateBinaryOp(leftValue, rightValue, node.Token.Type)
    }
}

func (codegen *CodeGen) DisplayVariables() {
    for name := range codegen.variables {
        fmt.Printf("Variable: %s\n", name)
    }

    for name := range codegen.functions {
        fmt.Printf("Function: %s\n", name)
    }
}
