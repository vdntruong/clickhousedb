package event

import (
	"time"
)

type Event struct {
	ID        string    `json:"id" ch:"id"`
	Operation string    `json:"operation" ch:"operation"`
	Detail    Detail    `json:"detail" ch:"detail"`   // JSON
	Options   []Option  `json:"options" ch:"options"` // Nested
	Records   []Record  `json:"records" ch:"records"` // Array(JSON)
	Timestamp time.Time `json:"timestamp" ch:"timestamp"`
}

type Record struct {
	Name        string `json:"name" ch:"name"`
	IsDefault   bool   `json:"is_default" ch:"is_default"`
	Description string `json:"description" ch:"description"`
}

type Option struct {
	Name   string `json:"name" ch:"name"`
	Type   string `json:"type" ch:"type"`
	Action string `json:"action" ch:"action"`
}

type Detail struct {
	Source string `json:"source" ch:"source"`
	Target string `json:"target" ch:"target"`
}
