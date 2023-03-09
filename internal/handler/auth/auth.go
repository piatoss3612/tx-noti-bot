package auth

import (
	"net/http"

	"github.com/piatoss3612/tx-noti-bot/internal/handler"
)

type authHandler struct {
	// TODO: appropriate fields
}

func New() (handler.Handler, error) {
	return &authHandler{}, nil
}

func (a *authHandler) Routes() http.Handler {
	return nil
}

func (a *authHandler) Cleanup() error {
	return nil
}
