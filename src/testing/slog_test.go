package testing_test

import (
	"log/slog"
	"os"
	"testing"
)

func TestSlog(t *testing.T) {
	// Slog log lines as they are currently printed out in tests.
	logger1 := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	logger1.Info("slog logging which is not indented and not under parent test header")

	/*
		t.Slog() allows:
		- the indentation of slog output to match t.Log() output
		- printing of the output under the correct test header
	*/
	logger2 := t.Slog(nil)
	logger2.Info("t.Slog log which is indented and which is under parent test header")
	t.Error("t.Log in parent test for comparison")

	// Additionally, t.Slog() indents slog output depending on the nesting level of the test.
	t.Run("Subtest", func(t *testing.T) {
		logger3 := t.Slog(nil)
		logger3.Info("t.Slog log which is indented and which is under subtest header")
		t.Error("t.Log in subtest for comparison")

		// Without t.Slog(), the slog log does not take into account the nesting level.
		// This is in addition to the log being printed in the wrong section.
		logger4 := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
		logger4.Info("slog logging which is not indented and not under subtest header")
	})
}
