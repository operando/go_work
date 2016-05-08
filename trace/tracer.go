package trace

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事を記録できるオブジェクトを表すインタフェース
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func New(w io.Writer) Tracer {
	// Tracerインタフェースに合致したオブジェクトを受け取るだけで
	// プライベートなtracer型については関知しないので公開していないtracerを返しても問題ない
	return &tracer{out: w}
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}
