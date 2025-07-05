package sqlstorage

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

func mustEvent(id string, start time.Time, dur time.Duration) storage.Event {
	return storage.Event{ID: id, UserID: "u1", StartTime: start, Duration: dur}
}

func newMock() (*Storage, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	return New(sqlx.NewDb(db, "sqlmock")), mock, func() { _ = db.Close() }
}

func TestCreateEvent_Overlap(t *testing.T) {
	s, mock, cleanup := newMock()
	defer cleanup()

	ctx := context.Background()
	ev := mustEvent("1", time.Now(), time.Hour)

	// expect overlap query returning row => ErrDateBusy
	overlapRe := regexp.QuoteMeta(
		`SELECT true FROM events WHERE user_id=$1 AND start_time < $3 AND (start_time + duration) > $2 LIMIT 1`)
	mock.ExpectBegin()
	mock.ExpectQuery(overlapRe).
		WithArgs(ev.UserID, ev.StartTime, ev.StartTime.Add(ev.Duration)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectRollback()

	if err := s.CreateEvent(ctx, ev); !errors.Is(err, storage.ErrDateBusy) {
		t.Fatalf("want ErrDateBusy, got %v", err)
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	s, mock, cleanup := newMock()
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM events WHERE id=$1`)).
		WithArgs("42").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.DeleteEvent(context.Background(), "42"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
