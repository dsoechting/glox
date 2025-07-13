package test

import (
	"dsoechting/glox/interpret"
	"dsoechting/glox/parse"
	"dsoechting/glox/scanner"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type Interpreter = interpret.Interpreter

type expressionTestCase struct {
	name     string
	input    string
	expected string
}

func TestExpressions(t *testing.T) {
	// TODO: Dynamically iterate through all test data
	expressionTests := []expressionTestCase{
		{name: "Add numbers", input: "add_number_test", expected: "33"},
		{name: "Add strings", input: "add_strings_test", expected: "firstsecond"},
		{name: "Division", input: "div_test", expected: "11"},
		{name: "Group", input: "group_test", expected: "true"},
		{name: "Multiplication", input: "mult_test", expected: "27"},
		{name: "Subtract", input: "subtract_test", expected: "1"},
		{name: "Ternary", input: "ternary_test", expected: "27"},
	}

	for _, test := range expressionTests {
		inputText, parseErr := readExpressionSnippet(test.input)
		if parseErr != nil {
			t.Errorf("Failed to read test input: %v\nError: %v\n", test.name, parseErr)
		}
		// Ideally my interpreter tests wouldn't rely on the scanner and the parser
		// But I'm not typing out all of that test data
		scanner := scanner.Create(inputText)
		interpreter := Interpreter{}

		tokens, _ := scanner.ScanTokens()
		parser := parse.Create(tokens)
		expression, _ := parser.Parse()

		// Run test
		actual, evalErr := interpreter.Interpret(expression)

		if evalErr != nil {
			t.Errorf("Error while running test: %v\nError: %v\n", test.name, evalErr)
			continue

		}
		if actual != test.expected {
			t.Errorf("Test '%s' failed.\nExpected: %v\nActual: %v\n", test.name, test.expected, actual)
			continue
		}
		t.Logf("Test '%s' passed.\nExpected: %v\nActual: %v\n", test.name, test.expected, actual)
	}
}

func readExpressionSnippet(snippetName string) (string, error) {
	filePath := filepath.Join("data", fmt.Sprintf("%s.txt", snippetName))
	return readTestFile(filePath)
}

func readTestFile(path string) (string, error) {

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", err
	}

	contentString := string(contentBytes)

	return contentString, nil
}
