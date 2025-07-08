package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

var globalLogger Logger

func SetLogger(l Logger) {
	globalLogger = l
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		latency := time.Since(start)
		globalLogger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
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
