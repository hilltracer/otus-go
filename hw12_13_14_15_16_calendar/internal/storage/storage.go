package storage

import "context"

type Repository interface {
	CreateEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	DeleteEvent(ctx context.Context, id string) error
}
