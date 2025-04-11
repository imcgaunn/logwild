package logmaker

import (
	"log/slog"
	"strings"
	"testing"
)

func TestGetFakeSentenceWithNWords(t *testing.T) {
	mySentence := GetFakeSentence(30)
	sentenceWords := strings.Split(mySentence, " ")
	slog.Info("words", "sentenceWords", sentenceWords, "len(words)", len(sentenceWords), "len(sentence)", len(mySentence))
	if len(sentenceWords) != 30 {
		t.FailNow()
	}
	lastWord := sentenceWords[len(sentenceWords)-1]
	if lastWord[len(lastWord)-1] != '.' {
		t.FailNow()
	}
}

func TestProfileWhenManyManyFakeSentence(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		GetFakeSentence(64)
	}
}
