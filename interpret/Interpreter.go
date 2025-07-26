package interpret

import (
	"dsoechting/glox/ast"
	"dsoechting/glox/environment"
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Environment = environment.Environment
type Stmt = ast.Stmt
type ExpressionStmt = ast.ExpressionStmt
type IfStmt = ast.IfStmt
type PrintStmt = ast.PrintStmt
type WhileStmt = ast.WhileStmt
type VarStmt = ast.VarStmt
type BlockStmt = ast.BlockStmt
type Expr = ast.Expr
type TernaryExpr = ast.TernaryExpr
type BinaryExpr = ast.BinaryExpr
type LogicalExpr = ast.LogicalExpr
type UnaryExpr = ast.UnaryExpr
type VariableExpr = ast.VariableExpr
type GroupingExpr = ast.GroupingExpr
type LiteralExpr = ast.LiteralExpr
type Token = token.Token
type GloxError = glox_error.GloxError

// Implements ExprVisitor and StmtVisitor
type Interpreter struct {
	environment Environment
}

func Create() Interpreter {
	env := environment.Create()
	return Interpreter{
		environment: env,
	}
}

func (i *Interpreter) Interpret(statements []Stmt) (string, error) {
	var sb strings.Builder

	for _, statement := range statements {
		value, err := i.execute(statement)
		if err != nil {
			return "", err
		}
		// Do no add new line for empty values
		if value != "" {
			sb.WriteString(fmt.Sprintf("%v\n", value))
		}
	}
	// This is going to break my tests :)
	// Maybe we will modify this later, to make things easier to test
	return sb.String(), nil
}

func (i *Interpreter) VisitExpression(stmt *ExpressionStmt) (any, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrint(stmt *PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	//Don't print in REPL
	return "", nil
}

func (i *Interpreter) VisitWhile(stmt *WhileStmt) (any, error) {
	cond, condErr := i.evaluate(stmt.Condition)
	if condErr != nil {
		return nil, condErr
	}

	for isTruthy(cond) {
		i.execute(stmt.Body)
		cond, condErr = i.evaluate(stmt.Condition)
		if condErr != nil {
			return nil, condErr
		}
	}
	return "", nil
}

func (i *Interpreter) VisitIf(stmt *IfStmt) (any, error) {
	cond, condErr := stmt.Condition.Accept(i)
	if condErr != nil {
		return nil, condErr
	}
	if isTruthy(cond) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return "", nil
}

func (i *Interpreter) VisitVar(stmt *VarStmt) (any, error) {
	var value any
	var err error
	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	// Don't print in REPL
	return "", nil
}

func (i *Interpreter) VisitTernary(expr *TernaryExpr) (any, error) {
	first, firstExpErr := i.evaluate(expr.First)
	if firstExpErr != nil {
		return nil, firstExpErr
	}

	if isTruthy(first) {
		return i.evaluate(expr.Second)
	} else {
		return i.evaluate(expr.Third)
	}
}

func (i *Interpreter) VisitBinary(expr *BinaryExpr) (any, error) {
	left, leftErr := i.evaluate(expr.Left)
	right, rightErr := i.evaluate(expr.Right)

	if leftErr != nil || rightErr != nil {
		return nil, errors.Join(leftErr, rightErr)
	}

	switch expr.Operator.TokenType {
	case token.EQUAL_EQUAL:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return isEqual(left, right), nil
	case token.BANG_EQUAL:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return !isEqual(left, right), nil
	case token.GREATER:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) <= right.(float64), nil
	case token.MINUS:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) - right.(float64), nil
	case token.STAR:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		operandsErr := checkNumberOperands(expr.Operator, left, right)
		if operandsErr != nil {
			return nil, operandsErr
		}
		return left.(float64) / right.(float64), nil
	case token.PLUS:
		leftStr, isLeftStr := left.(string)
		rightStr, isRightStr := right.(string)
		leftFloat, isLeftFloat := left.(float64)
		rightFloat, isRightFloat := right.(float64)

		if isLeftStr && isRightStr {
			return leftStr + rightStr, nil
		}
		if isLeftFloat && isRightFloat {
			return leftFloat + rightFloat, nil
		}
		return nil, createInterpreterError(expr.Operator, "Operands must be two numbers or string", left, right)
	}
	return nil, fmt.Errorf("Unsupporter binary operator %s\n", expr.Operator.TokenType)
}

func (i *Interpreter) VisitGrouping(expr *GroupingExpr) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteral(expr *LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogical(expr *LogicalExpr) (any, error) {
	left, leftErr := i.evaluate(expr.Left)
	if leftErr != nil {
		return nil, leftErr
	}

	if expr.Operator.TokenType == token.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitUnary(expr *UnaryExpr) (any, error) {
	var unaryError error
	right, rightErr := i.evaluate(expr.Right)
	if rightErr != nil {
		return nil, rightErr
	}

	switch expr.Operator.TokenType {
	case token.MINUS:
		unaryError = checkNumberOpernad(expr.Operator, right)
		if unaryError != nil {
			return nil, unaryError
		}
		return -right.(float64), nil
	case token.BANG:
		unaryError = checkNumberOpernad(expr.Operator, right)
		if unaryError != nil {
			return nil, unaryError
		}
		return !isTruthy(right), nil
	}
	// We should be unreachable here
	return nil, fmt.Errorf("Invalid Unary operator %s\n", expr.Operator.TokenType)

}

func (i *Interpreter) VisitVariable(expr *VariableExpr) (any, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssign(expr *ast.AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	assignErr := i.environment.Assign(expr.Name, value)
	if assignErr != nil {
		return nil, assignErr
	}

	// Don't want my REPL to print assignments
	return "", nil
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBlock(stmt *BlockStmt) (any, error) {
	blockEnv := environment.CreateWithEnclosing(i.environment)
	return i.executeBlock(stmt.Statements, blockEnv)
}

func (i *Interpreter) executeBlock(statements []Stmt, blockEnv Environment) (any, error) {
	previous := i.environment

	i.environment = blockEnv

	defer func() {
		i.environment = previous
	}()

	var sb strings.Builder
	for _, stmt := range statements {
		value, executeErr := i.execute(stmt)
		if executeErr != nil {
			return nil, executeErr
		}
		sb.WriteString(fmt.Sprintf("%v\n", value))
	}

	return sb.String(), nil
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	b, ok := value.(bool)
	if ok {
		return b
	}
	return true
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	// please don't be slices, maps, or functions here
	return a == b
}

func checkNumberOpernad(operator Token, operand any) error {
	_, isFloat := operand.(float64)
	if isFloat {
		return nil
	}
	return createInterpreterError(operator, "Operand must be a number", operand)
}

func checkNumberOperands(operator Token, left any, right any) error {
	_, isLeftFloat := left.(float64)
	_, isRightFloat := right.(float64)
	if isLeftFloat && isRightFloat {
		return nil
	}

	return createInterpreterError(operator, "Operands must be numbers", left, right)
}

func createInterpreterError(operator Token, message string, operands ...any) *GloxError {
	return glox_error.Create(operator.Line, fmt.Sprintf("%s on %s", operands, operator.Lexeme), message)
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	floatVal, isFloat := object.(float64)
	if isFloat {
		text := strconv.FormatFloat(floatVal, 'f', -1, 64)
		return text
	}
	return fmt.Sprintf("%v", object)

}
