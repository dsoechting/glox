package environment

import (
	glox_error "dsoechting/glox/error"
	"dsoechting/glox/token"
	"fmt"
	"strings"
)

type Token = token.Token
type GloxError = glox_error.GloxError

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func Create() Environment {
	return Environment{
		values:    make(map[string]any),
		enclosing: nil,
	}
}

func CreateWithEnclosing(enclosing Environment) Environment {
	return Environment{
		values:    make(map[string]any),
		enclosing: &enclosing,
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	value, isPresent := e.values[name.Lexeme]
	if isPresent {
		return value, nil
	}

	if e.enclosing != nil {

		enclosedValue, enclosedError := e.enclosing.Get(name)
		if enclosedError != nil {
			return nil, enclosedError
		}
		return enclosedValue, nil
	}
	return nil, glox_error.Create(name.Line, "", fmt.Sprintf("Undefined variable '%v'.\n", name.Lexeme))
}

func (e *Environment) Assign(name Token, value any) error {
	_, isPresent := e.values[name.Lexeme]
	if isPresent {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		encloseErr := e.enclosing.Assign(name, value)
		if encloseErr != nil {
			return encloseErr
		} else {
			return nil
		}

	}
	return glox_error.Create(name.Line, "", fmt.Sprintf("Undefined variable '%v'.", name.Lexeme))
}

func (e *Environment) String() string {
	var sb strings.Builder
	sb.WriteString("----------")
	sb.WriteString("Current Env:\n")
	helper(&sb, e)
	sb.WriteString("----------")
	return sb.String()
}

func helper(sb *strings.Builder, env *Environment) {
	sb.WriteString(fmt.Sprintln(env.values))
	if env.enclosing != nil {
		sb.WriteString("Found Enclosing\n")
		helper(sb, env.enclosing)
	}
}
