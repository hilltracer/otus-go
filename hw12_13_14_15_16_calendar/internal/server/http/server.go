package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger Logger
	app    Application
	srv    *http.Server
}

type Logger interface {
	Info(string)
	Error(string)
}

type Application interface {
	// CreateEvent(ctx context.Context, id, title string) error
	CreateFullEvent(ctx context.Context, e storage.Event) error
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, id string) error

	ListDay(ctx context.Context, userID string, date time.Time) ([]storage.Event, error)
	ListWeek(ctx context.Context, userID string, weekStart time.Time) ([]storage.Event, error)
	ListMonth(ctx context.Context, userID string, monthStart time.Time) ([]storage.Event, error)
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		latency := time.Since(start)
		s.logger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			rw.status,
			latency.Milliseconds(),
			r.UserAgent(),
		))
	})
}

func NewServer(logger Logger, app Application, addr string) *Server {
	mux := http.NewServeMux()
	s := &Server{logger: logger, app: app}

	// hello world endpoint
	mux.Handle("/", s.loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})))

	mux.Handle("/events", s.loggingMiddleware(http.HandlerFunc(s.handleCreate)))          // POST
	mux.Handle("/events/", s.loggingMiddleware(http.HandlerFunc(s.handleUpdateDelete)))   // PUT / DELETE
	mux.Handle("/events/day", s.loggingMiddleware(http.HandlerFunc(s.handleListDay)))     // GET
	mux.Handle("/events/week", s.loggingMiddleware(http.HandlerFunc(s.handleListWeek)))   // GET
	mux.Handle("/events/month", s.loggingMiddleware(http.HandlerFunc(s.handleListMonth))) // GET

	s.srv = &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http listen: " + err.Error())
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// ---------- handlers ----------

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req createOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	e := storage.Event{
		ID: req.ID, Title: req.Title, StartTime: req.StartTime, Duration: req.Duration,
		Description: req.Description, UserID: req.UserID, NotifyBefore: req.NotifyBefore,
	}
	if err := s.app.CreateFullEvent(r.Context(), e); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleUpdateDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/events/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodDelete:
		if err := s.app.DeleteEvent(r.Context(), id); err != nil {
			s.writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut:
		var req createOrUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		req.ID = id
		e := storage.Event{
			ID: id, Title: req.Title, StartTime: req.StartTime, Duration: req.Duration,
			Description: req.Description, UserID: req.UserID, NotifyBefore: req.NotifyBefore,
		}
		if err := s.app.UpdateEvent(r.Context(), e); err != nil {
			s.writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleListDay(w http.ResponseWriter, r *http.Request) {
	s.handleListGeneric(w, r, s.app.ListDay, "date")
}

func (s *Server) handleListWeek(w http.ResponseWriter, r *http.Request) {
	s.handleListGeneric(w, r, s.app.ListWeek, "start")
}

func (s *Server) handleListMonth(w http.ResponseWriter, r *http.Request) {
	s.handleListGeneric(w, r, s.app.ListMonth, "start")
}

// helpers --------------------------------------------------------------

func (s *Server) handleListGeneric(
	w http.ResponseWriter,
	r *http.Request,
	fn func(context.Context, string, time.Time) ([]storage.Event, error),
	param string,
) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.URL.Query().Get("userId")
	raw := r.URL.Query().Get(param)
	if userID == "" || raw == "" {
		http.Error(w, "missing query params", http.StatusBadRequest)
		return
	}
	tm, err := time.Parse("2006-01-02", raw)
	if err != nil {
		http.Error(w, "bad date", http.StatusBadRequest)
		return
	}
	evs, err := fn(r.Context(), userID, tm)
	if err != nil {
		s.writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listResponse{Events: evs})
}

func (s *Server) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrDateBusy):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, storage.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
