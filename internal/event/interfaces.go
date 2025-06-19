package event

import "context"

type Repository interface {
	Create(ctx context.Context, event *Event) error
	Get(ctx context.Context, id string) (*Event, error)
}

type Service interface {
	Repository
}
