package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"urlshortDB"
)

func main() {
	flag.Parse()
	mux := defaultMux()

	jsonByte := []byte(`[
{"path": "/gh", "url": "https://github.com"},
{"path": "/gm", "url": "https://gmail.com"}
]`)
	data, err := urlshortDB.ParseJSON(jsonByte)
	if err != nil {
		fmt.Printf("urlshortDB.ParseJSON(jsonByte): %v", err)
	}
	for k, v := range data {
		if err := urlshortDB.AddEntry(k, v); err != nil {
			fmt.Printf("urlshortDB.AddEntry(%v,%v): %v", k, v, err)
			continue
		}
	}

	// Build the DBhandler using the mux as the fallback
	_, err = urlshortDB.SetupDB()
	if err != nil {
		fmt.Printf("could not setup db: %v", err)
		os.Exit(1)
	}

	dbHandler := urlshortDB.DBHandler(mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
