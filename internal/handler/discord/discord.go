package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/tx-notification/internal/handler"
)

var ErrTargetUnsupported = errors.New("target is unsupported")

type discordHandler struct{}

func New() handler.Handler {
	return &discordHandler{}
}

func (d *discordHandler) Inject(target any) error {
	session, ok := (target).(*discordgo.Session)
	if !ok {
		return ErrTargetUnsupported
	}

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages
	session.AddHandler(d.ping)

	return nil
}
