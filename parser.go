package main

import "slices"

type InstType string

const (
	INST_ASSIGN InstType = "INST_ASSIGN"
	INST_IF     InstType = "INST_IF"
	INST_PRINT  InstType = "INST_PRINT"
	INST_ELSE   InstType = "INST_ELSE"
	INST_METHOD InstType = "INST_METHOD"
	INST_CLASS  InstType = "INST_CLASS"
	INST_RETURN InstType = "INST_RETURN"
)

type ExprType string

const (
	EXPR_TERM ExprType = "EXPR_TERM"
	EXPR_PLUS ExprType = "EXPR_PLUS"
)

type RelType string

const (
	REL_LESS_THAN RelType = "REL_LESS_THAN"
)

type TermType string

const (
	TERM_INPUT TermType = "TERM_INPUT"
	TERM_INT   TermType = "TERM_INT"
	TERM_IDENT TermType = "TERM_IDENT"
)

type ExprNode struct {
	exprType       ExprType
	termBinaryNode TermBinaryNode
	termNode       TermNode
}

type TermBinaryNode struct {
	rhs TermNode
	lhs TermNode
}

type RelNode struct {
	relType        RelType
	termBinaryNode TermBinaryNode
}

type TermNode struct {
	termType TermType
	value    string
}

type AssignNode struct {
	identifier string
	typeName   string
	expr       ExprNode
}

type IfNode struct {
	relNode       RelNode
	ifBlockNode   BlockNode
	elseBlockNode BlockNode
}

type PrintNode struct {
	termNode TermNode
}

type BlockNode struct {
	instructions []InstNode
}

type ClassNode struct {
	className     string
	functionNames []string
	varNames      []string
	blockNode     BlockNode
}

type MethodNode struct {
	methodName string
	parameters map[string]string
	returnType string
	varNames   []string
	blockNode  BlockNode
}

type ReturnNode struct {
	exprNode ExprNode
}

type InstNode struct {
	instType   InstType
	assignNode AssignNode
	ifNode     IfNode
	printNode  PrintNode
	methodNode MethodNode
	classNode  ClassNode
	returnNode ReturnNode
}

type ProgramNode struct {
	instructions []InstNode
	fileName     string
}

type Parser struct {
	tokens []Token
	index  int
}

func (b BlockNode) getFunctionNames() []string {
	functionNames := []string{}
	return functionNames
}

func (b BlockNode) getVarNames() []string {
	varNames := []string{}
	return varNames
}

func newParser(tokens []Token) Parser {
	return Parser{tokens, 0}
}

func (p Parser) parserCurrent() Token {
	if p.index < len(p.tokens) {
		return p.tokens[p.index]
	}
	return Token{tokenType: END}
}

func (p *Parser) parserAdvance() {
	if p.index >= len(p.tokens) {
		panic("Error finished all tokens")
	}
	p.index++
}

func (p *Parser) parseParameters() map[string]string {
	parameters := make(map[string]string)
	token := p.parserCurrent()
	var key string
	var val string
	for token.tokenType != CLOSE_PAREN {
		p.parserAdvance()
		key = p.parserCurrent().value
		p.parserAdvance()
		if p.parserCurrent().tokenType != COLON {
			panic("Expected : in method parameter decleration")
		}
		p.parserAdvance()
		val = p.parserCurrent().value
		parameters[key] = val
		p.parserAdvance()
		token = p.parserCurrent()
	}
	p.parserAdvance()

	return parameters
}

func (p *Parser) parseTerm() TermNode {
	token := p.parserCurrent()
	termNode := TermNode{}
	if token.tokenType == INPUT {
		termNode.termType = TERM_INPUT
	} else if token.tokenType == INT {
		termNode.termType = TERM_INT
		termNode.value = token.value
	} else if token.tokenType == IDENTIFIER {
		termNode.termType = TERM_IDENT
		termNode.value = token.value
	} else {
		panic("Expected a term (input, int or ident) but found " + token.tokenType)
	}
	p.parserAdvance()
	return termNode
}

func (p *Parser) parseExpr() ExprNode {
	exprNode := ExprNode{}
	lhs := p.parseTerm()
	token := p.parserCurrent()

	if token.tokenType == PLUS {
		p.parserAdvance()
		rhs := p.parseTerm()
		exprNode.exprType = EXPR_PLUS
		exprNode.termBinaryNode.lhs = lhs
		exprNode.termBinaryNode.rhs = rhs
	} else {
		exprNode.exprType = EXPR_TERM
		exprNode.termNode = lhs
	}
	return exprNode
}

func (p *Parser) parseRel() RelNode {
	relNode := RelNode{}
	lhs := p.parseTerm()
	token := p.parserCurrent()

	if token.tokenType == LESS_THAN {
		p.parserAdvance()
		rhs := p.parseTerm()
		relNode.relType = REL_LESS_THAN
		relNode.termBinaryNode.lhs = lhs
		relNode.termBinaryNode.rhs = rhs
	} else {
		panic("Expected relational (<) found " + token.tokenType)
	}
	return relNode
}

func (p *Parser) parseBlock() BlockNode {
	blockNode := BlockNode{}
	currentNode := p.parserCurrent()
	for {
		if currentNode.tokenType == BLOCK_END {
			p.parserAdvance()
			return blockNode
		} else {
			instNode := p.parseInst()
			if instNode.instType != "" {
				blockNode.instructions = append(blockNode.instructions, instNode)
			}
			currentNode = p.parserCurrent()
		}
	}
}

func (p *Parser) checkTypeValid(typeName string) {
	primitaveTypes := []string{"int", "string"}
	// TODO: add more primitaves and user types
	if !slices.Contains(primitaveTypes, typeName) {
		panic("Unknown type: " + typeName)
	}
}

func (p *Parser) parseAssign() InstNode {
	p.parserAdvance()
	token := p.parserCurrent()
	instNode := InstNode{}
	instNode.instType = INST_ASSIGN
	instNode.assignNode.typeName = token.value
	p.checkTypeValid(token.value)
	p.parserAdvance()
	token = p.parserCurrent()
	instNode.assignNode.identifier = token.value
	p.parserAdvance()
	token = p.parserCurrent()
	if token.tokenType != EQUAL {
		panic("Expected equal but found " + token.tokenType)
	}
	p.parserAdvance()
	instNode.assignNode.expr = p.parseExpr()
	return instNode
}

func (p *Parser) parseIf() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_IF
	p.parserAdvance()
	instNode.ifNode.relNode = p.parseRel()
	token := p.parserCurrent()
	if token.tokenType == BLOCK_START {
		p.parserAdvance()
		instNode.ifNode.ifBlockNode = p.parseBlock()
	}
	if p.parserCurrent().tokenType == ELSE {
		p.parserAdvance()
		if token.tokenType == BLOCK_START {
			p.parserAdvance()
			instNode.ifNode.elseBlockNode = p.parseBlock()
		}
	}
	return instNode
}

func (p *Parser) parsePrint() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_PRINT
	p.parserAdvance()
	instNode.printNode.termNode = p.parseTerm()
	return instNode
}

func (p *Parser) parseClass() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_CLASS
	p.parserAdvance()
	nameToken := p.parserCurrent()
	p.parserAdvance()
	classBlockNode := p.parseBlock()
	instNode.classNode.blockNode = classBlockNode
	instNode.classNode.className = nameToken.value
	instNode.classNode.functionNames = classBlockNode.getFunctionNames()
	instNode.classNode.varNames = classBlockNode.getVarNames()
	return instNode
}

func (p *Parser) parseMethod() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_METHOD
	p.parserAdvance()
	nameToken := p.parserCurrent()
	p.parserAdvance()
	if p.parserCurrent().tokenType == OPEN_PAREN {
		p.parserAdvance()
		instNode.methodNode.parameters = p.parseParameters()
	} else {
		panic("Expected ( in method definition")
	}
	if p.parserCurrent().tokenType == COLON {
		p.parserAdvance()
		instNode.methodNode.returnType = p.parserCurrent().value
		p.parserAdvance()
	} else {
		instNode.methodNode.returnType = "void"
	}

	methodBlockNode := p.parseBlock()
	instNode.methodNode.blockNode = methodBlockNode
	instNode.methodNode.methodName = nameToken.value
	instNode.methodNode.varNames = methodBlockNode.getVarNames()
	return instNode
}

func (p *Parser) parseReturn() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_RETURN
	p.parserAdvance()
	instNode.returnNode.exprNode = p.parseExpr()
	return instNode
}

func (p *Parser) parseInst() InstNode {
	var token Token
	var instNode InstNode
	token = p.parserCurrent()
	switch token.tokenType {
	case LET:
		instNode = p.parseAssign()
	case IF:
		instNode = p.parseIf()
	case PRINT:
		instNode = p.parsePrint()
	case CLASS:
		instNode = p.parseClass()
	case METHOD:
		instNode = p.parseMethod()
	case RETURN:
		instNode = p.parseReturn()
	default:
		p.parserAdvance()
	}
	return instNode
}

func (p *Parser) parseProgram() ProgramNode {
	programNode := ProgramNode{[]InstNode{}, "test.yeol"}

	var instNode InstNode
	for p.index < len(p.tokens) {
		instNode = p.parseInst()
		// fmt.Println(instNode.instType)
		if instNode.instType != "" {
			// fmt.Println(instNode.instType)
			programNode.instructions = append(programNode.instructions, instNode)
		}
	}

	return programNode
}
