package interpret

type GloxCallabe interface {
	call(interpreter Interpreter, arguments []any)
}
