package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Best struct {
	Score int       `json:"score"`
	WPM   int       `json:"wpm"`
	Combo int       `json:"combo"`
	When  time.Time `json:"when"`
}

func bestPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "termrain", "best.json")
}

func loadBest() Best {
	var b Best
	data, err := os.ReadFile(bestPath())
	if err == nil {
		json.Unmarshal(data, &b)
	}
	return b
}

func saveBest(b Best) error {
	p := bestPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(b, "", "  ")
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}
