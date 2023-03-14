package discord

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/tx-notification/internal/app"
	"github.com/piatoss3612/tx-notification/internal/handler"
)

type discord struct {
	name    string
	handler handler.Handler
	session *discordgo.Session
}

func New(name, token string, handler handler.Handler) (app.App, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &discord{name: name, handler: handler, session: session}, nil
}

func (d *discord) Setup() app.App {
	d.handler.Inject(d.session)

	return d
}

func (d *discord) Open() (<-chan bool, error) {
	err := d.session.Open()
	if err != nil {
		return nil, err
	}

	shutdown := make(chan bool)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-gracefulShutdown
		close(gracefulShutdown)
		close(shutdown)
	}()

	return shutdown, nil
}

func (d *discord) Close() error {
	return d.session.Close()
}
