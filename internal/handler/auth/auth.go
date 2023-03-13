package auth

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/piatoss3612/tx-noti-bot/internal/handler"
	"github.com/piatoss3612/tx-noti-bot/internal/middlewares"
	"github.com/piatoss3612/tx-noti-bot/internal/repository/user"
)

var ErrTargetUnsupported = errors.New("target is unsupported")

type authHandler struct {
	repo user.UserRepository
}

func New(repo user.UserRepository) handler.Handler {
	return &authHandler{repo: repo}
}

func (a *authHandler) Inject(target any) error {
	mux, ok := target.(*chi.Mux)
	if !ok {
		return ErrTargetUnsupported
	}

	mux.Use(middlewares.Logger)
	mux.Use(middlewares.Recovery(true))
	mux.Use(middlewares.Ping("/auth/v1/ping"))

	mux.Route("/auth/v1", func(r chi.Router) {
		r.Use(middlewares.RestrictContentTypeToJSON(false))

		r.Route("/user", func(sr chi.Router) {
			sr.Post("/register", a.registerUser)
			sr.Post("/login", a.loginUser)
			sr.Post("/delete", a.deleteUser)
		})

		r.Route("/otp", func(sr chi.Router) {
			sr.Post("/generate", a.generateOTP)
			sr.Post("/verify", a.verifyOTP)
			sr.Post("/validate", a.validateOTP)
			sr.Post("/disable", a.disableOTP)
		})
	})

	return nil
}
