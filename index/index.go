package index

import (
	"go-search-engine/data"
	"go-search-engine/lexer"
	"go-search-engine/utils"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

func CreateIndexFileOfDir(dirPath string, indexFilePath string) {
	filePaths := getFilePaths(dirPath)
	ch := make(chan map[string][]string)
	var wg sync.WaitGroup
	data := data.Data{
		FileTermFreq: make(map[string]map[string]int),
		FileTermCount: make(map[string]int),
	}

    totalItems := len(filePaths)
    itemsInBatch := 200

    for i := 0; i < totalItems; i += itemsInBatch {
        limit := itemsInBatch
        if i + itemsInBatch > totalItems {
            limit = totalItems - i 
        }

		targetFiles := filePaths[0:limit]
		filePaths = filePaths[limit:]

		wg.Add(1)
		go getFilesData(targetFiles, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		for filePath, terms := range msg {
			for _, term := range terms {
				data.AddFileTermFreqItem(filePath, term)
				data.AddFileTermCount(filePath)
			}
		}
	}
	utils.CacheData(data, indexFilePath)
}

func getDataFromFile(filePath string) []string {
	content := utils.GrabTextFromFile(filePath)
	var terms []string

	if len(content) < 1 {
		return []string{}
	}

	l := lexer.Lexer{Content: strings.Split(content, "")}

	for l.GetNextToken() {
		term := l.Value

		terms = append(terms, term)
	}

	return terms
}

func getFilePaths(dirPath string) []string {
	var pathsList []string
	items, err := ioutil.ReadDir(dirPath)

	if err != nil {
		panic(err)
	}

	for _, item := range items {
		itemPath := dirPath + "/" + item.Name()

		if item.IsDir() {
			paths := getFilePaths(itemPath)
			pathsList = append(pathsList, paths...)
			continue
		}
		itemExt := path.Ext(itemPath)
		// TODO: Fix condition bellow to more universal thing
		if itemExt == ".xhtml" || itemExt == ".txt" {
			pathsList = append(pathsList, itemPath)
		}
	}

	return pathsList
}

func getFilesData(paths []string, ch chan map[string][]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, filePath := range paths {
		terms := getDataFromFile(filePath)
		ch <- map[string][]string{filePath: terms}
	}
}
