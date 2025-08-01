package token

import "fmt"

// Just run go generate in this directory to update the enum string
//
//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
type TokenType int

const (
	// Single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	QUESTION
	COLON

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	// TODO can I type param this?
	Literal any
	Line    int
}

func Create(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.TokenType, t.Lexeme, t.Literal)
}
