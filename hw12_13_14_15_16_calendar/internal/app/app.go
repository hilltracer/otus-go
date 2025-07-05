package app

import (
	"context"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger Logger
	store  storage.Repository
}

type Logger interface {
	Info(string)
	Error(string)
}

func New(logger Logger, storage storage.Repository) *App {
	return &App{logger: logger, store: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	e := storage.Event{ID: id, Title: title}
	return a.store.CreateEvent(ctx, e)
}
