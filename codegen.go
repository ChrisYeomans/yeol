package main

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Compiler struct {
	programNode ProgramNode
	module      *ir.Module
}

type Context struct {
	*ir.Block
	parent   *Context
	vars     map[string]value.Value
	compiler *Compiler
}

func NewCString(s string) *constant.CharArray {
	return constant.NewCharArrayFromString(s + "\x00")
}

func newCompiler(programNode ProgramNode) Compiler {
	return Compiler{programNode, ir.NewModule()}
}

func newContext(b *ir.Block, compiler *Compiler) *Context {
	return &Context{
		Block:    b,
		parent:   nil,
		vars:     make(map[string]value.Value),
		compiler: compiler,
	}
}

func (c *Context) newContext(b *ir.Block, compiler *Compiler) *Context {
	ctx := newContext(b, compiler)
	ctx.parent = c
	return ctx
}

func (c *Compiler) compileProgram() {
	printf := c.module.NewFunc("printf", types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	mainFunc := c.module.NewFunc("main", types.I32)
	b := mainFunc.NewBlock("")
	mainContext := newContext(b, c)
	mainContext.parent = mainContext
	for _, inst := range c.programNode.instructions {
		mainContext.compileInst(inst)
	}
	b.NewRet(constant.NewInt(types.I32, 0))
}

func (c *Context) compileBlock(blockNode BlockNode) {
	for _, inst := range blockNode.instructions {
		c.compileInst(inst)
	}
}

func (c *Context) compileInst(instNode InstNode) {
	f := c.Parent
	switch instNode.instType {
	case INST_ASSIGN:
		v := c.NewAlloca(types.I32)
		c.NewStore(c.compileExpr(instNode.assignNode.expr), v)
	case INST_IF:
		thenCtx := c.newContext(f.NewBlock("if.then"), c.compiler)
		thenCtx.compileBlock(instNode.ifNode.ifBlockNode)
		elseBlock := f.NewBlock("if.else")
		c.newContext(elseBlock, c.compiler).compileBlock(instNode.ifNode.elseBlockNode)
		c.NewCondBr(c.compileRel(instNode.ifNode.relNode), thenCtx.Block, elseBlock)
	case INST_PRINT:
		zero := constant.NewInt(types.I32, 0)
		printIntegerFormat := c.compiler.module.NewGlobalDef("tmp", NewCString("%d\n"))
		pointerToString := c.NewGetElementPtr(types.NewArray(3, types.I8), printIntegerFormat, zero, zero)
		c.NewCall(c.compiler.module.Funcs[0], pointerToString, c.compileTerm(instNode.printNode.termNode), constant.NewInt(types.I32, 0))
	}
}

func (c *Context) compileRel(relNode RelNode) value.Value {
	switch relNode.relType {
	case REL_LESS_THAN:
		l := c.compileTerm(relNode.termBinaryNode.lhs)
		r := c.compileTerm(relNode.termBinaryNode.rhs)
		return c.NewICmp(enum.IPredSLT, l, r)
	}
	panic("Unimplemented Relational")
}

func (c *Context) compileTerm(termNode TermNode) value.Value {
	switch termNode.termType {
	case TERM_INT:
		value, _ := strconv.ParseInt(termNode.value, 10, 32)
		return constant.NewInt(types.I32, value)
	}

	panic("Unknown Term")
}

func (c *Context) compileExpr(exprNode ExprNode) value.Value {
	switch exprNode.exprType {
	case EXPR_TERM:
		return c.compileTerm(exprNode.termNode)
	case EXPR_PLUS:
		l := c.compileTerm(exprNode.termBinaryNode.lhs)
		r := c.compileTerm(exprNode.termBinaryNode.rhs)
		return c.NewAdd(l, r)
	}

	panic("Unknown Expression")
}
