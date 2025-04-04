package main

import (
	"fmt"
	"os"
	"os/exec"
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

	a := Assembler{programNode, []string{}, 0, strings.Builder{}}
	a.assembleProgram()

	os.WriteFile("./tmp.asm", []byte(a.fileSb.String()), 0644)
	cmd := exec.Command("fasm", "./tmp.asm -s", outputFileName)
	stdout, _ := cmd.Output()
	fmt.Println(string(stdout))
	cmd = exec.Command("rm", "tmp.asm")
	stdout, _ = cmd.Output()
	fmt.Println(string(stdout))
	//	stdout, _ := cmd.Output()
	// os.WriteFile("./test", stdout, 0644)
	//	fmt.Println(string(stdout))
	// fmt.Print(a.fileSb.String())
}
