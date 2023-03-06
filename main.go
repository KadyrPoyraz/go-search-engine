package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Data struct {
	fileTermFreq map[string]map[string]int
	fileTermCount map[string]int
}

type DirContent map[string]map[string]int

type Lexer struct {
	content []string
	value   []string
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
	l.value = token
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
				is_alphanumeric := regexp.MustCompile(`[$&+,:;=?@#|'<>.^*()%!-]`).MatchString(l.content[n])
				if !is_alphanumeric && !isWhitespacePresent {
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

const (
	tagStart = 60 // Unicode `<`
	tagEnd   = 62 // Unicode `>`
)

var (
	targetDirectory = "./target_files/docs.gl"
)

func main() {
	start := time.Now()
	allDocuments := make(DirContent)
	readDir(targetDirectory, allDocuments)

	//fmt.Println(allDocuments)
	fmt.Println(time.Since(start))
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

func readDir(dirPath string, dirContent DirContent) {
	items, err := ioutil.ReadDir(dirPath)

	if err != nil {
		panic(err)
	}

	for _, item := range items {
		itemPath := dirPath + "/" + item.Name()

		if item.IsDir() {
			readDir(itemPath, dirContent)
			continue
		}

		filePath := itemPath
		fileExtension := filepath.Ext(filePath)

		if fileExtension == ".xhtml" {
			//fmt.Println("Reading file:", filePath)
			if _, ok := dirContent[filePath]; !ok {
				dirContent[filePath] = make(map[string]int)
			}

			content := parseXHTMLFile(filePath)
			lexer := Lexer{content: strings.Split(content, "")}

			for lexer.getNextToken() {
				token := lexer.value
				result := strings.ToUpper(strings.Join(token, ""))

				if _, ok := dirContent[filePath][result]; ok {
					dirContent[filePath][result] += 1
				} else {
					dirContent[filePath][result] = 1
				}

			}
		}
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
