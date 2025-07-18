package ast

import "dsoechting/glox/token"

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitBlock(stmt *BlockStmt) (any, error)
	VisitExpression(stmt *ExpressionStmt) (any, error)
	VisitPrint(stmt *PrintStmt) (any, error)
	VisitVar(stmt *VarStmt) (any, error)
}

type BlockStmt struct {
	Statements []Stmt
}

func (e *BlockStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlock(e)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpression(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e *PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrint(e)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (e *VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVar(e)
}
