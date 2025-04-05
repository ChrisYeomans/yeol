package main

import (
	"fmt"
	"strings"
)

type Assembler struct {
	programNode ProgramNode
	variables   []string
	ifCount     int
	fileSb      strings.Builder
}

func (a *Assembler) instrDeclareVariables(instNode InstNode) {
	switch instNode.instType {
	case INST_ASSIGN:
		for _, variable := range a.variables {
			if instNode.assignNode.identifier == variable {
				return
			}
		}
		a.variables = append(a.variables, instNode.assignNode.identifier)
	case INST_IF:
		a.relDeclareVariables(instNode.ifNode.relNode)
		a.blockDeclareVariables(instNode.ifNode.ifBlockNode)
	case INST_PRINT:
		a.termDeclareVariables(instNode.printNode.termNode)
	}
}

func (a *Assembler) blockDeclareVariables(blockNode BlockNode) {
	for _, inst := range blockNode.instructions {
		a.instrDeclareVariables(inst)
	}
}

func (a *Assembler) exprDeclareVariables(exprNode ExprNode) {
	switch exprNode.exprType {
	case EXPR_TERM:
		a.termDeclareVariables(exprNode.termNode)
	case EXPR_PLUS:
		a.termDeclareVariables(exprNode.termBinaryNode.lhs)
		a.termDeclareVariables(exprNode.termBinaryNode.rhs)
	}
}

func (a *Assembler) relDeclareVariables(relNode RelNode) {
	switch relNode.relType {
	case REL_LESS_THAN:
		a.termDeclareVariables(relNode.termBinaryNode.lhs)
		a.termDeclareVariables(relNode.termBinaryNode.rhs)
	}
}

func (a *Assembler) termDeclareVariables(termNode TermNode) {
	switch termNode.termType {
	case TERM_INPUT:
		break
	case TERM_INT:
		break
	case TERM_IDENT:
		// TODO: Add check for existing variable and panic
		break
	}
}

func (a *Assembler) assembleProgram() {
	for _, instNode := range a.programNode.instructions {
		a.instrDeclareVariables(instNode)
	}

	a.fileSb.WriteString("LINE_MAX equ 1024\n")
	a.fileSb.WriteString("%include \"string.inc\"\n")
	a.fileSb.WriteString("%include \"util.inc\"\n")
	a.fileSb.WriteString("SECTION .text\n")
	a.fileSb.WriteString("global _start\n")
	a.fileSb.WriteString("_start:\n")

	a.fileSb.WriteString("mov rbp, rsp\n")
	a.fileSb.WriteString(fmt.Sprintf("    sub rsp, %d\n", len(a.variables)*8))

	for _, instNode := range a.programNode.instructions {
		a.assembleInst(instNode)
	}

	a.fileSb.WriteString(fmt.Sprintf("    add rsp, %d\n", len(a.variables)*8))

	a.fileSb.WriteString("    exit_program 0\n")

	a.fileSb.WriteString("SECTION .bss\n")
	a.fileSb.WriteString("    line: resb LINE_MAX\n")
}

func (a *Assembler) assembleInst(instNode InstNode) {
	switch instNode.instType {
	case INST_ASSIGN:
		a.assembleExpr(instNode.assignNode.expr)
		index := a.findVariableIndex(instNode.assignNode.identifier)
		a.fileSb.WriteString(fmt.Sprintf("    mov qword [rbp - %d], rax\n", index*8+8))
	case INST_IF:
		a.assembleRel(instNode.ifNode.relNode)
		label := a.ifCount
		a.ifCount++
		if len(instNode.ifNode.elseBlockNode.instructions) > 0 {
			a.fileSb.WriteString("    test rax, rax\n")
			a.fileSb.WriteString(fmt.Sprintf("    jz .else%d\n", label))
			a.assembleBlock(instNode.ifNode.ifBlockNode)
			a.fileSb.WriteString(fmt.Sprintf("    jmp .endif%d\n", label))
			a.fileSb.WriteString(fmt.Sprintf(".else%d:\n", label))
			a.assembleBlock(instNode.ifNode.elseBlockNode)
			a.fileSb.WriteString(fmt.Sprintf(".endif%d:\n", label))
		} else {
			a.fileSb.WriteString("    test rax, rax\n")
			a.fileSb.WriteString(fmt.Sprintf("    jz .endif%d\n", label))
			a.assembleBlock(instNode.ifNode.ifBlockNode)
			a.fileSb.WriteString(fmt.Sprintf(".endif%d:\n", label))
		}
	case INST_PRINT:
		a.assembleTerm(instNode.printNode.termNode)
		a.fileSb.WriteString("    mov rdi, 1\n")
		a.fileSb.WriteString("    mov rsi, rax\n")
		a.fileSb.WriteString("    call write_uint\n")
	}
}

func (a *Assembler) assembleBlock(blockNode BlockNode) {
	for _, inst := range blockNode.instructions {
		a.assembleInst(inst)
	}
}

func (a *Assembler) assembleExpr(exprNode ExprNode) {
	switch exprNode.exprType {
	case EXPR_TERM:
		a.assembleTerm(exprNode.termNode)
	case EXPR_PLUS:
		a.assembleTerm(exprNode.termBinaryNode.lhs)
		a.fileSb.WriteString("    mov rdx, rax\n")
		a.assembleTerm(exprNode.termBinaryNode.rhs)
		a.fileSb.WriteString("    add rax, rdx\n")
	}
}

func (a *Assembler) assembleTerm(termNode TermNode) {
	switch termNode.termType {
	case TERM_INPUT:
		a.fileSb.WriteString("    read 0, line, LINE_MAX\n")
		a.fileSb.WriteString("    mov rdi, line\n")
		a.fileSb.WriteString("    call strlen\n")
		a.fileSb.WriteString("    mov rdi, line\n")
		a.fileSb.WriteString("    mov rsi, rax\n")
		a.fileSb.WriteString("    call parse_uint\n")
	case TERM_INT:
		a.fileSb.WriteString(fmt.Sprintf("    mov rax, %s\n", termNode.value))
	case TERM_IDENT:
		index := a.findVariableIndex(termNode.value)
		a.fileSb.WriteString(fmt.Sprintf("    mov rax, qword [rbp - %d]\n", index*8+8))
	}
}

func (a *Assembler) assembleRel(relNode RelNode) {
	switch relNode.relType {
	case REL_LESS_THAN:
		a.assembleTerm(relNode.termBinaryNode.lhs)
		a.fileSb.WriteString("    mov rdx, rax\n")
		a.assembleTerm(relNode.termBinaryNode.rhs)
		a.fileSb.WriteString("    cmp rdx, rax\n")
		a.fileSb.WriteString("    setl al\n")
		a.fileSb.WriteString("    and al, 1\n")
		a.fileSb.WriteString("    movzx rax, al\n")
	}
}

func (a *Assembler) findVariableIndex(variable string) int {
	for i, varble := range a.variables {
		if varble == variable {
			return i
		}
	}
	return -1
}
