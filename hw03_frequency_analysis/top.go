package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCounter struct {
	word  string
	count int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	uniqueWords := make([]wordCounter, 0, len(counts))
	for word, count := range counts {
		uniqueWords = append(uniqueWords, wordCounter{word, count})
	}
	sort.Slice(uniqueWords, func(i, j int) bool {
		switch {
		case uniqueWords[i].count > uniqueWords[j].count:
			return true
		case uniqueWords[i].count < uniqueWords[j].count:
			return false
		default: // uniqueWords[i].count == uniqueWords[j].count
			return uniqueWords[i].word < uniqueWords[j].word
		}
	})
	leng := len(uniqueWords)
	if leng > 10 {
		leng = 10
	}
	top := make([]string, 0, leng)
	for i := 0; i < leng; i++ {
		top = append(top, uniqueWords[i].word)
	}
	return top
}
