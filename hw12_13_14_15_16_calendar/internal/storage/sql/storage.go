package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage { return &Storage{db: db} }

func Connect(ctx context.Context, dsn string) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return &Storage{db: db}, nil
}

func (s *Storage) Close(_ context.Context) error { return s.db.Close() }

func (s *Storage) CreateEvent(ctx context.Context, e storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // safe if already committed

	// overlap check
	var exists bool
	end := e.StartTime.Add(e.Duration)
	queryOverlap := `SELECT true FROM events
        WHERE user_id=$1 AND start_time < $3 AND (start_time + duration) > $2 LIMIT 1`
	if err = tx.QueryRowContext(ctx, queryOverlap, e.UserID, e.StartTime, end).Scan(&exists); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if exists {
		return storage.ErrDateBusy
	}

	// insert
	insert := `INSERT INTO events
        (id, title, start_time, duration, description, user_id, notify_before)
        VALUES (:id, :title, :start_time, :duration, :description, :user_id, :notify_before)`
	if _, err = tx.NamedExecContext(ctx, insert, map[string]any{
		"id":            e.ID,
		"title":         e.Title,
		"start_time":    e.StartTime,
		"duration":      e.Duration,
		"description":   e.Description,
		"user_id":       e.UserID,
		"notify_before": e.NotifyBefore,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// exists?
	var origCount int
	if err = tx.GetContext(ctx, &origCount, `SELECT 1 FROM events WHERE id=$1`, e.ID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}

	// overlap check (exclude self)
	end := e.StartTime.Add(e.Duration)
	var exists bool
	queryOverlap := `SELECT true FROM events
        WHERE user_id=$1 AND id <> $4 AND start_time < $3 AND (start_time + duration) > $2 LIMIT 1`
	if err = tx.QueryRowContext(ctx, queryOverlap, e.UserID, e.StartTime, end, e.ID).Scan(&exists); err != nil &&
		!errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if exists {
		return storage.ErrDateBusy
	}

	// update
	upd := `UPDATE events
        SET title=:title, start_time=:start_time, duration=:duration,
			description=:description, user_id=:user_id, notify_before=:notify_before
        WHERE id=:id`
	if _, err = tx.NamedExecContext(ctx, upd, map[string]any{
		"id":            e.ID,
		"title":         e.Title,
		"start_time":    e.StartTime,
		"duration":      e.Duration,
		"description":   e.Description,
		"user_id":       e.UserID,
		"notify_before": e.NotifyBefore,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return storage.ErrNotFound
	}
	return nil
}
