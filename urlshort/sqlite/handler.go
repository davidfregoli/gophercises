package urlshort

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type redirect struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(db map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		from := req.URL.String()
		to, match := db[from]
		if !match {
			fallback.ServeHTTP(res, req)
			return
		}
		http.Redirect(res, req, to, http.StatusTemporaryRedirect)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var data []redirect

	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}
	var pathsToUrls = make(map[string]string)
	for _, redirect := range data {
		pathsToUrls[redirect.Path] = redirect.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//    [{
//			"path": "/some-path",
//      "url": "https://www.some-url.com/demo"
//		}]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var data []redirect

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	var pathsToUrls = make(map[string]string)
	for _, redirect := range data {
		pathsToUrls[redirect.Path] = redirect.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func SQLHandler(fallback http.Handler) (http.HandlerFunc, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	query, err := load("queries/create")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return func(res http.ResponseWriter, req *http.Request) {
		var to string
		from := req.URL.String()
		row := db.QueryRow("SELECT `to` FROM redirects WHERE `from` = ?", from)
		if err := row.Scan(&to); err != nil {
			if err == sql.ErrNoRows {
				fallback.ServeHTTP(res, req)
				return
			}
			panic(err)
		}

		http.Redirect(res, req, to, http.StatusTemporaryRedirect)
	}, nil
}

func load(path string) (string, error) {
	file, err := os.Open(path + ".sql")
	if err != nil {
		return "", err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(byteValue), nil
}
