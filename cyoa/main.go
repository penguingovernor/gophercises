package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/penguingovernor/gophercises/cyoa/internal/adventure"
)

func main() {
	storyFile := flag.String("story_json", "./data/gopher.json", "the json story to load")
	flag.Parse()
	storyFD, err := os.Open(*storyFile)
	if err != nil {
		log.Fatalf("failed to load story: %v\n", err)
	}
	defer storyFD.Close()

	story, err := adventure.LoadStory(storyFD)
	if err != nil {
		log.Fatalf("failed to parse story: %v\n", err)
	}

	if err := http.ListenAndServe("localhost:8080", story); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
