package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/piatoss3612/tx-noti-bot/internal/routes"
)

type authRouteController struct {
	// TODO: appropriate fields
}

func New() (routes.RouteController, error) {
	return &authRouteController{}, nil
}

func (a *authRouteController) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Heartbeat("/api/auth/v1/ping"))

	mux.Route("/api/auth/v1", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		})
	})

	return mux
}
