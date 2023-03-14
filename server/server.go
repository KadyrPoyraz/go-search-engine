package server

import (
	"log"
	"net/http"
	"strconv"
)

func StartServer(port int) {
	http.HandleFunc("/search/", QuerySearchHandler)

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}