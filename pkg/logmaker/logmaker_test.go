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
	if mkr.PerSecondRate != 1000 {
		t.FailNow()
	}
	if mkr.PerMessageSize != 2*1024 {
		t.FailNow()
	}
}

func TestThatLogMakerLogsToConfiguredLogger(t *testing.T) {
	// get a temporary file and setup logger to write to it
	f, err := os.CreateTemp("", "logmakrtest")
	if err != nil {
		t.Errorf("something bad happened trying to open temp file %s\n", err)
	}
	// if we did open the file, clean it up at the end of the test
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err) // really too late to do anything at this point :/
		}
	}(f.Name())

	donech := make(chan int)
	hdl := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})
	mkr := NewLogMaker(WithLogger(slog.New(hdl)),
		WithPerSecondRate(5000),
		WithPerMessageSize(1024),
		WithBurstDuration(5*time.Second))
	go func() {
		err := mkr.StartWriting(donech)
		if err != nil {
			panic(err)
		}
	}()
	wroteMessages := <-donech
	t.Logf("finished test and wrote %d\n", wroteMessages)
}
