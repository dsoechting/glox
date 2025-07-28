package parse

import (
	"dsoechting/glox/ast"
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
	"fmt"
	"log"
)

type Expr = ast.Expr
type Stmt = ast.Stmt
type FunctionStmt = ast.FunctionStmt
type WhileStmt = ast.WhileStmt
type VarStmt = ast.VarStmt
type ExpressionStmt = ast.ExpressionStmt
type TernaryExpr = ast.TernaryExpr
type BinaryExpr = ast.BinaryExpr
type LogicalExpr = ast.LogicalExpr
type UnaryExpr = ast.UnaryExpr
type LiteralExpr = ast.LiteralExpr
type GroupingExpr = ast.GroupingExpr
type CallExpr = ast.CallExpr
type VariableExpr = ast.VariableExpr
type TokenType = token.TokenType
type Token = token.Token
type GloxError = glox_error.GloxError

type Parser struct {
	tokens  []token.Token
	current int
}

func Create(tokens []token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, stmtErr := p.declaration()
		if stmtErr != nil {
			// Do I want to end parsing here? I might want to keep going
			return nil, stmtErr
		}
		statements = append(statements, stmt)

	}
	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(token.FUN) {
		funcVal, funcErr := p.function("function")
		if funcErr != nil {
			return nil, funcErr
		}
	}
	if p.match(token.VAR) {
		varDecl, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		return varDecl, nil
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.LEFT_BRACE) {
		blockStmts, blockErr := p.block()
		if blockErr != nil {
			return nil, blockErr
		}

		return &ast.BlockStmt{
			Statements: blockStmts,
		}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	_, leftParenErr := p.consume(token.LEFT_PAREN, "Expect '(' after 'for',")
	if leftParenErr != nil {
		return nil, leftParenErr
	}

	var initializer Stmt
	var initializerErr error
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, initializerErr = p.varDeclaration()
	} else {
		initializer, initializerErr = p.expressionStatement()
	}
	if initializerErr != nil {
		return nil, initializerErr
	}

	var condition Expr
	var conditionErr error
	if !p.check(token.SEMICOLON) {
		condition, conditionErr = p.expression()
	}
	if conditionErr != nil {
		return nil, conditionErr
	}

	_, semiColonErr := p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if semiColonErr != nil {
		return nil, semiColonErr
	}

	var increment Expr
	var incrementErr error
	if !p.check(token.RIGHT_PAREN) {
		increment, incrementErr = p.expression()
	}
	if incrementErr != nil {
		return nil, incrementErr
	}

	_, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if rightParenErr != nil {
		return nil, rightParenErr
	}

	body, bodyErr := p.statement()
	if bodyErr != nil {
		return nil, bodyErr
	}

	// Append our incrementer to the end of the body
	if increment != nil {
		incrementStmt := &ExpressionStmt{
			Expression: increment,
		}
		body = &ast.BlockStmt{
			Statements: []Stmt{body, incrementStmt},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{
			Value: true,
		}
	}

	// Make a while loop with our condition and the body (with incrementer)
	body = &WhileStmt{
		Condition: condition,
		Body:      body,
	}

	// Prepend the body and run the initializer once before the while loop
	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {

	_, leftParenErr := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if leftParenErr != nil {
		return nil, leftParenErr
	}

	// Condition
	condition, condErr := p.expression()
	if condErr != nil {
		return nil, condErr
	}

	_, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if rightParenErr != nil {
		return nil, rightParenErr
	}

	// Then branch of code
	thenBranch, thenBranchErr := p.statement()
	if thenBranchErr != nil {
		return nil, thenBranchErr
	}

	// Else block if present
	var maybeElseBranch Stmt
	if p.match(token.ELSE) {
		elseBranch, elseBranchErr := p.statement()
		if elseBranchErr != nil {
			return nil, elseBranchErr
		}
		maybeElseBranch = elseBranch
	}
	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		// Possibly nil
		ElseBranch: maybeElseBranch,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, semiColonErr := p.consume(token.SEMICOLON, "Expect ';' after value")
	if semiColonErr != nil {
		return nil, semiColonErr
	}

	return &ast.PrintStmt{
		Expression: value,
	}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {

	_, leftParenErr := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if leftParenErr != nil {
		return nil, leftParenErr
	}

	// Condition
	condition, condErr := p.expression()
	if condErr != nil {
		return nil, condErr
	}

	_, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after while condition.")
	if rightParenErr != nil {
		return nil, rightParenErr
	}

	// Then branch of code
	body, bodyErr := p.statement()
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil

}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, nameErr := p.consume(token.IDENTIFIER, "Expected variable name.")
	if nameErr != nil {
		return nil, nameErr
	}

	var initializer Expr
	var err error
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}

	}
	_, semiColonErr := p.consume(token.SEMICOLON, "Expect a ';' after variable declaration.")
	if semiColonErr != nil {
		return nil, semiColonErr
	}

	return &VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, semiColonErr := p.consume(token.SEMICOLON, "Expect ';' after value")
	if semiColonErr != nil {
		return nil, semiColonErr
	}

	return &ExpressionStmt{
		Expression: value,
	}, nil
}

func (p *Parser) function(kind string) (*FunctionStmt, error) {
	name, nameErr := p.consume(token.IDENTIFIER, fmt.Sprintf("Expect %s name\n", kind))
	if nameErr != nil {
		return nil, nameErr
	}
	_, leftParenErr := p.consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name\n", kind))
	if leftParenErr != nil {
		return nil, leftParenErr
	}

	var parameters []Token
	if !p.check(token.RIGHT_PAREN) {
		for true {
			if len(parameters) >= 255 {
				errMessage := fmt.Sprintf("%s Can't have more than 255 parameters.", p.peek())
				log.Println(errMessage)
			}

			identifier, identifierErr := p.consume(token.IDENTIFIER, "Expect a parameter name")
			if identifierErr != nil {
				return nil, identifierErr
			}

			parameters = append(parameters, identifier)
			if !p.match(token.COMMA) {
				break
			}

		}
	}
	_, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if rightParenErr != nil {
		return nil, rightParenErr
	}

	_, leftBraceErr := p.consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", kind))
	if leftBraceErr != nil {
		return nil, leftBraceErr
	}

	body, bodyErr := p.block()
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, stmtErr := p.declaration()
		if stmtErr != nil {
			return nil, stmtErr
		}

		statements = append(statements, stmt)
	}
	_, rightBraceErr := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if rightBraceErr != nil {
		return nil, rightBraceErr
	}

	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, exprErr := p.ternary()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, valueErr := p.assignment()
		if valueErr != nil {
			return nil, valueErr
		}

		varExpr, isVarExpr := expr.(*VariableExpr)
		if isVarExpr {
			name := varExpr.Name
			return &ast.AssignExpr{
				Name:  name,
				Value: value,
			}, nil
		}
		return nil, createParseError(equals, "Invalid assignment target.")
	}
	return expr, exprErr
}

func (p *Parser) ternary() (Expr, error) {
	expr, orErr := p.or()
	if orErr != nil {
		return nil, orErr
	}

	if p.match(token.QUESTION) {
		operator := p.previous()
		second, secondErr := p.equality()
		if secondErr != nil {
			return nil, secondErr
		}

		_, err := p.consume(token.COLON, "Expecting a false value for the ternary")
		if err != nil {
			return nil, err
		}

		third, thirdErr := p.equality()
		if thirdErr != nil {
			return nil, thirdErr
		}

		expr = &TernaryExpr{
			Operator: operator,
			First:    expr,
			Second:   second,
			Third:    third,
		}
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, andErr := p.and()
	if andErr != nil {
		return nil, andErr
	}

	for p.match(token.OR) {
		operator := p.previous()
		rightExpr, rightErr := p.and()
		if rightErr != nil {
			return nil, rightErr
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    rightExpr,
		}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, exprErr := p.equality()
	if exprErr != nil {
		return nil, exprErr
	}

	for p.match(token.AND) {
		operator := p.previous()
		rightExpr, rightErr := p.equality()
		if rightErr != nil {
			return nil, rightErr
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    rightExpr,
		}
		return expr, nil
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, compErr := p.comparison()
	if compErr != nil {
		return nil, compErr
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, rightErr := p.comparison()
		if rightErr != nil {
			return nil, rightErr
		}
		expr = &BinaryExpr{
			Left: expr, Operator: operator, Right: right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, termErr := p.term()
	if termErr != nil {
		return nil, termErr
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, rightErr := p.term()
		if rightErr != nil {
			return nil, rightErr
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, factorErr := p.factor()
	if factorErr != nil {
		return nil, factorErr
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, rightErr := p.factor()
		if rightErr != nil {
			return nil, rightErr
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, unaryErr := p.unary()
	if unaryErr != nil {
		return nil, unaryErr
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return nil, rightErr
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(token.MINUS, token.BANG) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return nil, rightErr
		}

		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}
	result, callErr := p.call()
	if callErr != nil {
		return nil, callErr
	}

	return result, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var args []Expr
	if !p.check(token.RIGHT_PAREN) {
		for {
			expr, exprErr := p.expression()
			if exprErr != nil {
				return nil, exprErr
			}
			if len(args) >= 255 {
				log.Printf("%v Can't have more than 255 arguments.", p.peek())
			}
			args = append(args, expr)
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	rightParen, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if rightParenErr != nil {
		return nil, rightParenErr
	}
	return &CallExpr{
		Callee:    callee,
		Paren:     rightParen,
		Arguments: args,
	}, nil
}

func (p *Parser) call() (Expr, error) {
	expr, exprErr := p.primary()
	if exprErr != nil {
		return nil, exprErr
	}

	for true {
		if p.match(token.LEFT_PAREN) {
			expr, exprErr = p.finishCall(expr)
			if exprErr != nil {
				return nil, exprErr
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(token.FALSE) {
		return &LiteralExpr{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return &LiteralExpr{Value: true}, nil
	}
	if p.match(token.NIL) {
		return &LiteralExpr{Value: nil}, nil
	}

	if p.match(token.STRING, token.NUMBER) {
		return &LiteralExpr{Value: p.previous().Literal}, nil
	}

	if p.match(token.IDENTIFIER) {
		return &ast.VariableExpr{
			Name: p.previous(),
		}, nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, rightParenErr := p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if rightParenErr != nil {
			return nil, rightParenErr
		}

		return &GroupingExpr{
			Expression: expr,
		}, nil
	}
	return nil, createParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	// Add the error to the struct, and we keep trucking. We'll see if this becomes a problem
	consumeError := createParseError(p.peek(), message)
	return Token{}, consumeError
}

func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func createParseError(tokenWithError Token, message string) *GloxError {
	if tokenWithError.TokenType == token.EOF {
		return glox_error.Create(tokenWithError.Line, "at end", message)
	}
	return glox_error.Create(tokenWithError.Line, fmt.Sprintf(" at '%s' of token type '%s'", tokenWithError.Lexeme, tokenWithError.TokenType), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}
		p.advance()
	}

}
