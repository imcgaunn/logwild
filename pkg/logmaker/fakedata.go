package logmaker

import "github.com/brianvoe/gofakeit/v7"

func GetFakeSentence(numWords int) string {
	// if you try to generate too many random words
	// leads to bottlenecks
	return gofakeit.Sentence(numWords)
}
