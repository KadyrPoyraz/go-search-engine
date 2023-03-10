package search

import (
	"go-search-engine/data"
	"go-search-engine/lexer"
	"math"
	"sort"
	"strings"
)

type ResultItem struct {
	filePath string
	rank 	 float64
}

type Result []ResultItem

func GetSearchByQuery(stringQuery string, data data.Data) Result {
	query := strings.Split(stringQuery, "")

	var result Result

	for filePath, _ := range data.FileTermFreq {
		rank := float64(0)
		lexer := lexer.Lexer{Content: query}
		for lexer.GetNextToken() {
			term := lexer.Value
			tf := tf(term, data.FileTermFreq[filePath], data.FileTermCount[filePath])
			idf := idf(term, data.FileTermFreq)

			rank += tf * idf
		}

		result = append(result, ResultItem{filePath: filePath, rank: rank})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].rank > result[j].rank
	})
	topOfResults := 10

	if len(result) < topOfResults {
		topOfResults = len(result)
	}
	result = result[:topOfResults]

	return result
}

func getItemOrZero(key string, value map[string]int) int {
	if _, ok := value[key]; ok {
		return value[key]
	} else {
		return 0
	}
}

func tf(t string, d map[string]int, tc int) float64{
	a := float64(getItemOrZero(t, d))
	b := float64(tc)

	return a / b
}

func idf(t string, d map[string]map[string]int) float64 {
	N := float64(len(d))
	M := float64(0)

	for _, tfTable := range d {
		if _, ok := tfTable[t]; ok {
			M += 1
		}
	}

	if M < 1 {
		M = 1
	}

	return math.Log10(N / M)
}