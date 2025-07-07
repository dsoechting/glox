package glox_error

import "fmt"

type GloxError struct {
	line    int
	where   string
	message string
}

func (e *GloxError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", e.line, e.where, e.message)
}

func Create(line int, where string, message string) *GloxError {
	return &GloxError{line, where, message}
}
