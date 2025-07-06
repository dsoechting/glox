package main

import (
	"fmt"
	"os"
	"strings"
)

var exprTypes = []string{
	"Binary : left Expr, operator token.Token, right Expr",
	"Grouping : expression Expr",
	"Literal : value token.Token",
	"Unary : operator token.Token, right Expr",
}

func main() {
	args := os.Args[1:]
	argCount := len(args)

	if argCount != 1 {
		error := fmt.Errorf("Usage: generate_ast <output directory>")
		fmt.Println(error)
		os.Exit(64)
	}
	outputDir := args[0]

	defineAst(outputDir, "Expr", exprTypes)
}

func defineAst(outputDir string, baseName string, types []string) {

	path := outputDir + "/" + baseName + ".go"
	var b strings.Builder

	fmt.Fprintln(&b, "package ast")
	fmt.Fprintln(&b, "")
	fmt.Fprintln(&b, "import \"dsoechting/glox/token\"")
	fmt.Fprintln(&b, "")

	fmt.Fprintln(&b, "type Expr interface {")
	fmt.Fprintln(&b, "	Accept(visitor ExprVisitor[any]) any")
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

func defineVisitor(b *strings.Builder, baseName string, types []string) {

	// type ExprVisitor interface {
	// 	VisitBinary(expr *BinaryExpr)
	// 	VisitGrouping(expr *GroupingExpr)
	// 	VisitLiteral(expr *LiteralExpr)
	// 	VisitUnary(expr *UnaryExpr)
	// }

	fmt.Fprintf(b, "type %sVisitor[T any] interface {\n", baseName)
	for _, v := range types {
		splits := strings.Split(v, ":")
		structName := strings.TrimSpace(splits[0])
		fmt.Fprintf(b, "	Visit%s(expr *%s%s) T\n", structName, structName, baseName)
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
		fieldName := strings.Title(strings.ToLower(splitField[0]))
		fieldType := splitField[1]
		fmt.Fprintf(b, "	%s %s\n", fieldName, fieldType)
	}
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")

	fmt.Fprintf(b, "func (e *%s) Accept(visitor ExprVisitor[any]) any {\n", fullStructName)
	fmt.Fprintf(b, "	return visitor.Visit%s(e)\n", structName)
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")

}

func saveFile(path string, b strings.Builder) {
	body := b.String()
	os.WriteFile(path, []byte(body), 0644)
}
