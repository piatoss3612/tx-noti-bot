package middlewares

import (
	"net/http"
	"strings"
)

// what's difference between r.URL.Path and r.RequestURI?

func Ping(endpoint string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == http.MethodGet || r.Method == http.MethodHead) && strings.EqualFold(r.URL.Path, endpoint) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("pong"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
