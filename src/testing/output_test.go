package testing_test

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestWriterUsingSlogHandler(t *testing.T) {
	// Slog log lines as they are currently printed out in tests.
	logger1 := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	logger1.Info("slog using standard output in parent test")

	/*
		t.Output() allows:
		- the indentation of slog output to match t.Log() output
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

func TestWriteShouldAwaitTrailingNewline(t *testing.T) {
	w := t.Output()

	w.Write([]byte("Hel"))
	w.Write([]byte("lo\nWorld\nInput to log\n\n\nMore logging\nShouldn't be logged"))
	w.Write([]byte("Also shouldn't be logged"))

	t.Error()
}

func TestShouldWriteEntireMultipleLineInputWhenSubtestIsDone(t *testing.T) {
	t.Run("Multiple lines from goroutine", func(t *testing.T) {
		go func() {
			time.Sleep(50 * time.Millisecond)
			w := t.Output()
			w.Write([]byte("First line\nSecond line\nThird line\n"))
		}()
	})
	time.Sleep(1 * time.Second)
	t.Error()
}
