package parse

import (
	"dsoechting/glox/ast"
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
)

type Expr = ast.Expr
type BinaryExpr = ast.BinaryExpr
type UnaryExpr = ast.UnaryExpr
type LiteralExpr = ast.LiteralExpr
type GroupingExpr = ast.GroupingExpr
type TokenType = token.TokenType
type Token = token.Token

type Parser struct {
	tokens  []token.Token
	current int
}

type ParserError struct {
	line    int
	message string
}

func Create(tokens []token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{
			Left: expr, Operator: operator, Right: right,
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(token.MINUS, token.BANG) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(token.FALSE) {
		return &LiteralExpr{Value: false}
	}
	if p.match(token.TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(token.NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(token.STRING, token.NUMBER) {
		return &LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		return &GroupingExpr{
			Expression: expr,
		}
	}

	return &LiteralExpr{
		Value: "temp",
	}
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

	return Token{}, createParseError(p.peek(), message)
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

func createParseError(tokenWithError Token, message string) ParseError {
	if tokenWithError.TokenType == token.EOF {
		glox_error.Create
	}

}
