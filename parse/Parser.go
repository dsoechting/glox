package parse

import (
	"dsoechting/glox/ast"
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
	"fmt"
)

type Expr = ast.Expr
type Stmt = ast.Stmt
type VarStmt = ast.VarStmt
type TernaryExpr = ast.TernaryExpr
type BinaryExpr = ast.BinaryExpr
type UnaryExpr = ast.UnaryExpr
type LiteralExpr = ast.LiteralExpr
type GroupingExpr = ast.GroupingExpr
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
		// stmt, stmtErr := p.statement()
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
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
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

	return &ast.ExpressionStmt{
		Expression: value,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.ternary()
}

func (p *Parser) ternary() (Expr, error) {
	expr, eqErr := p.equality()
	if eqErr != nil {
		return nil, eqErr
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
	result, primaryErr := p.primary()
	if primaryErr != nil {
		return nil, primaryErr
	}

	return result, nil
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
