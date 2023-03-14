package main

import (
	"flag"
	"fmt"
	"go-search-engine/server"
	"go-search-engine/utils"
	"io/ioutil"
	"os"
	"strings"
	"time"

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

		l := lexer.Lexer{Content: strings.Split(content, "")}

		for l.GetNextToken() {
			term := l.Value

			dataStruct.AddFileTermFreqItem(itemPath, term)
			dataStruct.AddFileTermCount(itemPath)
		}
	}
}

// TODO: Add saving of indexed data in postgreSQL database
// TODO: Add stemming for the search