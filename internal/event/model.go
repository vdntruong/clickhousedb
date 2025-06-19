package event

import (
	"time"
)

type Event struct {
	ID        string    `json:"id" ch:"id"`
	Operation string    `json:"operation" ch:"operation"`
	Detail    Detail    `json:"detail" ch:"detail"`   // JSON
	Options   []Option  `json:"options" ch:"options"` // Nested
	Records   []Record  `json:"records"`              // Array(JSON)
	Timestamp time.Time `json:"timestamp" ch:"timestamp"`
}

type Record struct {
	Name        string `json:"name"`
	IsDefault   bool   `json:"is_default"`
	Description string `json:"description"`
}

type Option struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

type Detail struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
