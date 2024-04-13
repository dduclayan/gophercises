package main

import (
	"encoding/json"
	"os"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func loadStory(path string) (Story, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var s Story
	d := json.NewDecoder(f)
	if err := d.Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}
