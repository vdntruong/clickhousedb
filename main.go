package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("Starting application...")
	MustLoadEnv()

	if err := Migrate(); err != nil {
		log.Fatal(fmt.Errorf("migration failed: %w", err))
	}
	log.Println("Migrations complete. Application can proceed.")

	chConn := MustLoadClickHouseConn()
	defer chConn.Close()

	samplePrv := NewSampleProvider(chConn)
	sampleSrv := NewSampleAdapter(samplePrv)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := sampleSrv.AddData(ctx); err != nil {
		log.Printf("failed to add sample data to database: %v", err)
	}
}
