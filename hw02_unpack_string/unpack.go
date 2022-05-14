package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var b strings.Builder
	var shielding bool
	var lastSimbol string

	for _, r := range input {
		isBackslash := r == '\\'

		// Экранирование следующего символа.
		if isBackslash && !shielding {
			shielding = true
			continue
		}

		isDigit := unicode.IsDigit(r)

		// Проверка: экранировать можно только цифры и `\`.
		if shielding && !isDigit && !isBackslash {
			return "", ErrInvalidString
		}

		s := string([]rune{r})

		// Распаковка предыдущего символа (lastSimbol) входной строки (в зависимости от текущего).
		if !shielding && isDigit {
			// lastSimbol=="" означает, что одно из следующих утверждений верно:
			// 		1) текущий символ первый (начало цикла);
			// 		2) предыдущий символ является не экранированной цифрой.
			//
			// Не экранированная цифра не может следовать за не экранированной
			// цифрой и строка не может начинаться с не экранированной цифры.
			if lastSimbol == "" {
				return "", ErrInvalidString
			}

			n, err := strconv.Atoi(s)
			if err != nil {
				// Этот код никогда не должен быть выполнен (предохранитель).
				panic(fmt.Sprintf("попытка распарсить \"%[1]v\" (тип %[1]T) в тип int", s))
			}

			b.WriteString(strings.Repeat(lastSimbol, n))

			lastSimbol = "" // Последний символ является не экранированной цифрой.
		} else {
			b.WriteString(lastSimbol)
			lastSimbol = s
		}
		shielding = false
	}

	// Строка не может заканчиваться не экранированным '/'.
	if shielding {
		return "", ErrInvalidString
	}

	// На каждом шаге цикла добавляется от 0 до 9 символов, полученных из input на
	// одном из предыдущих шагов, что влечет необходимость добавления последнего символа к
	// выходному буферу.
	b.WriteString(lastSimbol)

	return b.String(), nil
}
