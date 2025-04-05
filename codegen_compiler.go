package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	inputFileName := os.Args[1]
	var outputFileName string
	if len(os.Args) < 3 {
		outputFileName = strings.Split(inputFileName, ".")[0]
	} else {
		outputFileName = os.Args[2]
	}

	buffer, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Print(err)
	}
	l := newLexer(string(buffer))
	tokens := l.tokenize()

	p := newParser(tokens)
	programNode := p.parseProgram()

	c := newCompiler(programNode)
	c.compileProgram()

	os.WriteFile("./"+outputFileName+".ll", []byte(c.module.String()), 0644)
}
