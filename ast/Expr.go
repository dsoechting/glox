package ast

import "dsoechting/glox/token"

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinary(expr *BinaryExpr) any
	VisitGrouping(expr *GroupingExpr) any
	VisitLiteral(expr *LiteralExpr) any
	VisitUnary(expr *UnaryExpr) any
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(e)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(e)
}
