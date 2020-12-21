package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

// Helper structure to accumulate information
// about certain word.
type WordCount struct {
	// Countable lexeme
	Word string
	// Number of Word occurrences
	Count int
}

func prepareText(text string) (string, bool) {
	if text == "" {
		return text, false
	}

	// Make everything lowercase
	text = strings.ToLower(text)

	// remove punctuation
	punct := regexp.MustCompile("[[:punct:]]+")
	text = punct.ReplaceAllString(text, "")

	// remove extra whitespaces
	re := regexp.MustCompile("[[:space:]]+")
	text = re.ReplaceAllString(text, " ")

	return text, true
}

func Top10(text string) []string {
	result := make([]string, 0, 10)

	text, ok := prepareText(text)
	if !ok {
		return result
	}

	counter := make(map[string]*WordCount)

	for _, word := range strings.Split(text, " ") {
		count, ok := counter[word]
		if !ok {
			counter[word] = &WordCount{Word: word, Count: 1}
		} else {
			count.Count++
		}
	}

	slice := make([]*WordCount, 0, len(counter))
	for _, count := range counter {
		slice = append(slice, count)
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count
	})

	for i := 0; i < 10 && i < len(slice); i++ {
		result = append(result, slice[i].Word)
	}

	return result
}
