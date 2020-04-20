package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/penguingovernor/gophercises/url/internal/urlshort"
)

func main() {

	yamlFile := flag.String("yaml_paths", "./data/paths.yml", "the yaml file that contains the shortened URLs")
	yamlFD, err := os.Open(*yamlFile)
	if err != nil {
		log.Fatalf("unable to open file %s: %v\n", *yamlFile, err)
	}
	defer yamlFD.Close()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yamlFD, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", yamlHandler); err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
