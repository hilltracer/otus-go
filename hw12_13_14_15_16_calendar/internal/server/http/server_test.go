package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestEndpoints(t *testing.T) {
	st := memorystorage.New()
	ap := app.New(logger.New("error"), st)
	srv := NewServer(logger.New("error"), ap, "")
	ts := httptest.NewServer(srv.srv.Handler)
	defer ts.Close()

	base := time.Date(2025, 7, 3, 12, 0, 0, 0, time.UTC)

	// --- create ---
	body, _ := json.Marshal(map[string]any{
		"id": "e1", "title": "demo", "startTime": base,
		"duration": int64(time.Hour), "userId": "u1",
	})
	//nolint:noctx
	resp, _ := http.Post(ts.URL+"/events", "application/json", bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: want 201, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// --- list day ---
	dayURL := ts.URL + "/events/day?date=2025-07-03&userId=u1"
	//nolint:noctx
	resp, _ = http.Get(dayURL)
	var lr listResponse
	_ = json.NewDecoder(resp.Body).Decode(&lr)
	if len(lr.Events) != 1 {
		t.Fatalf("list day: want 1, got %d", len(lr.Events))
	}
	defer resp.Body.Close()
}
