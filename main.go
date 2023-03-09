package main

import (
	"flag"
	"fmt"
	"go-search-engine/utils"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"go-search-engine/data"
	"go-search-engine/lexer"
	"go-search-engine/search"
)
func main() {
	entry()
}

func entry() {
	fmt.Println("Starting the application...")
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	indexDirToFilePath := indexCmd.String("dirPath", "", "The path to the dir to be indexed")
	indexFilePath := indexCmd.String("indexFilePath", "index.json", "Name of file with indexed data")

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchIndexFile := searchCmd.String("indexFile", "", "Path to file to search for")
	searchQuery := searchCmd.String("query", "", "Query within the index file")

	if len(os.Args) < 2 {
		fmt.Println("Expected \"index\" subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "index":
		fmt.Println("Index in process...")
		err := indexCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}

		isIndexFileExists := checkForIndexedData(*indexFilePath)
		if !isIndexFileExists {
			indexDirToFile(*indexDirToFilePath, *indexFilePath)
		}

		fmt.Printf("Dir %s has been indexed into file %s\n", *indexDirToFilePath, *indexFilePath)
	case "search":
		fmt.Println("Search in process...")
		err := searchCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}

		data := utils.GetDataFromCache(*searchIndexFile)

		query := strings.Split(*searchQuery, "")

		var result []struct{fp string; tf float64}

		for filePath, _ := range data.FileTermFreq {
			rank := float64(0)
			lexer := lexer.Lexer{Content: query}
			for lexer.GetNextToken() {
				term := lexer.Value
				tf := search.Tf(term, data.FileTermFreq[filePath], data.FileTermCount[filePath])
				idf := search.Idf(term, data)

				rank = tf * idf
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
		result = result[:10]

		for i := 0; i < len(result); i++ {
			fmt.Println(result[i])
		}

		fmt.Printf("Search in %s file by query %s", *searchIndexFile, *searchQuery)
	default:
		fmt.Println("Expected \"index\" subcommand")
		os.Exit(1)
	}
}

func checkForIndexedData(indexFilePath string) bool {
	_, err := os.Stat(indexFilePath)

	return err == nil
}


func indexDirToFile(dirPath string, indexFilePath string) {
	data := data.Data{
		FileTermFreq: make(map[string]map[string]int),
		FileTermCount: make(map[string]int),
	}

	collectDirToData(dirPath, data)
	utils.CacheData(data, indexFilePath)
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
