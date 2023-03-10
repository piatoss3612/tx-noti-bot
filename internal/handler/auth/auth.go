package auth

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/piatoss3612/tx-noti-bot/internal/handler"
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

	// TODO: add middlewares?

	mux.Route("/auth/v1", func(r chi.Router) {
		r.Route("/user", func(sr chi.Router) {
			sr.Post("/register", a.RegisterUser)
			sr.Post("/delete", a.DeleteUser)
		})

		r.Route("/otp", func(sr chi.Router) {
			sr.Post("/enable", a.EnableOTP)
			sr.Post("/disable", a.DisableOTP)
			sr.Post("/verify", a.VerifyOTP)
			sr.Post("/validate", a.ValidateOTP)
		})
	})

	return nil
}

func (a *authHandler) Cleanup() error {
	return nil
}
