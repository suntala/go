package testing

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
)

// Slog returns a structured logger that writes to the test case output, like t.Log does.
//
// TODO: Figure out if we want Slog() behaviour to include the following:
//
//	When used with go test -json, any log messages are printed in JSON form.
//	A TestEvent corresponding to a log message has OutputType "slog", and the Output field contains the JSON text.
func (t *T) Slog(opts *slog.HandlerOptions) *slog.Logger {
	return slog.New(t.newOutputHandler(opts))
}

// outputHandler is a Handler that wraps around TextHandler to
// format Records as test logs and writes the result to common.output.
type outputHandler struct {
	th slog.Handler
	bb *bytes.Buffer
	c  *common
}

// newOutputHandler creates an outputHandler using the given options.
// If opts is nil, the default options are used.
func (c *common) newOutputHandler(opts *slog.HandlerOptions) *outputHandler {
	var b bytes.Buffer
	return &outputHandler{
		th: slog.NewTextHandler(&b, opts),
		bb: &b,
		c:  c,
	}
}

// Handle formats its argument Record as a test log.
// #TODO: provide specifics when we finalize the approach.
func (h *outputHandler) Handle(ctx context.Context, r slog.Record) error {
	err := h.th.Handle(ctx, r)
	if err != nil {
		return err
	}

	pcs := make([]uintptr, 10)
	n := runtime.Callers(1, pcs[:])
	f := runtime.CallersFrames(pcs[:n])

	rFrames := runtime.CallersFrames([]uintptr{r.PC})
	rFrame, _ := rFrames.Next()

	var depth int
	for {
		depth++

		frame, more := f.Next()
		if frame.PC == rFrame.PC {
			break
		}
		if !more {
			break
		}
	}

	h.c.logDepth(h.bb.String(), depth)

	return err
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *outputHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.th.Enabled(ctx, level)
}

// WithAttrs returns a new outputHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *outputHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &outputHandler{th: h.th.WithAttrs(attrs)}
}

// WithGroup returns a new outputHandler with the given group appended to
// h's existing groups.
func (h *outputHandler) WithGroup(name string) slog.Handler {
	return &outputHandler{th: h.th.WithGroup(name)}
}
