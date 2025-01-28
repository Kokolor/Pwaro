package codegen

import (
	"tinygo.org/x/go-llvm"
	"strconv"
	"Pwaro/lexer"
	"Pwaro/parser"
	"fmt"
)

func GenerateIR(node *parser.Node, ctx llvm.Context, module llvm.Module, builder llvm.Builder) llvm.Value {
	switch node.Token.Type {
	case lexer.TokenNumber:
		value, err := strconv.Atoi(node.Token.Value)
		if err != nil {
			panic(fmt.Sprintf("Invalid number: %s", node.Token.Value))
		}

		return llvm.ConstInt(ctx.Int32Type(), uint64(value), false)
	case lexer.TokenPrint:
		exprValue := GenerateIR(node.Left, ctx, module, builder)

		printfType := llvm.FunctionType(ctx.Int32Type(), []llvm.Type{llvm.PointerType(ctx.Int8Type(), 0)}, true)
		printfFunc := llvm.AddFunction(module, "printf", printfType)

		formatStr := builder.CreateGlobalStringPtr("%d\n", "formatStr")

		builder.CreateCall(printfType, printfFunc, []llvm.Value{formatStr, exprValue}, "")

		return llvm.ConstInt(ctx.Int32Type(), 0, false)
	default:
		leftValue := GenerateIR(node.Left, ctx, module, builder)
		rightValue := GenerateIR(node.Right, ctx, module, builder)

		switch node.Token.Type {
		case lexer.TokenPlus:
			return builder.CreateAdd(leftValue, rightValue, "addtmp")
		case lexer.TokenMinus:
			return builder.CreateSub(leftValue, rightValue, "subtmp")
		case lexer.TokenStar:
			return builder.CreateMul(leftValue, rightValue, "multmp")
		case lexer.TokenSlash:
			return builder.CreateSDiv(leftValue, rightValue, "divtmp")
		default:
			panic(fmt.Sprintf("Unknown operator: %s", node.Token.Type.ToString()))
		}
	}
}
