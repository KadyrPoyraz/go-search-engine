package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-search-engine/utils"
	"io/ioutil"
	"os"
	"strings"

	"go-search-engine/data"
	"go-search-engine/lexer"
)

const (
	indexJsonFileName = "index.json"
)

var (
	targetDirectory = "./target_files/docs.gl"
)

func main() {
	entry()

	//indexJsonFile, err := ioutil.ReadFile(indexJsonFileName)
	//
	//if err != nil {
	//	readDir(targetDirectory, data)
	//	cacheData(data)
	//}
	//
	//err = json.Unmarshal(indexJsonFile, &data)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println(data)
}

func entry() {
	fmt.Println("Starting the application...")
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	indexDirToFilePath := indexCmd.String("dirPath", "", "The path to the dir to be indexed")
	indexFilePath := indexCmd.String("indexFilePath", "index.json", "Name of file with indexed data")

	if len(os.Args) < 2 {
		fmt.Println("Expected \"index\" subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "index":
		fmt.Println("Index in process...")
		indexCmd.Parse(os.Args[2:])

		isIndexFileExists := checkForIndexedData(*indexFilePath)
		if !isIndexFileExists {
			indexDirToFile(*indexDirToFilePath, *indexfilePath)
		}

		fmt.Printf("Dir %s has been indexed into file %s\n", *indexDirToFilePath, *indexFilePath)
	default:
		fmt.Println("Expected \"index\" subcommand")
		os.Exit(1)
	}
}

func checkForIndexedData(indexFilePath string) bool {
	_, err := os.Stat(indexFilePath)

	return err == nil
}

func cacheData(data data.Data, indexFilePath string) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(indexFilePath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func indexDirToFile(dirPath string, indexFilePath string) {
	data := data.Data{
		FileTermFreq: make(map[string]map[string]int),
		FileTermCount: make(map[string]int),
	}

	collectDirToData(dirPath, data)
	cacheData(data, indexFilePath)
}

func collectDirToData(dirPath string, dataStruct data.Data) {
	items, err := ioutil.ReadDir(dirPath)

	if err != nil {
		panic(err)
	}

	for _, item := range items {
		itemPath := dirPath + "/" + item.Name()

		if item.IsDir() {
			collectDirToData(itemPath, dataStruct)
			continue
		}
		content := utils.GrabTextFromFile(itemPath)

		if len(content) < 1 {
			continue
		}

		collectData(dataStruct, itemPath, content)
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
