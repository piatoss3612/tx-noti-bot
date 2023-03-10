package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/piatoss3612/tx-noti-bot/internal/handler"
	"github.com/piatoss3612/tx-noti-bot/internal/routes"
)

type authRouteController struct {
	hdr handler.Handler
}

func New(hdr handler.Handler) routes.RouteController {
	return &authRouteController{hdr: hdr}
}

func (a *authRouteController) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	a.hdr.Inject(mux)

	return mux
}
