package main

import (
	"dsoechting/glox/ast"
	"fmt"
	"log"
	"strings"
)

type Expr = ast.Expr
type TernaryExpr = ast.TernaryExpr
type BinaryExpr = ast.BinaryExpr
type UnaryExpr = ast.UnaryExpr
type GroupingExpr = ast.GroupingExpr
type LiteralExpr = ast.LiteralExpr

type AstPrinter struct{}

func (printer *AstPrinter) Print(expr Expr) string {
	if expr == nil {
		nilError := "The AST Printer was given a nil expression"
		log.Println(nilError)
		return nilError
	}
	result, printErr := expr.Accept(printer)
	if printErr != nil {
		log.Println(printErr.Error())
		return printErr.Error()
	}
	return result.(string)
}

func (printer *AstPrinter) VisitTernary(expr *TernaryExpr) (any, error) {
	return printer.parenthesizeTernary(expr.Operator.Lexeme, expr.First, expr.Second, expr.Third), nil
}

func (printer *AstPrinter) VisitBinary(expr *BinaryExpr) (any, error) {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (printer *AstPrinter) VisitGrouping(expr *GroupingExpr) (any, error) {
	return printer.parenthesize("group", expr.Expression), nil
}

func (printer *AstPrinter) VisitLiteral(expr *LiteralExpr) (any, error) {
	return fmt.Sprint(expr.Value), nil
}

func (printer *AstPrinter) VisitUnary(expr *UnaryExpr) (any, error) {
	return printer.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (printer *AstPrinter) parenthesize(name string, exprs ...Expr) string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		sb.WriteString(" ")
		res, _ := expr.Accept(printer)
		sb.WriteString(res.(string))
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *AstPrinter) parenthesizeTernary(name string, one Expr, two Expr, three Expr) string {
	var sb strings.Builder

	oneStr, _ := one.Accept(p)
	twoStr, _ := two.Accept(p)
	threeStr, _ := three.Accept(p)

	sb.WriteString(fmt.Sprintf("(%s ", oneStr.(string)))
	sb.WriteString(fmt.Sprintf("%s ", name))
	sb.WriteString(fmt.Sprintf("%s ", twoStr.(string)))
	sb.WriteString(": ")
	sb.WriteString(fmt.Sprintf("%s )", threeStr.(string)))

	return sb.String()
}
