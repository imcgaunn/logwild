package logmaker

import "github.com/brianvoe/gofakeit/v7"

func GetFakeSentence() string {
	return gofakeit.Sentence(250)
}
