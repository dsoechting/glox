package interpret

type GloxCallable interface {
	call(interpreter *Interpreter, arguments []any) (any, error)
	arity() int
}
