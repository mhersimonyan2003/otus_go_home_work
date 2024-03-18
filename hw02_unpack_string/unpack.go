package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var resultString strings.Builder
	var currentSymbol string
	currentNumber := -1
	runesCount := utf8.RuneCountInString(str)

	for index, symbol := range str {
		strSymbol := string(symbol)

		if strInt, err := strconv.Atoi(strSymbol); err == nil {
			if currentSymbol == "" || currentNumber != -1 {
				return "", ErrInvalidString
			}

			currentNumber = strInt
			resultString.WriteString(strings.Repeat(currentSymbol, currentNumber))

			currentNumber = -1
			currentSymbol = ""

			continue
		}

		resultString.WriteString(currentSymbol)
		currentSymbol = strSymbol

		if index == runesCount-1 {
			resultString.WriteString(strSymbol)
		}
	}

	return resultString.String(), nil
}
