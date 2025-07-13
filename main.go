package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	glox_error "dsoechting/glox/error"
	"dsoechting/glox/interpret"
	"dsoechting/glox/parse"
	"dsoechting/glox/scanner"
)

type GloxError = glox_error.GloxError

type Glox struct {
	compileError error
	runtimeError error
}

func main() {
	args := os.Args[1:]
	argCount := len(args)
	glox := Glox{}

	if argCount > 1 {
		fmt.Println("Usage glox [script]")
		os.Exit(64)
	} else if argCount == 1 {
		glox.runFile(args[0])
	} else {
		glox.runPrompt()
	}
}

func (g *Glox) runFile(path string) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		// Can't read file
		os.Exit(66)
	}
	g.run(string(data))
	if g.compileError != nil {
		os.Exit(65)
	}
	if g.runtimeError != nil {
		os.Exit(70)
	}
}

func (g *Glox) runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if line == nil {
			break
		}
		g.run(string(line))
	}
	return nil
}

func (g *Glox) run(source string) {
	scanner := scanner.Create(source)
	tokens, scanErr := scanner.ScanTokens()
	if scanErr != nil {
		g.setCompileError(scanErr)
		return
	}

	// Token printing code
	// for _, token := range tokens {
	// 	log.Println(token)
	// }

	// printer := AstPrinter{}
	parser := parse.Create(tokens)
	// We need to make this static in the future, I am just hacking this in for now
	interpreter := interpret.Interpreter{}

	expression, parseError := parser.Parse()
	if parseError != nil {
		g.setCompileError(parseError)
		return
	}
	_, evalErr := interpreter.Interpret(expression)
	if evalErr != nil {
		g.setRuntimeError(evalErr)
	}
}

func (g *Glox) setCompileError(error error) {
	g.compileError = error
	log.Println(error.Error())
}

func (g *Glox) setRuntimeError(error error) {
	g.runtimeError = error
	log.Println(error.Error())
}
