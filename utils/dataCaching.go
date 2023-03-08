package utils

import (
	"encoding/json"
	"fmt"
	"go-search-engine/data"
	"io/ioutil"
)

func CacheData(data data.Data, indexFilePath string) {
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

func GetDataFromCache(indexFilePath string) data.Data {
	indexJsonFile, err := ioutil.ReadFile(indexFilePath)
	data := data.Data{}

	err = json.Unmarshal(indexJsonFile, &data)
	if err != nil {
		fmt.Println(err)
	}

	return data
}
