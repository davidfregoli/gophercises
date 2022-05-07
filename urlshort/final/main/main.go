package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davidfregoli/gophercises/urlshort/final/urlshort"
)

func main() {
	var yaml string
	var json string
	yamlPath := flag.String("yaml", "", "The path to the redirects YAML file.")
	jsonPath := flag.String("json", "", "The path to the redirects JSON file.")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/facebook":  "https://facebook.com",
		"/instagram": "https://instagram.com",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	if *yamlPath != "" {
		yamlFile, err := os.Open(*yamlPath)
		if err == nil {
			defer yamlFile.Close()

			byteValue, err := ioutil.ReadAll(yamlFile)
			if err != nil {
				panic(err)
			}
			yaml = string(byteValue)
			fmt.Print(yaml)

			// Build the YAMLHandler using the mapHandler as the
			// fallback
			handler, err = urlshort.YAMLHandler([]byte(yaml), handler)
			if err != nil {
				panic(err)
			}
		} else {
			// Cannot open YAML file
			fmt.Print(err)
		}
	}

	if *jsonPath != "" {
		jsonFile, err := os.Open(*jsonPath)
		if err == nil {
			defer jsonFile.Close()

			byteValue, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				panic(err)
			}
			json = string(byteValue)
			fmt.Print(json)

			// Build the JSONHandler using the previous handler as the
			// fallback
			handler, err = urlshort.JSONHandler([]byte(json), handler)
			if err != nil {
				panic(err)
			}
		} else {
			// Cannot open JSON file
			fmt.Print(err)
		}
	}

	fmt.Println("Starting the server on :80")
	http.ListenAndServe(":80", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
