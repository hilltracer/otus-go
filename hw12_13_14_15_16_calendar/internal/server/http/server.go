package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
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
	CreateEvent(ctx context.Context, id, title string) error
}

func NewServer(logger Logger, app Application, addr string) *Server {
	SetLogger(logger)

	mux := http.NewServeMux()
	s := &Server{logger: logger, app: app}

	// hello world endpoint
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})))

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
