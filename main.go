package main

import (
	"flag"
	"fmt"
	"go-search-engine/data"
	"go-search-engine/lexer"
	"go-search-engine/search"
	"go-search-engine/server"
	"go-search-engine/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
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

	serveCmd := flag.NewFlagSet("server", flag.ExitOnError)
	servePort := serveCmd.Int("port", 8000, "Port of running server")

	if len(os.Args) < 2 {
		fmt.Println("Expected \"index\", \"search\" or \"serve\" subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "index":
		fmt.Println("Index in process...")
		err := indexCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}

		indexDirToFile(*indexDirToFilePath, *indexFilePath)

		fmt.Printf("Dir %s has been indexed into file %s\n", *indexDirToFilePath, *indexFilePath)
	case "search":
		fmt.Println("Search in process...")
		err := searchCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		data := utils.GetDataFromCache(*searchIndexFile)
		search.GetSearchByQuery(*searchQuery, data)
	case "serve":
		fmt.Printf("Serving on port %d...", *servePort)
		server.StartServer(*servePort)
	default:
		fmt.Println("Expected \"index\" subcommand")
		os.Exit(1)
	}
}

func indexDirToFile(dirPath string, indexFilePath string) {
	start := time.Now()
	filePaths := getFilePaths(dirPath)
	ch := make(chan map[string][]string)
	data := data.Data{
		FileTermFreq: make(map[string]map[string]int),
		FileTermCount: make(map[string]int),
	}

	for ;len(filePaths) > 0; {
		itemsInBatch := 200
		if len(filePaths) < itemsInBatch {
			itemsInBatch = len(filePaths)
		}

		targetFiles := filePaths[0:itemsInBatch]
		filePaths = filePaths[itemsInBatch:]

		go getFilesData(targetFiles, ch)
		//utils.CacheData(data, indexFilePath)
	}

	for msg := range ch {
		for filePath, terms := range msg {
			for _, term := range terms {
				data.AddFileTermFreqItem(filePath, term)
				data.AddFileTermCount(filePath)
			}
		}
	}
	fmt.Println("Time:", time.Since(start))
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
		if itemExt == ".xhtml" || itemExt == ".txt" {
			pathsList = append(pathsList, itemPath)
		}
		//pathsList = append(pathsList, itemPath)
	}

	return pathsList
}

func getFilesData(paths []string, ch chan map[string][]string) {
	for _, filePath := range paths {
		terms := getDataFromFile(filePath)
		ch <- map[string][]string{filePath: terms}
	}
	close(ch)
}

// TODO: Wrap indexing with goroutines
// TODO: Add saving of indexed data in postgreSQL database