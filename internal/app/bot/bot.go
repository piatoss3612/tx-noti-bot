package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	name    string
	session *discordgo.Session
}

func New(name, token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{name, session}, nil
}

func (b *Bot) Setup() *Bot {
	b.session.AddHandler(ping)
	b.session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages
	return b
}

func (b *Bot) Open() (<-chan bool, error) {
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

func (b *Bot) Close() error {
	return b.session.Close()
}
