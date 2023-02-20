package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Info(a ...any)
	Error(a ...any)
	Fatal(a ...any)
}

func NewLogger(w io.Writer) Logger {
	return &logger{
		writer: w,
	}
}

type logger struct {
	writer io.Writer
}

func (l logger) Info(a ...any) {
	_, _ = fmt.Fprintln(l.writer, a...)
}

func (l logger) Error(a ...any) {
	_, _ = fmt.Fprintln(l.writer, "Error:", a)
}

func (l logger) Fatal(a ...any) {
	_, _ = fmt.Fprintln(l.writer, "Fatal:", a)
	os.Exit(1)
}
