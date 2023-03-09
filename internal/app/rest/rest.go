package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piatoss3612/tx-noti-bot/internal/app"
)

type rest struct {
	// TODO: appropriate fields
	addr string
	srv  *http.Server
}

func New(addr string) (app.App, error) {
	// TODO: address validation
	// TODO: add handler
	return &rest{addr: addr}, nil
}

func (a *rest) Setup() app.App {
	a.srv = &http.Server{
		Addr:    a.addr,
		Handler: nil,
	}
	return a
}

func (a *rest) Open() (<-chan bool, error) {
	shutdown := make(chan bool)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		_ = a.srv.ListenAndServe()

		<-gracefulShutdown

		close(gracefulShutdown)
		close(shutdown)
	}()

	return shutdown, nil
}

func (a *rest) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
