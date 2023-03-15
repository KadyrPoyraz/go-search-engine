package search

import (
	"go-search-engine/data"
	"go-search-engine/lexer"
	"math"
	"sort"
	"strings"
)

type ResultItem struct {
	FilePath string
	Rank 	 float64
}

type Result []ResultItem

func GetSearchByQuery(stringQuery string, data data.Data) Result {
	query := strings.Split(stringQuery, "")
	ch := make(chan ResultItem)
	quit := make(chan bool)

	var result Result
	dataFilePaths := make([]string, len(data.FileTermFreq))

	i := 0
	for filePath, _ := range data.FileTermFreq {
		dataFilePaths[i] = filePath
		i += 1
	}

	for ;len(dataFilePaths) > 0; {
		itemsInBatch := 500
		if len(dataFilePaths) < itemsInBatch {
			itemsInBatch = len(dataFilePaths)
		}

		targetFiles := dataFilePaths[0:itemsInBatch]
		dataFilePaths = dataFilePaths[itemsInBatch:]

		go func(closeChanel bool){

			for _, filePath := range targetFiles {
				rank := float64(0)
				lexer := lexer.Lexer{Content: query}
				for lexer.GetNextToken() {
					term := lexer.Value
					tf := tf(term, data.FileTermFreq[filePath], data.FileTermCount[filePath])
					idf := idf(term, data.FileTermFreq)

					rank += tf * idf
				}

				ch <- ResultItem{FilePath: filePath, Rank: rank}
			}
			if closeChanel {
				quit <- true
			}
		}(len(dataFilePaths) == 0)
	}

	func(){
		for {
			select {
			case res := <-ch:
				//fmt.Println(res)
				result = append(result, res)
			case <-quit:
				return
			}
		}
	}()

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rank > result[j].Rank
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