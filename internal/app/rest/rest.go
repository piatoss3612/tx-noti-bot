package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/piatoss3612/tx-noti-bot/internal/app"
	"github.com/piatoss3612/tx-noti-bot/internal/handler"
	"golang.org/x/exp/slog"
)

var ErrInvalidPortNumber = errors.New("invalid port number")

type rest struct {
	// TODO: appropriate fields
	name string
	port string
	hdr  handler.Handler
	srv  *http.Server
}

func New(name, port string, handler handler.Handler) (app.App, error) {
	n, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Join(ErrInvalidPortNumber, err)
	}

	if n < 1 || n > 65535 {
		return nil, ErrInvalidPortNumber
	}

	return &rest{name: name, port: port, hdr: handler}, nil
}

func (r *rest) Setup() app.App {
	slog.Info(fmt.Sprintf("Start setup %s server", r.name))

	r.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", r.port),
		Handler: r.hdr.Routes(),
	}

	slog.Info(fmt.Sprintf("%s server is now available", r.name))

	return r
}

func (r *rest) Open() (<-chan bool, error) {
	slog.Info(fmt.Sprintf("Open %s server...", r.name))

	shutdown := make(chan bool)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		_ = r.srv.ListenAndServe()
	}()

	go func() {
		<-gracefulShutdown

		close(gracefulShutdown)
		close(shutdown)
	}()

	slog.Info(fmt.Sprintf("%s is now running!", r.name))

	return shutdown, nil
}

func (r *rest) Close() error {
	slog.Info(fmt.Sprintf("Close %s server...", r.name))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	slog.Info("Done closing process!")

	return nil
}
