package search

import (
	"go-search-engine/data"
	"math"
)

func getItemOrZero(key string, value map[string]int) int {
	if _, ok := value[key]; ok {
		return value[key]
	} else {
		return 0
	}
}

func Tf(t string, d map[string]int, tc int) float64{
	a := float64(getItemOrZero(t, d))
	b := float64(tc)

	return a / b
}

func Idf(t string, d data.Data) float64 {
	docs := d.FileTermFreq
	N := float64(len(docs))
	M := float64(1)

	for _, tfTable := range docs {
		for term, _ := range tfTable {
			if t == term {
				M += 1
			}
		}
	}

	return math.Log10(N / M)
}