package search

import (
	"fmt"
	"go-search-engine/data"
	"go-search-engine/lexer"
	"go-search-engine/utils"
	"math"
	"sort"
	"strings"
)

func GetSearchByQuery(stringQuery string, indexFile string) {
	data := utils.GetDataFromCache(indexFile)

	query := strings.Split(stringQuery, "")

	var result []struct{fp string; tf float64}

	for filePath, _ := range data.FileTermFreq {
		rank := float64(0)
		lexer := lexer.Lexer{Content: query}
		for lexer.GetNextToken() {
			term := lexer.Value
			tf := tf(term, data.FileTermFreq[filePath], data.FileTermCount[filePath])
			idf := idf(term, data)

			rank += tf * idf
		}

		tfTable := struct{fp string; tf float64}{
			fp: filePath,
			tf: rank,
		}
		result = append(result, tfTable)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].tf > result[j].tf
	})
	topOfResults := 10

	if len(result) < topOfResults {
		topOfResults = len(result)
	}
	result = result[:topOfResults]

	for i := 0; i < len(result); i++ {
		fmt.Println(result[i])
	}

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

func idf(t string, d data.Data) float64 {
	docs := d.FileTermFreq
	N := float64(len(docs))
	M := float64(0)

	for _, tfTable := range docs {
		if _, ok := tfTable[t]; ok {
			M += 1
		}
	}

	if M < 1 {
		M = 1
	}

	return math.Log10(N / M)
}