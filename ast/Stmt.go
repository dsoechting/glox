package ast

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpression(stmt *ExpressionStmt) (any, error)
	VisitPrint(stmt *PrintStmt) (any, error)
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
