package ports

import "io"

type Logger interface {
	SetDebugLevel(dl bool)
	SetWriter(w io.Writer)
	Error(v ...any)
	Info(v ...any)
	Debug(v ...any)
	Panic(v ...any)
}
