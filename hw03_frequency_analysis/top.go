package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordTop struct {
	word  string
	count int
}

func Top10(text string) []string {
	wordsSlice := strings.Fields(text)
	frequentWordsMap := make(map[string]int)
	wordPopularitySlice := []WordTop{}
	resultWordsSlice := []string{}

	for _, word := range wordsSlice {
		frequentWordsMap[word]++
	}

	for word, count := range frequentWordsMap {
		wordPopularitySlice = append(wordPopularitySlice, WordTop{word: word, count: count})
	}

	sort.Slice(wordPopularitySlice, func(i, j int) bool {
		if wordPopularitySlice[i].count == wordPopularitySlice[j].count {
			return strings.Compare(wordPopularitySlice[i].word, wordPopularitySlice[j].word) == -1
		}

		return wordPopularitySlice[i].count > wordPopularitySlice[j].count
	})

	for index, wordPopularity := range wordPopularitySlice {
		if index > 9 {
			break
		}

		resultWordsSlice = append(resultWordsSlice, wordPopularity.word)
	}

	return resultWordsSlice
}
