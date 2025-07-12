package main

import (
	"dsoechting/glox/ast"
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (printer *AstPrinter) Print(expr ast.Expr) string {
	if expr == nil {
		return fmt.Sprintln("The AST Printer was given a null expression")
	}
	return expr.Accept(printer).(string)
}

func (printer *AstPrinter) VisitBinary(expr *ast.BinaryExpr) any {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (printer *AstPrinter) VisitGrouping(expr *ast.GroupingExpr) any {
	return printer.parenthesize("group", expr.Expression)
}

func (printer *AstPrinter) VisitLiteral(expr *ast.LiteralExpr) any {
	return fmt.Sprint(expr.Value)
}

func (printer *AstPrinter) VisitUnary(expr *ast.UnaryExpr) any {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (printer *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.Accept(printer).(string))
	}
	sb.WriteString(")")
	return sb.String()
}
