package main

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	IDENTIFIER         TokenType = "IDENTIFIER"
	LET                TokenType = "LET"
	IF                 TokenType = "IF"
	FOR                TokenType = "FOR"
	ELSE               TokenType = "ELSE"
	PRINT              TokenType = "PRINT"
	INPUT              TokenType = "INPUT"
	INT                TokenType = "INT"
	EQUAL              TokenType = "EQUAL"
	PLUS               TokenType = "PLUS"
	MINUS              TokenType = "MINUS"
	DIVIDE             TokenType = "DIVIDE"
	MULTIPLY           TokenType = "MULTIPLY"
	MODULO             TokenType = "MODULO"
	LESS_THAN          TokenType = "LESS_THAN"
	GREATER_THAN       TokenType = "GREATER_THAN"
	LESS_THAN_EQUAL    TokenType = "LESS_THAN_EQUAL"
	GREATER_THAN_EQUAL TokenType = "GREATER_THAN_EQUAL"
	EQUAL_EQUAL        TokenType = "EQUAL_EQUAL"
	NOT_EQUAL          TokenType = "NOT_EQUAL"
	AND                TokenType = "AND"
	OR                 TokenType = "OR"
	NOT                TokenType = "NOT"
	INVALID            TokenType = "INVALID"
	END                TokenType = "END"
	SPACE              TokenType = "SPACE"
	BLOCK_START        TokenType = "BLOCK_START"
	BLOCK_END          TokenType = "BLOCK_END"
	METHOD             TokenType = "METHOD"
	TYPE               TokenType = "TYPE"
	CLASS              TokenType = "CLASS"
	OPEN_PAREN         TokenType = "OPEN_PAREN"
	CLOSE_PAREN        TokenType = "CLOSE_PAREN"
	COLON              TokenType = "COLON"
	RETURN             TokenType = "RETURN"
)

type Token struct {
	tokenType TokenType
	value     string
}

type Lexer struct {
	buffer string
	pos    int
}

func (l Lexer) currChar() byte {
	return l.buffer[l.pos]
}

func (l Lexer) isBufferNotEmpty() bool {
	return l.pos < len(l.buffer)
}

func (l *Lexer) nextToken() Token {
	var value strings.Builder
	if unicode.IsSpace(rune(l.currChar())) {
		l.pos++
		return Token{SPACE, ""}
	} else if l.currChar() == '{' {
		l.pos++
		return Token{BLOCK_START, ""}
	} else if l.currChar() == '}' {
		l.pos++
		return Token{BLOCK_END, ""}
	} else if l.currChar() == '(' {
		l.pos++
		return Token{OPEN_PAREN, ""}
	} else if l.currChar() == ')' {
		l.pos++
		return Token{CLOSE_PAREN, ""}
	} else if l.currChar() == ':' {
		l.pos++
		return Token{COLON, ""}
	} else if l.currChar() == '=' {
		l.pos++
		return Token{EQUAL, ""}
	} else if l.currChar() == '+' {
		l.pos++
		return Token{PLUS, ""}
	} else if l.currChar() == '%' {
		l.pos++
		return Token{MODULO, ""}
	} else if l.currChar() == '-' {
		l.pos++
		return Token{MINUS, ""}
	} else if l.currChar() == '*' {
		l.pos++
		return Token{MULTIPLY, ""}
	} else if l.currChar() == '/' {
		l.pos++
		return Token{DIVIDE, ""}
	} else if l.currChar() == '<' {
		l.pos++
		return Token{LESS_THAN, ""}
	} else if unicode.IsDigit(rune(l.currChar())) {
		for l.isBufferNotEmpty() && unicode.IsDigit(rune(l.currChar())) {
			value.WriteString(string(l.currChar()))
			l.pos++
		}
		return Token{INT, value.String()}
	} else if unicode.IsLetter(rune(l.currChar())) || l.currChar() == '_' {
		for l.isBufferNotEmpty() && (unicode.IsLetter(rune(l.currChar())) || l.currChar() == '_') {
			value.WriteString(string(l.currChar()))
			l.pos++
		}
		if value.String() == "input" {
			return Token{INPUT, ""}
		} else if value.String() == "print" {
			return Token{PRINT, ""}
		} else if value.String() == "if" {
			return Token{IF, ""}
		} else if value.String() == "for" {
			return Token{FOR, ""}
		} else if value.String() == "else" {
			return Token{ELSE, ""}
		} else if value.String() == "let" {
			return Token{LET, ""}
		} else if value.String() == "method" {
			return Token{METHOD, ""}
		} else if value.String() == "class" {
			return Token{CLASS, ""}
		} else if value.String() == "return" {
			return Token{RETURN, ""}
		} else {
			return Token{IDENTIFIER, value.String()}
		}
	} else {
		value.WriteString(string(l.currChar()))
		l.pos++
		return Token{INVALID, value.String()}
	}
}

func (l Lexer) tokenize() []Token {
	var token Token
	var tokens []Token

	for l.isBufferNotEmpty() {
		token = l.nextToken()
		if token.tokenType != SPACE {
			// fmt.Println(token.tokenType)
			tokens = append(tokens, token)
		}

	}

	return tokens
}

func newLexer(buffer string) Lexer {
	return Lexer{buffer, 0}
}
