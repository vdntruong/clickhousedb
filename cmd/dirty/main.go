package main

import (
	"clickhousedb/db"
	"clickhousedb/infra/env"
	"clickhousedb/internal/event"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("Starting application...")
	env.MustLoadEnv()
	env.DBMigrationPath = "/Volumes/external/Dev/go/clickhousedb/db/migrations"

	if err := db.Migrate(); err != nil {
		log.Fatal(fmt.Errorf("migration failed: %w", err))
	}
	log.Println("Migrations complete. Application can proceed.")

	chConn := db.MustLoadClickHouseConn()
	defer chConn.Close()

	chEventProvider := event.NewClickHouseEventProvider(chConn)
	//chEventSrv := event.NewAdapter(chEventProvider)

	newEvent := event.Event{
		ID:        "123",
		Operation: "example_operation",
		Detail: event.Detail{
			Source: "source-1",
			Target: "target-1",
		},
		Options: []event.Option{
			{
				Name:   "option-1",
				Type:   "type-1",
				Action: "action-1",
			},
			{
				Name:   "option-2",
				Type:   "type-2",
				Action: "action-2",
			},
		},
		Records: []event.Record{
			{
				Name:        "record-1",
				IsDefault:   true,
				Description: "description-1",
			},
			{
				Name:        "record-2",
				IsDefault:   false,
				Description: "description-2",
			},
		},
		Timestamp: time.Now(),
	}

	if err := chEventProvider.Create(context.Background(), &newEvent); err != nil {
		log.Fatal(fmt.Errorf("failed to create event: %w", err))
	}

	e, err := chEventProvider.Get(context.Background(), newEvent.ID)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get event: %w", err))
	}
	fmt.Println(e)

}
