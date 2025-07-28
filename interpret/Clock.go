package interpret

import "time"

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return time.Now().UTC().UnixMilli() / 1000.0, nil
}

func (c *Clock) String() string {
	return "<native fn>"
}
