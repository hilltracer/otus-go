package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

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

func (s *Storage) inRange(userID string, from, to time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []storage.Event
	for _, ev := range s.events {
		if ev.UserID != userID {
			continue
		}
		if ev.StartTime.Before(to) && !ev.StartTime.Before(from) {
			out = append(out, ev)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartTime.Before(out[j].StartTime) })
	return out
}

func (s *Storage) ListDay(_ context.Context, userID string, date time.Time) ([]storage.Event, error) {
	from := date.Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)
	return s.inRange(userID, from, to), nil
}

func (s *Storage) ListWeek(_ context.Context, userID string, weekStart time.Time) ([]storage.Event, error) {
	from := weekStart.Truncate(24 * time.Hour)
	to := from.AddDate(0, 0, 7)
	return s.inRange(userID, from, to), nil
}

func (s *Storage) ListMonth(_ context.Context, userID string, monthStart time.Time) ([]storage.Event, error) {
	from := time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	to := from.AddDate(0, 1, 0)
	return s.inRange(userID, from, to), nil
}
