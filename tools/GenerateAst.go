package main

import (
	"fmt"
	"os"
	"strings"
)

var exprTypes = []string{
	"Ternary : Operator token.Token, First Expr, Second Expr, Third Expr",
	"Binary : Left Expr, Operator token.Token, Right Expr",
	"Grouping : Expression Expr",
	"Literal : Value any",
	"Unary : Operator token.Token, Right Expr",
	"Variable : Name token.Token",
}

var stmtTypes = []string{
	"Expression : Expression Expr",
	"Print : Expression Expr",
	"Var : Name token.Token, Initializer Expr",
}

const EXPR string = "Expr"
const STMT string = "Stmt"

func main() {
	args := os.Args[1:]
	argCount := len(args)

	if argCount != 1 {
		error := fmt.Errorf("Usage: generate_ast <output directory>")
		fmt.Println(error)
		os.Exit(64)
	}
	outputDir := args[0]

	defineAst(outputDir, EXPR, exprTypes)
	defineAst(outputDir, STMT, stmtTypes)
}

func defineAst(outputDir string, baseName string, types []string) {

	path := outputDir + "/" + baseName + ".go"
	var b strings.Builder

	fmt.Fprintln(&b, "package ast")
	fmt.Fprintln(&b, "")
	defineImports(&b, baseName)

	fmt.Fprintf(&b, "type %s interface {\n", baseName)
	fmt.Fprintf(&b, "	Accept(visitor %sVisitor) (any, error)", baseName)
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b, "")

	defineVisitor(&b, baseName, types)

	for _, v := range types {
		splits := strings.Split(v, ":")
		structName := strings.TrimSpace(splits[0])
		fields := strings.TrimSpace(splits[1])
		defineType(&b, baseName, structName, fields)
	}

	saveFile(path, b)

}

func defineImports(b *strings.Builder, baseName string) {
	switch baseName {
	case EXPR:
		fallthrough
	case STMT:
		fmt.Fprintln(b, "import \"dsoechting/glox/token\"")
		fmt.Fprintln(b, "")
		break
	}
}

func defineVisitor(b *strings.Builder, baseName string, types []string) {

	fmt.Fprintf(b, "type %sVisitor interface {\n", baseName)
	for _, v := range types {
		splits := strings.Split(v, ":")
		structName := strings.TrimSpace(splits[0])
		fmt.Fprintf(b, "	Visit%s(%s *%s%s) (any, error)\n", structName, strings.ToLower(baseName), structName, baseName)
	}
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")
}

func defineType(b *strings.Builder, baseName string, structName string, fieldList string) {

	fullStructName := structName + baseName
	fmt.Fprintf(b, "type %s struct {\n", fullStructName)

	fields := strings.Split(fieldList, ", ")
	for _, v := range fields {
		splitField := strings.Split(v, " ")
		fieldName := splitField[0]

		fieldType := splitField[1]
		fmt.Fprintf(b, "	%s %s\n", fieldName, fieldType)
	}
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")

	fmt.Fprintf(b, "func (e *%s) Accept(visitor %sVisitor) (any, error) {\n", fullStructName, baseName)
	fmt.Fprintf(b, "	return visitor.Visit%s(e)\n", structName)
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")

}

func saveFile(path string, b strings.Builder) {
	body := b.String()
	os.WriteFile(path, []byte(body), 0644)
}
