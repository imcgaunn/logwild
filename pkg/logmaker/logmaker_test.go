package logmaker

import (
	"log/slog"
	"os"
	"strings"
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
	if mkr.PerMessageSize != 48 {
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

func TestThatLogMakerAppendsToExistingFile(t *testing.T) {
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

	// write some initial content to the file
	initialContent := "initial content"
	_, err = f.WriteString(initialContent)
	if err != nil {
		t.Errorf("something bad happened trying to write to temp file %s\n", err)
	}

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
	t.Logf("finished test and wrote %d log messages\n", wroteMessages)

	// read the file and verify that the original content is still present
	// and that the new logs have been appended to the end of the file
	_, err = f.Seek(0, 0)
	if err != nil {
		t.Errorf("something bad happened trying to seek to the beginning of the temp file %s\n", err)
	}
	fileContent, err := os.ReadFile(f.Name())
	if err != nil {
		t.Errorf("something bad happened trying to read the temp file %s\n", err)
	}
	if !strings.HasPrefix(string(fileContent), initialContent) {
		t.Errorf("expected file to start with '%s', but it didn't", initialContent)
	}
	if len(string(fileContent)) <= len(initialContent) {
		t.Errorf("expected file to contain more than just the initial content, but it didn't")
	}
}
