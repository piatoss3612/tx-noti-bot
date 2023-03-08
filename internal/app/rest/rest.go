package rest

import (
	"github.com/piatoss3612/tx-noti-bot/internal/app"
)

type rest struct {
	// TODO: appropriate fields
}

func NewRest() app.App {
	return &rest{}
}

func (a *rest) Setup() app.App {
	return a
}

func (a *rest) Open() (<-chan bool, error) {
	return nil, nil
}

func (a *rest) Close() error {
	return nil
}
