package event

import (
	"encoding/json"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/chcol"
)

type DBEvent struct {
	Event

	OptionsName   []string     `ch:"options.name"`   // Nested type (options.name)		for inserting
	OptionsType   []string     `ch:"options.type"`   // Nested type (options.type)		for inserting
	OptionsAction []string     `ch:"options.action"` // Nested type (options.action) 	for inserting
	Records       []chcol.JSON `ch:"records"`        // Array(JSON) 					for inserting and receiving
}

// ToEvent converts a DBEvent to an Event.
func (dbe *DBEvent) ToEvent() (*Event, error) {
	// Copy simple fields directly from the embedded Event
	e := dbe.Event

	// Unmarshal Records from []*chcol.JSON
	e.Records = make([]Record, len(dbe.Records))
	for i, chJSON := range dbe.Records {
		if err := CHJSONToStruct(&chJSON, &e.Records[i]); err != nil {
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

	// Populate Nested fields
	dbe.OptionsName = make([]string, len(e.Options))
	dbe.OptionsType = make([]string, len(e.Options))
	dbe.OptionsAction = make([]string, len(e.Options))
	for i, opt := range e.Options {
		dbe.OptionsName[i] = opt.Name
		dbe.OptionsType[i] = opt.Type
		dbe.OptionsAction[i] = opt.Action
	}

	// Marshal Records to []*chcol.JSON
	dbe.Records = make([]chcol.JSON, len(e.Records))
	for i, rec := range e.Records {
		chJSON, err := StructToCHJSON(rec)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Record to chcol.JSON for record %d: %w", i, err)
		}

		dbe.Records[i] = *chJSON
	}

	return &dbe, nil
}

func StructToCHJSON[T any](v T) (*chcol.JSON, error) {
	var vMap map[string]any
	vBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %T to JSON: %w", v, err)
	}
	if err := json.Unmarshal(vBytes, &vMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T JSON to map[string]any: %w", v, err)
	}

	chJSON := chcol.NewJSON()
	if err := chJSON.Scan(vMap); err != nil {
		return nil, err
	}
	return chJSON, nil
}

func CHJSONToStruct[T any](chJSON *chcol.JSON, v *T) error {
	vBytes, err := chJSON.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal chcol.JSON to JSON: %w", err)
	}
	if err := json.Unmarshal(vBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to %T: %w", v, err)
	}
	return nil
}
