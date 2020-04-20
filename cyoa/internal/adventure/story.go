package adventure

import (
	"encoding/json"
	"io"
)

type Story map[string]Storyblock

type Storyblock struct {
	Title   string              `json:"title"`
	Story   []string            `json:"story"`
	Options []StoryBlockOptions `json:"options"`
}

type StoryBlockOptions struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func LoadStory(rd io.Reader) (*Story, error) {
	var st Story
	if err := json.NewDecoder(rd).Decode(&st); err != nil {
		return nil, err
	}
	return &st, nil
}
