package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/piatoss3612/tx-notification/internal/app"
	"github.com/piatoss3612/tx-notification/internal/routes"
)

var (
	ErrInvalidPortNumber     = errors.New("invalid port number")
	ErrRouteControllerMissed = errors.New("route controller is missed; but is required")
)

type rest struct {
	name string
	port string
	rc   routes.RouteController
	srv  *http.Server
}

func New(name, port string, rc routes.RouteController) (app.App, error) {
	n, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Join(ErrInvalidPortNumber, err)
	}

	if n < 1 || n > 65535 {
		return nil, ErrInvalidPortNumber
	}

	if rc == nil {
		return nil, ErrRouteControllerMissed
	}

	return &rest{name: name, port: port, rc: rc}, nil
}

func (r *rest) Setup() app.App {
	r.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", r.port),
		Handler: r.rc.Routes(),
	}

	return r
}

func (r *rest) Open() (<-chan bool, error) {
	shutdown := make(chan bool)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		<-gracefulShutdown

		close(gracefulShutdown)
		close(shutdown)
	}()

	return shutdown, nil
}

func (r *rest) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
