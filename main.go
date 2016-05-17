package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseURL = `https://en.wiktionary.org/w/api.php?` +
		`action=query` +
		`&format=json` +
		`&prop=` +
		`&export=1` +
		`&exportnowrap=1` +
		`&redirects=1` +
		`&utf8=1` +
		`&titles=`
)

var word string

func main() {
	flag.StringVar(&word, "word", "говорить", "the word to look for")
	flag.Parse()
	fmt.Println(getRussian(getText(baseURL + word)))
}

func getText(URL string) string {
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// TODO: Create a struct and unmarshal the XML into a struct
	scanner := bufio.NewScanner(resp.Body)
	var text string
	searching := true
	for scanner.Scan() {
		line := scanner.Text()
		if searching {
			if strings.Contains(line, "<text") {
				indx := strings.Index(line, ">")
				if indx == -1 {
					continue
				}
				searching = false
				if indx+1 != len(line)-1 {
					text += line[indx+1:]
				}
			}
		} else {
			if strings.Contains(line, "</text>") {
				indx := strings.Index(line, "<")
				text += line[:indx]
				break
			}
			text += line + string('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return text
}

func getRussian(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var russian string
	searching := true
	for scanner.Scan() {
		line := scanner.Text()
		if searching {
			indx := strings.Index(line, "==Russian==")
			if indx == -1 {
				continue
			}
			searching = false
			if indx+len("==Russian==") != len(line)-1 {
				russian += line[indx+len("==Russian=="):]
			}
		} else {
			otherLang := func(line string) bool {
				return strings.HasPrefix(line, "==") && !strings.HasPrefix(line, "===")
			}
			if strings.HasPrefix(line, "[[") || otherLang(line) {
				break
			}
			russian += line + string('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return russian
}
