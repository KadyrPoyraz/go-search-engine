package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"go-search-engine/utils"
	"go-search-engine/lexer"
	"go-search-engine/data"
)

const (
	indexJsonFileName = "index.json"
)

var (
	targetDirectory = "./target_files/docs.gl"
)

func main() {
	data := data.Data{
		FileTermFreq: make(map[string]map[string]int),
		FileTermCount: make(map[string]int),
	}
	indexJsonFile, err := ioutil.ReadFile(indexJsonFileName)

	if err != nil {
		readDir(targetDirectory, data)
		cacheData(data)
	}

	err = json.Unmarshal(indexJsonFile, &data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
}

func cacheData(data data.Data) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(indexJsonFileName, b, 0644)
	if err != nil {
		panic(err)
	}
}

func collectData(data data.Data, filePath string, content string) {
	data.FileTermFreq[filePath] = make(map[string]int)
	lexer := lexer.Lexer{Content: strings.Split(content, "")}

	for lexer.GetNextToken() {
		token := lexer.Value

		data.AddFileTermFreqItem(filePath, token)
		data.AddFileTermCount(filePath)
	}
}

func readDir(dirPath string, data data.Data) {

	items, err := ioutil.ReadDir(dirPath)

	if err != nil {
		panic(err)
	}

	for _, item := range items {
		itemPath := dirPath + "/" + item.Name()

		if item.IsDir() {
			readDir(itemPath, data)
			continue
		}
		content := utils.GrabTextFromFile(itemPath)

		if len(content) < 1 {
			continue
		}
		collectData(data, itemPath, content)
	}
}
