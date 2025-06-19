package event

import "context"

type Adapter struct {
	repository Repository
}

func NewAdapter(repository Repository) *Adapter {
	return &Adapter{
		repository: repository,
	}
}

func (a *Adapter) Create(ctx context.Context, event *Event) error {
	return a.repository.Create(ctx, event)
}

func (a *Adapter) Get(ctx context.Context, id string) (*Event, error) {
	return a.repository.Get(ctx, id)
}
