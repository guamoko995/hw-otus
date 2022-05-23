package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// wordCounter хранит слово и его частотность.
type wordCounter struct {
	word  string
	count int
}

var (
	en = `[A-Za-z][A-Za-z-]*` // Регулярное выражение английского слова
	ru = `[А-Яа-я][А-Яа-я-]*` // Регулярное выражение русского слова
	d  = `[0-9]+`             // Регулярное выражение числа

	// Скомпилированное регулярное выражение, соответствующее
	// английскому слову, русскому слову или числу.
	re = regexp.MustCompile(en + "|" + ru + "|" + d)
)

// Top10 принимает на вход строку с текстом и возвращает слайс
// с 10-ю наиболее часто встречаемыми в тексте словами.
func Top10(text string) []string {
	// Получение массива всех слов из текста.
	words := re.FindAllString(text, -1)

	// Подсчет количеств каждого уникального слова в тексте.
	counts := make(map[string]int)
	for _, word := range words {
		counts[strings.ToLower(word)]++
	}

	// Количество уникальных слов.
	leng := len(counts)

	// Массив подсчитанных уникальных слов
	uniqueWords := make([]wordCounter, 0, leng)
	for word, count := range counts {
		uniqueWords = append(uniqueWords, wordCounter{word, count})
	}

	// Сортировка моссива по частотности (приоритет) и по алфавиту
	// (в случае одинаковой частотности).
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

	// Количество слов в топе
	if leng > 10 {
		leng = 10
	}

	// Создание топа самых встречаемых слов в тексте
	// (не более 10).
	top := make([]string, 0, leng)
	for i := 0; i < leng; i++ {
		top = append(top, uniqueWords[i].word)
	}

	return top
}
