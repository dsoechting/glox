package ast

import "dsoechting/glox/token"

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

type ExprVisitor interface {
	VisitTernary(expr *TernaryExpr) (any, error)
	VisitAssign(expr *AssignExpr) (any, error)
	VisitBinary(expr *BinaryExpr) (any, error)
	VisitGrouping(expr *GroupingExpr) (any, error)
	VisitLiteral(expr *LiteralExpr) (any, error)
	VisitUnary(expr *UnaryExpr) (any, error)
	VisitVariable(expr *VariableExpr) (any, error)
}

type TernaryExpr struct {
	Operator token.Token
	First    Expr
	Second   Expr
	Third    Expr
}

func (e *TernaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitTernary(e)
}

type AssignExpr struct {
	Name  token.Token
	Value Expr
}

func (e *AssignExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssign(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(e)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(e)
}

type VariableExpr struct {
	Name token.Token
}

func (e *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariable(e)
}
