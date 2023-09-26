package testing_test

import (
	"log/slog"
	"os"
	"testing"
)

func TestOutput(t *testing.T) {
	// Slog log lines as they are currently printed out in tests.
	logger1 := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	logger1.Info("slog using standard output in parent test")

	/*
		t.Output() allows:
		- the indentation of slog output to match t.Log() output
		- prepending the ouput with callsite information
		- printing of the output under the correct test header
	*/
	logger2 := slog.New(slog.NewTextHandler(t.Output(), nil))
	logger2.Info("slog using t.Output in parent test")
	t.Error("t.Log in parent test")

	// Additionally, t.Output() indents slog output depending on the nesting level of the test.
	t.Run("Subtest", func(t *testing.T) {
		logger3 := slog.New(slog.NewTextHandler(t.Output(), nil))
		logger3.Info("slog using t.Output in subtest")
		t.Error("t.Log in subtest")

		// Without t.Output(), the slog log does not take into account the nesting level.
		// This is in addition to the log being printed in the wrong section.
		logger4 := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
		logger4.Info("slog using standard output in subtest")
	})
}
