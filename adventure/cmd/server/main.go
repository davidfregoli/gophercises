package main

import (
	"encoding/json"
	htmlTpl "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	txtTpl "text/template"
)

const storyFile string = "./story.json"

func main() {
	jsonFile, err := os.Open(storyFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var story Story
	err = json.Unmarshal(byteValue, &story)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rawPath := req.URL.Path
		path := strings.TrimPrefix(rawPath, "/")
		isCli := req.URL.Query().Get("cli") == "true"
		if isCli {
			tplFile, _ := os.Open("./indexCli.tpl")

			byteValue, err := ioutil.ReadAll(tplFile)
			if err != nil {
				panic(err)
			}

			tpl, err := txtTpl.New("Chapter").Parse(string(byteValue))
			if err != nil {
				panic(err)
			}

			chapter, found := story[path]
			if !found {
				http.Redirect(w, req, "/intro?cli=true", http.StatusPermanentRedirect)
				return
			}

			tpl.Execute(w, chapter)
		} else {
			tplFile, _ := os.Open("./index.tpl")

			byteValue, err := ioutil.ReadAll(tplFile)
			if err != nil {
				panic(err)
			}

			tpl, err := htmlTpl.New("Chapter").Parse(string(byteValue))
			if err != nil {
				panic(err)
			}

			chapter, found := story[path]
			if !found {
				http.Redirect(w, req, "/intro", http.StatusPermanentRedirect)
				return
			}

			tpl.Execute(w, chapter)
		}
	})
	http.ListenAndServe(":8080", mux)
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"chapter"`
}
