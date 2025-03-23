package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Top10 returns the 10 most frequent words in the input string.
func Top10(str string) []string {
	// Convert the string to lowercase to make the function case-insensitive.
	str = strings.ToLower(str)

	// Split the string into words.
	rawWords := strings.Fields(str)
	freqMap := make(map[string]int)

	// Count the frequency of each word.
	for _, w := range rawWords {
		clean := trimWord(w)
		if clean == "" {
			continue
		}
		freqMap[clean]++
	}

	if len(freqMap) == 0 {
		return []string{}
	}

	// Create a slice of wordCount structs to sort the words by frequency.
	type wordCount struct {
		word  string
		count int
	}

	wordCounts := make([]wordCount, 0, len(freqMap))
	for word, count := range freqMap {
		wordCounts = append(wordCounts, wordCount{word: word, count: count})
	}

	// Sort the words by frequency and alphabetically.
	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	// Collect the top 10 words.
	result := make([]string, 0, 10)
	for i := 0; i < len(wordCounts) && i < 10; i++ {
		result = append(result, wordCounts[i].word)
	}

	return result
}

func trimWord(w string) string {
	start, end := 0, len(w)
	for start < end {
		r, size := utf8.DecodeRuneInString(w[start:])
		if !isEdgePunct(r) {
			break
		}
		start += size
	}
	for start < end {
		r, size := utf8.DecodeLastRuneInString(w[:end])
		if !isEdgePunct(r) {
			break
		}
		end -= size
	}

	res := w[start:end]

	if res == "" {
		if isMultipleDash(w) && len(w) > 1 {
			return w
		}
		return ""
	}

	if res == "-" {
		return ""
	}

	return res
}

func isEdgePunct(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

func isMultipleDash(s string) bool {
	if len(s) < 1 {
		return false
	}
	for _, r := range s {
		if r != '-' {
			return false
		}
	}
	return true
}
