package glox_error

import "fmt"

type GloxError struct {
	line    int
	message string
}

func (e *GloxError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", e.line, "", e.message)
}

func Create(line int, message string) *GloxError {
	return &GloxError{line, message}
}
