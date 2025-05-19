package main

import "time"

// SampleData represents a record of data to be stored in ClickHouse.
// This structure matches the schema of the corresponding ClickHouse table.
type SampleData struct {
	ID        uint64
	Name      string
	Timestamp time.Time
	Value     float64
	Tags      []string
}
