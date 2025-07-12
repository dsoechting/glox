package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"dsoechting/glox/parse"
	"dsoechting/glox/scanner"
)

type Glox struct {
	hadError bool
}

func main() {
	args := os.Args[1:]
	argCount := len(args)

	if argCount > 1 {
		fmt.Println("Usage glox [script]")
		os.Exit(64)
	} else if argCount == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		// Can't read file
		os.Exit(66)
	}
	runErr := run(string(data))
	if runErr != nil {
		fmt.Println(runErr.Error())
		os.Exit(65)
	}
}

func runPrompt() error {
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
		runErr := run(string(line))
		if runErr != nil {
			fmt.Println(runErr.Error())
		}
	}
	return nil
}

func run(source string) error {
	scanner := scanner.Create(source)
	tokens, scanErr := scanner.ScanTokens()
	// Token printing code
	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	printer := AstPrinter{}
	parser := parse.Create(tokens)
	expression, parseError := parser.Parse()
	if parseError != nil {
		log.Println(parseError)
	} else {
		log.Println(printer.Print(expression))
	}
	return scanErr
}
