package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tagStart      = 60 // Unicode `<`
	tagEnd         = 62 // Unicode `>`
	indexJsonFileName = "index.json"
)

var (
	targetDirectory = "./target_files/docs.gl"
)


type FileTermFreq map[string]int
type FileTermCount map[string]int

type Data struct {
	FileTermFreq map[string]FileTermFreq
	FileTermCount FileTermCount
}

func (d *Data) AddFileTermFreqItem(filePath string, term string) {
	if _, ok := d.FileTermFreq[filePath][term]; ok {
		d.FileTermFreq[filePath][term] += 1
	} else {
		d.FileTermFreq[filePath][term] = 1
	}
}

func (d *Data) AddFileTermCount(filePath string) {
	if _, ok := d.FileTermCount[filePath]; ok {
		d.FileTermCount[filePath] += 1
	} else {
		d.FileTermCount[filePath] = 1
	}
}

type Lexer struct {
	content []string
	value    string
}

func (l *Lexer) trimWhiteSpaces() {
	for ; len(l.content) > 0; {
		isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(l.content[0])
		if isWhitespacePresent {
			l.content = l.content[1:]
		} else {
			break
		}
	}
}

func (l *Lexer) chop(n int) {
	token := l.content[0:n]
	l.content = l.content[n:]
	l.value = strings.ToUpper(strings.Join(token, ""))
}

func (l *Lexer) getNextToken() bool {
	l.trimWhiteSpaces()

	if len(l.content) == 0 {
		return false
	}

	for _, r := range l.content[0] {
		if unicode.IsNumber(r) {
			n := 0
			for ; n < len(l.content) ; {
				_, err := strconv.ParseFloat(l.content[n], 64)
				if err == nil {
					n += 1
				} else {
					break
				}
			}
			l.chop(n)
			return true
		}

		if unicode.IsLetter(r) {
			n := 0
			for ; n < len(l.content) ; {
				//is_alphanumeric := regexp.MustCompile(``).MatchString(l.content[n])
				isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(l.content[n])
				isAlphanumeric := regexp.MustCompile(`[$&+,:;=?@#|'<>.^*()%!-]`).MatchString(l.content[n])
				if !isAlphanumeric && !isWhitespacePresent {
					n += 1
				} else {
					break
				}
			}
			l.chop(n)
			return true
		}
	}

	l.chop(1)

	return true
}

func cacheData(data Data) {
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

func main() {
	//start := time.Now()
	data := Data{
		FileTermFreq: make(map[string]FileTermFreq),
		FileTermCount: make(FileTermCount),
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

func parseXHTMLFile(filePath string) string {
	bytesFile, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	fileContent := string(bytesFile)
	fileContent = stripTags(fileContent)
	fileContent = clearSpaces(fileContent)

	return fileContent
}

// readFile Converts file to string with text from file
func readFile(filePath string) string {
	fileExtension := filepath.Ext(filePath)


	if fileExtension == ".xhtml" {
		content := parseXHTMLFile(filePath)

		return content
	}

	return ""
}

func collectData(data Data, filePath string, content string) {
	data.FileTermFreq[filePath] = make(FileTermFreq)
	lexer := Lexer{content: strings.Split(content, "")}

	for lexer.getNextToken() {
		token := lexer.value

		data.AddFileTermFreqItem(filePath, token)
		data.AddFileTermCount(filePath)
	}
}

func readDir(dirPath string, data Data) {

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
		content := readFile(itemPath)

		if len(content) < 1 {
			continue
		}
		collectData(data, itemPath, content)
	}
}

func stripTags(from string) string {
	targetString := from
	var builder strings.Builder
	builder.Grow(len(targetString) + utf8.UTFMax)

	in := false
	start := 0
	end := 0

	for i, c := range targetString {
		if (i+1) == len(targetString) && end >= start {
			builder.WriteString(targetString[end:])
		}

		if c != tagStart && c != tagEnd {
			continue
		}

		if c == tagStart {
			if !in {
				start = i
			}
			in = true

			builder.WriteString(targetString[end:start])
			continue
		}
		in = false
		end = i + 1
	}
	targetString = builder.String()

	return targetString
}

func clearSpaces(from string) string {
	space := regexp.MustCompile(`\s+`)
	from = space.ReplaceAllString(from , " ")

	return from
}