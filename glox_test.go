package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type expressionTestCase struct {
	name     string
	input    string
	expected string
}

// I will probably need to rework this later.
// The run() function will not always return the stringified interprete result
// But it is nice to have the the Scanner and Parser do a lot of work for me.
// Maybe I have Scanner and Parser tests like this, and only progress to the next test if the prev passes
func TestExpressions(t *testing.T) {
	// Can I auto parse the files from this dir?
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

		actual, evalErr := run(inputText)
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
	filePath := filepath.Join("test", "expressions", fmt.Sprintf("%s.txt", snippetName))
	return readTestFile(filePath)
}

func readTestFile(path string) (string, error) {

	// Read the entire file content into a byte slice
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", err
	}

	// Convert the byte slice to a string
	contentString := string(contentBytes)

	return contentString, nil
}
