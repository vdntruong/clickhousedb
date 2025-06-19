package event

import (
	"encoding/json"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/chcol"
)

type DBEvent struct {
	Event

	OptionsName   []string     `ch:"options.name"`   // Nested type (options.name)
	OptionsType   []string     `ch:"options.type"`   // Nested type (options.type)
	OptionsAction []string     `ch:"options.action"` // Nested type (options.action)
	Records       []chcol.JSON `ch:"records"`        // Array(JSON)
}

// ToEvent converts a DBEvent to an Event.
func (dbe *DBEvent) ToEvent() (*Event, error) {
	// Copy simple fields directly from the embedded Event
	e := dbe.Event

	//// Unmarshal Detail
	//if len(dbe.Detail) > 0 {
	//	if err := json.Unmarshal(dbe.Detail, &e.Detail); err != nil {
	//		return nil, fmt.Errorf("failed to unmarshal Detail from DBEvent: %w", err)
	//	}
	//}

	// Reconstruct Options from Nested fields
	if len(dbe.OptionsName) != len(dbe.OptionsType) || len(dbe.OptionsName) != len(dbe.OptionsAction) {
		return nil, fmt.Errorf("mismatch in lengths of nested option fields from DBEvent")
	}
	e.Options = make([]Option, len(dbe.OptionsName))
	for i := range dbe.OptionsName {
		e.Options[i] = Option{
			Name:   dbe.OptionsName[i],
			Type:   dbe.OptionsType[i],
			Action: dbe.OptionsAction[i],
		}
	}

	// Unmarshal Records from []*chcol.JSON
	e.Records = make([]Record, len(dbe.Records))
	for i, chJSON := range dbe.Records {
		jsonBytes, err := chJSON.MarshalJSON() // chcol.JSON has MarshalJSON
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chcol.JSON to bytes for record %d: %w", i, err)
		}
		if err := json.Unmarshal(jsonBytes, &e.Records[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Record from DBEvent for record %d: %w", i, err)
		}
	}

	return &e, nil
}

// FromEvent converts an Event to a DBEvent.
func FromEvent(e *Event) (*DBEvent, error) {
	dbe := DBEvent{
		Event: *e, // Copy simple fields from Event to the embedded Event
	}

	//// Marshal Detail
	//detailBytes, err := json.Marshal(e.Detail)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to marshal Detail for DBEvent: %w", err)
	//}
	//dbe.Detail = json.RawMessage(detailBytes)

	// Populate Nested fields for Options
	//dbe.OptionsName = make([]string, len(e.Options))
	//dbe.OptionsType = make([]string, len(e.Options))
	//dbe.OptionsAction = make([]string, len(e.Options))
	//for i, opt := range e.Options {
	//	dbe.OptionsName[i] = opt.Name
	//	dbe.OptionsType[i] = opt.Type
	//	dbe.OptionsAction[i] = opt.Action
	//}

	// Marshal Records to []*chcol.JSON
	dbe.Records = make([]chcol.JSON, len(e.Records))
	for i, rec := range e.Records {
		// 1. Marshal the Record struct to a map[string]any
		var recMap map[string]any
		recBytes, err := json.Marshal(rec)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Record to JSON for DBEvent: %w", err)
		}
		if err := json.Unmarshal(recBytes, &recMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Record JSON to map[string]any for DBEvent: %w", err)
		}

		// 2. Create a new chcol.JSON instance
		chJSON := chcol.NewJSON()

		// 3. Populate the chcol.JSON instance using its Scan method with the map
		if err := chJSON.Scan(recMap); err != nil {
			return nil, fmt.Errorf("failed to scan map into chcol.JSON for record %d: %w", i, err)
		}
		dbe.Records[i] = *chJSON
	}

	return &dbe, nil
}
