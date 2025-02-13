package codegen

import (
	"Pwaro/lexer"
	"Pwaro/parser"
	"fmt"
	"strconv"

	"tinygo.org/x/go-llvm"
)

type VarInfo struct {
	ptr llvm.Value
	typ llvm.Type
}

type CodeGen struct {
	ctx       llvm.Context
	module    llvm.Module
	builder   llvm.Builder
	variables map[string]VarInfo
	printf    llvm.Value
}

func (cg *CodeGen) InitCodeGen(ctx llvm.Context, module llvm.Module, builder llvm.Builder) {
	cg.ctx = ctx
	cg.module = module
	cg.builder = builder
	cg.variables = make(map[string]VarInfo)
	cg.printf = llvm.Value{}
}

func (cg *CodeGen) getLLVMType(typ string) llvm.Type {
	switch typ {
	case "i8":
		return cg.ctx.Int8Type()
	case "i16":
		return cg.ctx.Int16Type()
	case "i32":
		return cg.ctx.Int32Type()
	case "i64":
		return cg.ctx.Int64Type()
	default:
		panic(fmt.Sprintf("Unsupported type: %s", typ))
	}
}

func (cg *CodeGen) GenerateIR(node parser.Node) llvm.Value {
	switch n := node.(type) {
	case *parser.NodeIdent:
		return cg.genVar(n)
	case *parser.NodeExpr:
		return cg.genExpr(n)
	case *parser.NodeVar:
		return cg.genVarDecl(n)
	case *parser.NodePrint:
		return cg.genPrint(n)
	case *parser.NodeFunc:
		return cg.genFunc(n)
	case *parser.NodeFuncCall:
		return cg.genFuncCall(n)
	case *parser.NodeBlock:
		return cg.genBlock(n)
	default:
		panic(fmt.Sprintf("Unknown node: %T", node))
	}
}

func (cg *CodeGen) genExpr(n *parser.NodeExpr) llvm.Value {
	if n.Literal != "" {
		val, _ := strconv.ParseInt(n.Literal, 10, 32)
		return llvm.ConstInt(cg.getLLVMType(n.Type), uint64(val), false)
	}

	lhs := cg.GenerateIR(n.Left())
	rhs := cg.GenerateIR(n.Right())

	switch n.Operator {
	case lexer.TokenPlus:
		return cg.builder.CreateAdd(lhs, rhs, "addtmp")
	case lexer.TokenMinus:
		return cg.builder.CreateSub(lhs, rhs, "subtmp")
	case lexer.TokenStar:
		return cg.builder.CreateMul(lhs, rhs, "multmp")
	case lexer.TokenSlash:
		return cg.builder.CreateSDiv(lhs, rhs, "divtmp")
	default:
		panic("Unsupported operator: " + n.Operator.ToString())
	}
}

func (cg *CodeGen) genVar(n *parser.NodeIdent) llvm.Value {
	varInfo, exists := cg.variables[n.Name]
	if !exists {
		panic("Undefined variable: " + n.Name)
	}

	loaded := cg.builder.CreateLoad(varInfo.typ, varInfo.ptr, n.Name)
	width := varInfo.typ.IntTypeWidth()
	if width < 32 {
		loaded = cg.builder.CreateSExt(loaded, cg.ctx.Int32Type(), n.Name+"_sext")
	} /* else if width > 32 {
		loaded = cg.builder.CreateTrunc(loaded, cg.ctx.Int32Type(), n.Name+"_trunc")
	} */
	return loaded
}

func (cg *CodeGen) genVarDecl(n *parser.NodeVar) llvm.Value {
	varType := cg.getLLVMType(n.Type)
	alloca := cg.builder.CreateAlloca(varType, n.Name)
	value := cg.GenerateIR(n.Value)
	cg.builder.CreateStore(value, alloca)
	cg.variables[n.Name] = VarInfo{ptr: alloca, typ: varType}
	return alloca
}

func (cg *CodeGen) getPrintf() llvm.Value {
	if cg.printf.IsNil() {
		printfType := llvm.FunctionType(cg.ctx.Int32Type(), []llvm.Type{llvm.PointerType(cg.ctx.Int8Type(), 0)}, true)
		cg.printf = llvm.AddFunction(cg.module, "printf", printfType)
	}
	return cg.printf
}

func (cg *CodeGen) genPrint(n *parser.NodePrint) llvm.Value {
	printf := cg.getPrintf()
	formatStr := cg.builder.CreateGlobalStringPtr("%d\n", "printf_fmt")
	value := cg.GenerateIR(n.Expr)

	int8PtrTy := llvm.PointerType(cg.ctx.Int8Type(), 0)
	printfType := llvm.FunctionType(cg.ctx.Int32Type(), []llvm.Type{int8PtrTy}, true)
	cg.builder.CreateCall(printfType, printf, []llvm.Value{formatStr, value}, "printf_call")
	return llvm.Value{}
}

func (cg *CodeGen) genFunc(n *parser.NodeFunc) llvm.Value {
	oldInsert := cg.builder.GetInsertBlock()
	defer cg.builder.SetInsertPointAtEnd(oldInsert)

	funcType := llvm.FunctionType(cg.ctx.Int32Type(), nil, false)
	fn := llvm.AddFunction(cg.module, n.Name, funcType)
	entry := llvm.AddBasicBlock(fn, "entry")
	cg.builder.SetInsertPointAtEnd(entry)

	oldVars := cg.variables
	cg.variables = make(map[string]VarInfo)
	defer func() { cg.variables = oldVars }()

	stmts := n.Body.Statements()
	for i := 0; i < len(stmts)-1; i++ {
		cg.GenerateIR(stmts[i])
	}

	retVal := cg.GenerateIR(stmts[len(stmts)-1])
	cg.builder.CreateRet(retVal)

	return fn
}

func (cg *CodeGen) genFuncCall(n *parser.NodeFuncCall) llvm.Value {
	fn := cg.module.NamedFunction(n.FuncName)
	fnType := llvm.FunctionType(cg.ctx.Int32Type(), nil, false)

	if fn.IsNil() {
		fn = llvm.AddFunction(cg.module, n.FuncName, fnType)
	}

	return cg.builder.CreateCall(fnType, fn, []llvm.Value{}, "calltmp")
}

func (cg *CodeGen) genBlock(n *parser.NodeBlock) llvm.Value {
	for _, stmt := range n.Statements() {
		cg.GenerateIR(stmt)
	}
	return llvm.Value{}
}
