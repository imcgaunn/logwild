package logmaker

import (
	"log/slog"
	"testing"
)

func TestWhatDoesFakeDataLookLike(t *testing.T) {
	mySentence := GetFakeSentence()
	slog.Info("got sentence", "mySentence", mySentence)
}
