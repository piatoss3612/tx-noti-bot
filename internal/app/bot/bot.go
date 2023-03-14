package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/tx-notification/internal/app"
	"github.com/piatoss3612/tx-notification/internal/handler"
)

type bot struct {
	name    string
	handler handler.Handler
	session *discordgo.Session
}

func New(name, token string, handler handler.Handler) (app.App, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &bot{name: name, handler: handler, session: session}, nil
}

func (b *bot) Setup() app.App {
	b.handler.Inject(b.session)

	return b
}

func (b *bot) Open() (<-chan bool, error) {
	err := b.session.Open()
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

func (b *bot) Close() error {
	return b.session.Close()
}
