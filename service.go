package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type SampleRepository interface {
	BatchInsert(ctx context.Context, dataCh <-chan SampleData, batchSizeThreshold uint32) error
}

// SampleAdapter implements the business logic layer for sample data operations.
// It coordinates database operations and implements higher-level workflows.
type SampleAdapter struct {
	repo SampleRepository
}

// NewSampleAdapter creates a service adapter that implements business logic for sample data.
// It wraps a SampleProvider to provide higher-level operations.
// Parameters:
//   - provider: Data repository for sample operations
//
// Returns:
//   - *SampleAdapter: Service adapter for sample data operations
func NewSampleAdapter(repo SampleRepository) *SampleAdapter {
	return &SampleAdapter{repo: repo}
}

func (s *SampleAdapter) AddData(ctx context.Context) error {
	var dataCh = make(chan SampleData)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err = s.repo.BatchInsert(ctx, dataCh, 1_000)
		return err
	})

	eg.Go(func() error {
		defer close(dataCh)
		for i := 0; i < 10_000; i++ {
			data := SampleData{
				ID:        uint64(i),
				Name:      fmt.Sprintf("Event Name %d", i%100),
				Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
				Value:     float64(i) * 1.1,
				Tags:      []string{fmt.Sprintf("tag%d", i%5), fmt.Sprintf("region%d", i%3)},
			}
			select {
			case dataCh <- data:
			case <-ctx.Done():
				return nil
			}
		}
		return nil
	})

	return eg.Wait()
}
