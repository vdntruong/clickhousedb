package main

import (
	"context"
	_ "embed"
	"fmt"

	"clickhousedb/utils/concurrency"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

//go:embed queries/examples_insert.sql
var QueryInsertExample string

// SampleProvider implements database operations for sample data.
// It encapsulates ClickHouse-specific implementation details for data storage.
type SampleProvider struct {
	conn driver.Conn
}

// NewSampleProvider creates a new instance of SampleProvider with the given ClickHouse connection.
// The provider handles database operations for sample data.
// Parameters:
//   - conn: An established ClickHouse connection
//
// Returns:
//   - *SampleProvider: Data repository object for sample operations
func NewSampleProvider(conn driver.Conn) *SampleProvider {
	return &SampleProvider{conn: conn}
}

// BatchInsert adds multiple records to ClickHouse in efficient batches.
// It processes data from the input channel and sends in batches of the specified size.
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - dataCh: Channel of data records to insert
//   - batchSizeThreshold: Number of records per batch (defaults to 10,000 if zero)
//
// Returns:
//   - error: nil on success, or error details on failure
//
// Note: This method creates a new batch after each Send() operation to avoid "batch already sent" errors.
func (r *SampleProvider) BatchInsert(ctx context.Context, dataCh <-chan SampleData, batchSizeThreshold uint32) error {
	if batchSizeThreshold == 0 {
		batchSizeThreshold = 10_000 // ClickHouse recommends batches of at least 1,000 rows, ideally 10,000-100,000.
	}

	batch, err := r.conn.PrepareBatch(ctx, QueryInsertExample)
	if err != nil {
		return fmt.Errorf("could not prepare batch: %w", err)
	}

	var count uint32
	for data := range concurrency.OrDone(ctx, dataCh) {
		if err := batch.Append(
			data.ID,
			data.Name,
			data.Timestamp,
			data.Value,
			data.Tags,
		); err != nil {
			// Abort the batch if appending fails for some reason.
			// Depending on the error, you might want to retry or log and continue.
			if abortErr := batch.Abort(); abortErr != nil {
				return fmt.Errorf("failed to abort batch after append error: %w", abortErr)
			}
			return fmt.Errorf("could not append row to batch (ID %d): %w", data.ID, err)
		}

		count++

		if count%batchSizeThreshold == 0 {
			if err := batch.Send(); err != nil {
				return fmt.Errorf("could not send batch (ID %d): %w", data.ID, err)
			}

			// Create a new batch after sending
			batch, err = r.conn.PrepareBatch(ctx, QueryInsertExample)
			if err != nil {
				return fmt.Errorf("could not prepare new batch after sending: %w", err)
			}
		}
	}

	// Send any remaining data
	if batch.Rows() > 0 {
		if err := batch.Send(); err != nil {
			return fmt.Errorf("could not send final batch: %w", err)
		}
	}

	return nil
}
