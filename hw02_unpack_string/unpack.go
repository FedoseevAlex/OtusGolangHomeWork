package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func repeat(symbol rune, times int) []rune {
	result := make([]rune, times)

	for i := 0; i < times; i++ {
		result[i] = symbol
	}

	return result
}

func Unpack(packed string) (string, error) {
	runes := []rune(packed)
	result := []rune{}

	for current, next := 0, 1; current < len(runes); current, next = current+1, next+1 {
		if unicode.IsDigit(runes[current]) {
			// If we are here either string starts with digit or
			// there are two consecutive digits
			return "", ErrInvalidString
		}

		if runes[current] == '\\' {
			if next == len(runes) {
				// Return error if backslash is the last symbol in string
				return "", ErrInvalidString
			}

			if runes[next] != '\\' && !unicode.IsDigit(runes[next]) {
				// Return error if backslash precede something that is not
				// a backslash or digit
				return "", ErrInvalidString
			}

			// If all checks passed jump to escaped symbol
			current, next = current+1, next+1
		}

		if next == len(runes) {
			// If current is the last symbol just add it to result
			result = append(result, runes[current])
			continue
		}

		if unicode.IsDigit(runes[next]) {
			reps, err := strconv.Atoi(string(runes[next]))
			if err != nil {
				// Something unusual happens
				return "", err
			}

			// Add expanded symbol
			result = append(result, repeat(runes[current], reps)...)

			// Skip symbol multiplier
			current, next = current+1, next+1
		} else {
			// No multiplier. Just add the current symbol
			result = append(result, runes[current])
		}
	}
	return string(result), nil
}
