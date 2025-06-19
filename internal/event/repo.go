package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "embed"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/chcol"
)

type ClickHouseEventProvider struct {
	conn clickhouse.Conn
}

func NewClickHouseEventProvider(conn clickhouse.Conn) *ClickHouseEventProvider {
	return &ClickHouseEventProvider{
		conn: conn,
	}
}

func (c *ClickHouseEventProvider) Get(ctx context.Context, id string) (*Event, error) {
	//var ev Event
	//row := c.conn.QueryRow(ctx, selectEventByIDSQL, id)
	//if err := row.ScanStruct(&ev); err != nil {
	//	return nil, fmt.Errorf("could not scan event struct: %w", err)
	//}
	//return &ev, nil

	// Define a temporary struct for scanning.
	// We now expect Records to be []*chcol.JSON directly from ScanStruct.
	var scannedEvent struct {
		ID        string        `ch:"id"`
		Operation string        `ch:"operation"`
		Detail    Detail        `ch:"detail"`
		Options   []Option      `ch:"options"`
		Records   []*chcol.JSON `ch:"records"` // Changed to []*chcol.JSON
		Timestamp time.Time     `ch:"timestamp"`
	}

	row := c.conn.QueryRow(ctx, selectEventByIDSQL, id)
	// Use ScanStruct on the temporary struct.
	if err := row.ScanStruct(&scannedEvent); err != nil {
		return nil, fmt.Errorf("could not scan event struct: %w", err)
	}

	// Now, iterate through the []*chcol.JSON and unmarshal each into your Record struct
	var records []Record
	for _, chJSONRecord := range scannedEvent.Records {
		if chJSONRecord == nil {
			// Handle cases where a JSON object in the array might be NULL/empty
			continue
		}
		var record Record
		// chcol.JSON has a MarshalJSON method which provides the raw JSON bytes.
		jsonBytes, err := chJSONRecord.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chcol.JSON to bytes: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal record JSON bytes: %w", err)
		}
		records = append(records, record)
	}

	// Construct the final Event object with the properly unmarshaled Records
	ev := &Event{
		ID:        scannedEvent.ID,
		Operation: scannedEvent.Operation,
		Detail:    scannedEvent.Detail,
		Options:   scannedEvent.Options,
		Records:   records, // Assign the unmarshaled records
		Timestamp: scannedEvent.Timestamp,
	}

	return ev, nil
}

func (c *ClickHouseEventProvider) Create(ctx context.Context, event *Event) error {
	batch, err := c.conn.PrepareBatch(ctx, insertEventSQL)
	if err != nil {
		return fmt.Errorf("could not prepare batch: %w", err)
	}

	//if err := batch.AppendStruct(event); err != nil {
	//	return fmt.Errorf("could not append struct to batch: %w", err)
	//}

	// Manually prepare slices for Nested type fields
	var optionsNames []string
	var optionsTypes []string
	var optionsActions []string

	for _, opt := range event.Options {
		optionsNames = append(optionsNames, opt.Name)
		optionsTypes = append(optionsTypes, opt.Type)
		optionsActions = append(optionsActions, opt.Action)
	}

	// Use batch.Append to explicitly pass all column values
	err = batch.Append(
		event.ID,
		event.Operation,
		event.Detail,
		optionsNames,
		optionsTypes,
		optionsActions,
		event.Records,
		event.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("could not append struct to batch: %w", err)
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("could not send batch: %w", err)
	}

	return nil

}

var (
	//go:embed queries/events_insert.sql
	insertEventSQL string

	//go:embed queries/events_select_by_id.sql
	selectEventByIDSQL string
)
