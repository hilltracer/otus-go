package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

//nolint:unparam
func mustEvent(id string, start time.Time, dur time.Duration) storage.Event {
	return storage.Event{ID: id, UserID: "u1", StartTime: start, Duration: dur}
}

func TestStorage(t *testing.T) {
	s := New()
	ctx := context.Background()

	start := time.Now()
	e1 := mustEvent("1", start, time.Hour)

	// create
	if err := s.CreateEvent(ctx, e1); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// overlap create
	eOverlap := mustEvent("2", start.Add(30*time.Minute), time.Hour)
	if err := s.CreateEvent(ctx, eOverlap); !errors.Is(err, storage.ErrDateBusy) {
		t.Fatalf("expected ErrDateBusy, got %v", err)
	}

	// update not found
	if err := s.UpdateEvent(ctx, mustEvent("42", start, time.Hour)); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// update overlap
	e3 := mustEvent("3", start.Add(2*time.Hour), time.Hour)
	_ = s.CreateEvent(ctx, e3)
	if err := s.UpdateEvent(ctx, storage.Event{
		ID: "3", UserID: "u1", StartTime: start.Add(30 * time.Minute),
		Duration: time.Hour,
	}); !errors.Is(err, storage.ErrDateBusy) {
		t.Fatalf("expected ErrDateBusy on update, got %v", err)
	}

	// delete ok
	if err := s.DeleteEvent(ctx, "1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if err := s.DeleteEvent(ctx, "3"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	// delete again -> not found
	if err := s.DeleteEvent(ctx, "1"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	base := time.Date(2025, 7, 1, 10, 0, 0, 0, time.UTC)
	add := func(id string, shiftDays int) {
		_ = s.CreateEvent(ctx, mustEvent(id, base.AddDate(0, 0, shiftDays), time.Hour))
	}
	add("d1", 0)  // 1 july
	add("d2", 2)  // 3 july
	add("d3", 10) // 11 july

	day, _ := s.ListDay(ctx, "u1", base)
	if len(day) != 1 {
		t.Fatalf("want 1 event on day, got %d", len(day))
	}
	week, _ := s.ListWeek(ctx, "u1", base)
	if len(week) != 2 {
		t.Fatalf("want 2 events in week, got %d", len(week))
	}
	month, _ := s.ListMonth(ctx, "u1", base)
	if len(month) != 3 {
		t.Fatalf("want 3 events in month, got %d", len(month))
	}
}

func TestStorage_ConcurrentSafety(_ *testing.T) {
	s := New()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("%d", i)
			ev := mustEvent(id, time.Now().Add(time.Duration(i)*time.Hour), time.Hour)
			_ = s.CreateEvent(ctx, ev)
		}(i)
	}
	wg.Wait()
	// no race conditions detected with -race flag
}
