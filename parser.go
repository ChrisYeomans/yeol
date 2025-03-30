package main

type InstType string

const (
	INST_ASSIGN InstType = "INST_ASSIGN"
	INST_IF     InstType = "INST_IF"
	INST_PRINT  InstType = "INST_PRINT"
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
	expr       ExprNode
}

type IfNode struct {
	relNode  RelNode
	instNode *InstNode
}

type PrintNode struct {
	termNode TermNode
}

type InstNode struct {
	instType   InstType
	assignNode AssignNode
	ifNode     IfNode
	printNode  PrintNode
}

type ProgramNode struct {
	instructions []InstNode
	fileName     string
}

type Parser struct {
	tokens []Token
	index  int
}

func newParser(tokens []Token) Parser {
	return Parser{tokens, 0}
}

func (p Parser) parserCurrent() Token {
	return p.tokens[p.index]
}

func (p *Parser) parserAdvance() {
	p.index++
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
		panic("Expected reational (<) found " + token.tokenType)
	}
	return relNode
}

func (p *Parser) parseAssign() InstNode {
	token := p.parserCurrent()
	instNode := InstNode{}
	instNode.instType = INST_ASSIGN
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
	p.parserAdvance()
	ifInstNode := p.parseInst()
	instNode.ifNode.instNode = &ifInstNode
	return instNode
}

func (p *Parser) parsePrint() InstNode {
	instNode := InstNode{}
	instNode.instType = INST_PRINT
	p.parserAdvance()
	instNode.printNode.termNode = p.parseTerm()
	return instNode
}

func (p *Parser) parseInst() InstNode {
	var token Token
	var instNode InstNode
	token = p.parserCurrent()
	if token.tokenType == IDENTIFIER {
		instNode = p.parseAssign()
	} else if token.tokenType == IF {
		instNode = p.parseIf()
	} else if token.tokenType == PRINT {
		instNode = p.parsePrint()
	} else if token.tokenType == LET {
		p.parserAdvance()
	}
	return instNode
}

func (p *Parser) parseProgram() ProgramNode {
	programNode := ProgramNode{[]InstNode{}, "test.yeol"}

	var instNode InstNode
	for p.index < len(p.tokens) {
		instNode = p.parseInst()
		if instNode.instType != "" {
			// fmt.Println(instNode.instType)
			programNode.instructions = append(programNode.instructions, instNode)
		}
	}

	return programNode
}
