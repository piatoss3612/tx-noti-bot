package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/tx-noti-bot/internal/app"
	"golang.org/x/exp/slog"
)

type bot struct {
	name    string
	session *discordgo.Session
}

func New(name, token string) (app.App, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &bot{name, session}, nil
}

func (b *bot) Setup() app.App {
	slog.Info(fmt.Sprintf("Start setup process of %s", b.name))

	b.session.AddHandler(ping)
	b.session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages

	slog.Info(fmt.Sprintf("%s is now ready to start!", b.name))

	return b
}

func (b *bot) Open() (<-chan bool, error) {
	slog.Info("Open connection to discord server...")

	err := b.session.Open()
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("%s is now running!", b.name))

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
	slog.Info(fmt.Sprintf("Close connection of %s...", b.name))

	_ = b.session.Close()

	slog.Info("Done closing process!")

	return nil
}
