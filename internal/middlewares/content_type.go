package middlewares

import (
	"net/http"
	"strings"
)

func RestrictContentTypeToJSON(allowZeroContentLength bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				if allowZeroContentLength {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
				return
			}

			if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
