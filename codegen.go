package main

import (
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Compiler struct {
	programNode    ProgramNode
	module         *ir.Module
	currentContext Context
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
	return Compiler{programNode, ir.NewModule(), Context{}}
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
	c.module.NewGlobalDef("printIntegerFormat", NewCString("%d\n"))

	mainFunc := c.module.NewFunc("main", types.I32)
	b := mainFunc.NewBlock("")
	starterContext := newContext(b, c)
	c.currentContext = *starterContext
	// c.currentContext.NewRet(constant.NewInt(types.I32, 0))
	for _, inst := range c.programNode.instructions {
		c.currentContext.compileInst(inst)
	}
	c.currentContext.NewRet(constant.NewInt(types.I32, 0))
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
		c.vars[instNode.assignNode.identifier] = v
	case INST_IF:
		thenCtx := c.newContext(f.NewBlock(""), c.compiler)
		leaveBlock := f.NewBlock("")
		c.compiler.currentContext.NewBr(thenCtx.Block)
		c.compiler.currentContext = *c.newContext(leaveBlock, c.compiler)

		thenCtx.compileBlock(instNode.ifNode.ifBlockNode)
		if len(instNode.ifNode.elseBlockNode.instructions) > 0 {
			elseBlock := f.NewBlock("")
			c.newContext(elseBlock, c.compiler).compileBlock(instNode.ifNode.elseBlockNode)
			c.NewCondBr(c.compileRel(instNode.ifNode.relNode), thenCtx.Block, elseBlock)
			elseBlock.NewBr(leaveBlock)
		}

		thenCtx.NewBr(leaveBlock)

	case INST_PRINT:
		zero := constant.NewInt(types.I32, 0)
		printIntegerFormat := c.getPrintIntegerFormat()
		pointerToString := c.NewGetElementPtr(types.NewArray(3, types.I8), printIntegerFormat, zero, zero)
		c.NewCall(c.getPrintfFunc(),
			pointerToString,
			c.compileTerm(instNode.printNode.termNode))
	}
}

func (c *Context) getPrintfFunc() *ir.Func {
	for _, fun := range c.compiler.module.Funcs {
		if fun.GlobalName == "printf" {
			return fun
		}
	}
	panic("Couldn't find prinf function")
}

func (c Context) getPrintIntegerFormat() *ir.Global {
	for _, gl := range c.compiler.module.Globals {
		if gl.GlobalName == "printIntegerFormat" {
			return gl
		}
	}
	panic("Couldn't find printIntegerFormat global")
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
	case TERM_IDENT:
		return c.NewLoad(types.I32, c.lookupVariable(termNode.value))
	}

	panic("Unknown Term")
}

func (c *Context) lookupVariable(name string) value.Value {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookupVariable(name)
	} else {
		fmt.Printf("variable `%s`\n", name)
		panic("no such variable")
	}
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
