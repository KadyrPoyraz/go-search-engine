package server

import (
	"encoding/json"
	"go-search-engine/search"
	"go-search-engine/utils"
	"net/http"
	"strings"
)

// TODO: Improve structure of handlers, maybe with using DDD, just for practice

func QuerySearchHandler(rw http.ResponseWriter, r *http.Request) {
	url := r.URL
	query := url.Query()
	rw.Header().Set("Content-Type", "application/json")

	if _, ok := query["query"]; !ok {
		rw.Write([]byte("Please enter \"query\" query parameter to search"))
		return
	}

	searchQuery := strings.Join(query["query"], "")

	// TODO: implement grabbing indexFilePath from env variables
	data := utils.GetDataFromCache("index.json")
	searchResult := search.GetSearchByQuery(searchQuery, data)

	result, err := json.Marshal(searchResult)

	if err != nil {
		panic(err)
	}

	rw.Write(result)
	return
}