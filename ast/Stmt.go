package ast

import "dsoechting/glox/token"

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitBlock(stmt *BlockStmt) error
	VisitExpression(stmt *ExpressionStmt) error
	VisitPrint(stmt *PrintStmt) error
	VisitVar(stmt *VarStmt) error
}

type BlockStmt struct {
	Statements []Stmt
}

func (e *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlock(e)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpression(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrint(e)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (e *VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVar(e)
}
