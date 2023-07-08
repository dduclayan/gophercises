// Exercise - https://github.com/gophercises/urlshort
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"urlshortJSON"
)

var (
	file = flag.String("file", "", "path to yaml file")
)

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortJSON.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// NOTE: no tabs allowed in yaml. spacing is very important.

	jsonData, err := generateJSONInput()
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshortJSON.JSONHandler(jsonData, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func generateJSONInput() ([]byte, error) {
	if *file != "" {
		fmt.Printf("reading file %v\n", *file)
		data, err := readFile(*file)
		if err != nil {
			panic(err)
		}
		return data, nil
	}
	jsonByte := []byte(`[
{"path": "/gh", "url": "https://github.com"},
{"path": "/gm", "url": "https://gmail.com"}
]`)
	return jsonByte, nil
}

func readFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
