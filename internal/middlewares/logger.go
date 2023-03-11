package middlewares

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type WrappedResponseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func WrapResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{w, http.StatusOK, 0, false}
}

func (w *WrappedResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := WrapResponseWriter(w)
		start := time.Now()

		next.ServeHTTP(wrapped, r)

		slog.Info(
			"request completed",
			"method", r.Method,
			"status", wrapped.status,
			"uri", r.RequestURI,
			"from", r.RemoteAddr,
			"duration", time.Since(start),
			"size", wrapped.size,
		)
	})
}
