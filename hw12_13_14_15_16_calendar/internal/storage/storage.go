package storage

import (
	"context"
	"time"
)

type Repository interface {
	CreateEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	DeleteEvent(ctx context.Context, id string) error

	ListDay(ctx context.Context, userID string, date time.Time) ([]Event, error)
	ListWeek(ctx context.Context, userID string, weekStart time.Time) ([]Event, error)
	ListMonth(ctx context.Context, userID string, monthStart time.Time) ([]Event, error)
}
