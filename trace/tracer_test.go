package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("return nil")
	} else {
		tracer.Trace("hogehoge")
		if buf.String() != "hogehoge\n" {
			t.Error("'%s' mismatch", buf.String())
		}
	}
}
