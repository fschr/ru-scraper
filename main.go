package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"regexp"
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

var word, section string

func main() {
	flag.StringVar(&word, "word", "говорить", "the word to look for")
	flag.StringVar(&section, "section", "nil", "returns a specified section of the output")
	flag.Parse()
	russian := getRussian(getText(baseURL + word))
	var output string
	switch section {
	case "nil":
		output = russian
	default:
		output = getSection(russian, section)
	}
	fmt.Print(output)
}

func getSection(russian, s string) string {
	sectionTitle := "===" + s + "==="
	sectionTitleLen := len(sectionTitle)
	scanner := bufio.NewScanner(strings.NewReader(russian))
	var section string
	searching := true
	for scanner.Scan() {
		line := scanner.Text()
		if searching {
			indx := strings.Index(line, sectionTitle)
			if indx == -1 {
				continue
			}
			searching = false
			if indx+sectionTitleLen != len(line)-1 {
				section += line[indx+sectionTitleLen:]
			}
		} else {
			if strings.HasPrefix(line, "[[") || strings.HasPrefix(line, "===") {
				break
			}
			if line == "" || strings.HasPrefix(line, "{{") {
				continue
			}
			if strings.HasPrefix(line, "# ") {
				line = line[2:]
				line = strings.Replace(line, "[[", "", -1)
				line = strings.Replace(line, "]]", "", -1)
				r := regexp.MustCompile(`\W*{{.*}}\W*`)
				section += strings.TrimSpace(r.ReplaceAllString(line, "")) + "\n"
				continue
			}
			if strings.HasPrefix(line, "#:") || strings.HasPrefix(line, "#*") {
//				line = line[11:]
//				line = strings.Replace(line, "'''", "", -1)
//				r := regexp.MustCompile(`([^\|]*)\|`)
//				section += "\tExample: " + r.FindStringSubmatch(line)[1]
//				r = regexp.MustCompile(`\|t=(.*)\|`)
//				section += " — " + r.FindStringSubmatch(line)[1] + "\n"
				continue
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			section += line + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return section
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
