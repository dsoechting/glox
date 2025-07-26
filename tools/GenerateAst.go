package main

import (
	"fmt"
	"os"
	"strings"
)

type Ast struct {
	types       []string
	base_name   string
	return_type string //Is there a way to use type parameters here?
}

var exprTypes = []string{
	"Ternary : Operator token.Token, First Expr, Second Expr, Third Expr",
	"Assign : Name token.Token, Value Expr",
	"Binary : Left Expr, Operator token.Token, Right Expr",
	"Grouping : Expression Expr",
	"Literal : Value any",
	"Unary : Operator token.Token, Right Expr",
	"Variable : Name token.Token",
}

var stmtTypes = []string{
	"Block : Statements []Stmt",
	"Expression : Expression Expr",
	"Print : Expression Expr",
	"Var : Name token.Token, Initializer Expr",
}

const EXPR string = "Expr"
const STMT string = "Stmt"

var EXPR_AST = Ast{
	types:       exprTypes,
	base_name:   "Expr",
	return_type: "any, error",
}
var STMT_AST = Ast{
	types:       stmtTypes,
	base_name:   "Stmt",
	return_type: "error",
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

	defineAst(EXPR_AST, outputDir)
	defineAst(STMT_AST, outputDir)
}

func defineAst(ast Ast, outputDir string) {

	path := outputDir + "/" + ast.base_name + ".go"
	var b strings.Builder

	fmt.Fprintln(&b, "package ast")
	fmt.Fprintln(&b, "")
	defineImports(&b, ast)

	fmt.Fprintf(&b, "type %s interface {\n", ast.base_name)
	fmt.Fprintf(&b, "	Accept(visitor %sVisitor) (%s)", ast.base_name, ast.return_type)
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b, "")

	defineVisitor(&b, ast)
	defineTypes(&b, ast)

	saveFile(path, b)

}

func defineImports(b *strings.Builder, ast Ast) {
	switch ast.base_name {
	case EXPR:
		fallthrough
	case STMT:
		fmt.Fprintln(b, "import \"dsoechting/glox/token\"")
		fmt.Fprintln(b, "")
		break
	}
}

func defineVisitor(b *strings.Builder, ast Ast) {

	fmt.Fprintf(b, "type %sVisitor interface {\n", ast.base_name)
	for _, v := range ast.types {
		splits := strings.Split(v, ":")
		structName := strings.TrimSpace(splits[0])
		fmt.Fprintf(b, "	Visit%s(%s *%s%s) (%s)\n", structName, strings.ToLower(ast.base_name), structName, ast.base_name, ast.return_type)
	}
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")
}

func defineTypes(b *strings.Builder, ast Ast) {
	for _, v := range ast.types {
		splits := strings.Split(v, ":")
		structName := strings.TrimSpace(splits[0])
		fields := strings.TrimSpace(splits[1])
		defineType(b, ast, structName, fields)
	}
}

func defineType(b *strings.Builder, ast Ast, structName string, fieldList string) {

	fullStructName := structName + ast.base_name
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

	fmt.Fprintf(b, "func (e *%s) Accept(visitor %sVisitor) (%s) {\n", fullStructName, ast.base_name, ast.return_type)
	fmt.Fprintf(b, "	return visitor.Visit%s(e)\n", structName)
	fmt.Fprintln(b, "}")
	fmt.Fprintln(b, "")

}

func saveFile(path string, b strings.Builder) {
	body := b.String()
	os.WriteFile(path, []byte(body), 0644)
}
