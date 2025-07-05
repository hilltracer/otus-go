package memorystorage

import (
	"context"
	"sync"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) overlap(e storage.Event) bool {
	for _, ev := range s.events {
		if ev.UserID != e.UserID {
			continue
		}
		startA, endA := ev.StartTime, ev.StartTime.Add(ev.Duration)
		startB, endB := e.StartTime, e.StartTime.Add(e.Duration)
		if startB.Before(endA) && startA.Before(endB) {
			return true
		}
	}
	return false
}

func (s *Storage) CreateEvent(_ context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.overlap(e) {
		return storage.ErrDateBusy
	}
	s.events[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[e.ID]; !ok {
		return storage.ErrNotFound
	}
	if s.overlap(e) {
		return storage.ErrDateBusy
	}
	s.events[e.ID] = e
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}
	delete(s.events, id)
	return nil
}
