package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

func Recovery(enableStackTrace bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					var traced []byte

					if enableStackTrace {
						traced = debug.Stack()
					}

					var err error
					var ok bool

					if err, ok = rcv.(error); !ok {
						err = fmt.Errorf("%v", rcv)
					}

					slog.Error(
						"panic while handling request",
						err,
						"stack", string(traced),
					)

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
