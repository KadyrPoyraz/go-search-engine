package server

import (
	"net/http"
	"strings"
)

func QuerySearchHandler(rw http.ResponseWriter, r *http.Request) {
	url := r.URL
	query := url.Query()

	if _, ok := query["query"]; !ok {
		rw.Write([]byte("Please enter \"query\" query parameter to search"))
	}

	searchQuery := strings.Join(query["query"], "")
	rw.Write([]byte(searchQuery))
}