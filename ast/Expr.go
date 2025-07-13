package ast

import "dsoechting/glox/token"

type Token = token.Token

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

type ExprVisitor interface {
	VisitTernary(expr *TernaryExpr) (any, error)
	VisitBinary(expr *BinaryExpr) (any, error)
	VisitGrouping(expr *GroupingExpr) (any, error)
	VisitLiteral(expr *LiteralExpr) (any, error)
	VisitUnary(expr *UnaryExpr) (any, error)
}

type TernaryExpr struct {
	Operator Token
	First    Expr
	Second   Expr
	Third    Expr
}

func (e *TernaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitTernary(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
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
	Operator Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(e)
}
