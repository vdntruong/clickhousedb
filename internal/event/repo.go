package event

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
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
	var dbEvent DBEvent
	row := c.conn.QueryRow(ctx, selectEventByIDSQL, id)
	if err := row.ScanStruct(&dbEvent); err != nil {
		return nil, fmt.Errorf("could not scan DBEvent struct: %w", err)
	}

	event, err := dbEvent.ToEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to convert DBEvent to Event: %w", err)
	}

	return event, nil
}

func (c *ClickHouseEventProvider) Create(ctx context.Context, event *Event) error {
	dbEvent, err := FromEvent(event)
	if err != nil {
		return fmt.Errorf("failed to convert Event to DBEvent: %w", err)
	}

	batch, err := c.conn.PrepareBatch(ctx, insertEventSQL)
	if err != nil {
		return fmt.Errorf("could not prepare batch: %w", err)
	}

	if err := batch.AppendStruct(dbEvent); err != nil {
		return fmt.Errorf("could not append DBEvent struct to batch: %w", err)
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
