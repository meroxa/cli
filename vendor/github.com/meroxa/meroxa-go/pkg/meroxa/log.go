package meroxa

import (
	"time"
)

type Logs struct {
	Data     []LogData `json:"data"`
	Metadata Metadata  `json:"metadata"`
}

type LogData struct {
	Log       string    `json:"log"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type Metadata struct {
	Limit  int       `json:"limit"`
	Query  string    `json:"query"`
	Source string    `json:"source"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
