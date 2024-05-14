package logmaker

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestCanMakeLogMaker(t *testing.T) {
	levelVar := &slog.LevelVar{}
	hdl := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: levelVar})
	mkr := NewLogMaker(WithLogger(slog.New(hdl)))
	if mkr.BurstDuration != time.Duration(5*time.Second) {
		t.FailNow()
	}
}

func TestThatLogMakerLogsToConfiguredLogger(t *testing.T) {
	// get a temporary file and setup logger to write to it
	f, err := os.CreateTemp("", "logmakrtest")
	if err != nil {
		t.Errorf("something bad happened trying to open temp file %s\n", err)
	}
	defer os.Remove(f.Name())

	donech := make(chan int)
	hdl := slog.NewJSONHandler(f, nil)
	mkr := NewLogMaker(
		WithLogger(slog.New(hdl)),
		WithPerSecondRate(5000),
		WithBurstDuration(time.Second))
	go func() {
		err := mkr.StartWriting(donech)
		if err != nil {
			panic(err)
		}
	}()
	wroteMessages := <-donech
	t.Logf("finished test and wrote %d\n", wroteMessages)
}
