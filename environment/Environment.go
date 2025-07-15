package environment

import (
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
	"fmt"
)

type Token = token.Token
type GloxError = glox_error.GloxError

type Environment struct {
	values map[string]any
}

func Create() Environment {
	return Environment{
		values: make(map[string]any),
	}

}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {

	value, isPresent := e.values[name.String()]
	if !isPresent {
		return nil, glox_error.Create(name.Line, "", fmt.Sprintf("Undefined variable '%v'.", name.Lexeme))
	}
	return value, nil
}
