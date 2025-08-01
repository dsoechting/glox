package ast

import "dsoechting/glox/token"

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitBlock(stmt *BlockStmt) (any, error)
	VisitExpression(stmt *ExpressionStmt) (any, error)
	VisitIf(stmt *IfStmt) (any, error)
	VisitPrint(stmt *PrintStmt) (any, error)
	VisitVar(stmt *VarStmt) (any, error)
	VisitWhile(stmt *WhileStmt) (any, error)
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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIf(e)
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

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (e *WhileStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhile(e)
}
