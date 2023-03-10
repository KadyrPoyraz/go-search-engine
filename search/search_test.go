package search

import (
	"go-search-engine/data"
	"go-search-engine/lexer"
	"reflect"
	"strings"
	"testing"
)

type Test struct {
	name   string
	query  string
	files  []map[string]string
	result Result
}

func TestGetSearchByQuery(t *testing.T) {
	tests := []Test{
		{
			name: "Test search dog in documents",
			query: "where is the dog?",
			files: []map[string]string{
				{"fileOne": "Here is the dog in this document"},
				{"fileTwo": "Here is the fox in this document"},
			},
			result: Result{
				ResultItem{filePath: "fileOne", rank: 0.043004285094854454},
				ResultItem{filePath: "fileTwo", rank: 0},
			},
		},
		{
			name: "Test search fox in documents",
			query: "where is the fox?",
			files: []map[string]string{
				{"fileOne": "Here is the dog in this document"},
				{"fileTwo": "Here is the fox in this document"},
			},
			result: Result{
				ResultItem{filePath: "fileTwo", rank: 0.043004285094854454},
				ResultItem{filePath: "fileOne", rank: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := data.Data{
				FileTermFreq: make(map[string]map[string]int),
				FileTermCount: make(map[string]int),
			}
			for _, file := range test.files {
				for fileName, fileContent := range file {
					l := lexer.Lexer{Content: strings.Split(fileContent, "")}
					l.PutContentToData(d, fileName)
				}
			}

			result := GetSearchByQuery(test.query, d)

			if !reflect.DeepEqual(test.result, result) {
				t.Errorf( "FAILED: expected %v, got %v\n", test.result, result)
			}
		})
	}
}
