package data

import (
	"reflect"
	"testing"
)

type Test struct {
	name string
	input []map[string][]string
	result Data
}

func TestData_AddFileTermCount(t *testing.T) {
	tests := []Test{
		{
			name: "Test if Data adds file term count and frequency correctly",
			input: []map[string][]string{
				{
					"fileOne": {
						"test1", "test2", "test3", "test2",
					},
				},
				{
					"fileTwo": {
						"test1", "test2", "test3", "test1",
					},

				},
			},
			result: Data{
				FileTermFreq: map[string]map[string]int{
					"fileOne": {
						"test1": 1,
						"test2": 2,
						"test3": 1,
					},
					"fileTwo": {
						"test1": 2,
						"test2": 1,
						"test3": 1,
					},
				},
				FileTermCount: map[string]int{
					"fileOne": 4,
					"fileTwo": 4,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := Data{
				FileTermFreq: make(map[string]map[string]int),
				FileTermCount: make(map[string]int),
			}
			for _, file := range test.input {
				for fileName, terms := range file {
					for _, term := range terms {
						data.AddFileTermFreqItem(fileName, term)
						data.AddFileTermCount(fileName)
					}
				}
			}

			if !reflect.DeepEqual(test.result, data) {
				t.Errorf("FAILED: expected %v, got %v\n", test.result, data)
			}
		})
	}
}