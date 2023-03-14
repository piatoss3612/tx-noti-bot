package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/piatoss3612/tx-notification/internal/handler"
)

var ErrTargetUnsupported = errors.New("target is unsupported")

type botHandler struct {
}

func New() handler.Handler {
	return &botHandler{}
}

func (b *botHandler) Inject(target any) error {
	session, ok := (target).(*discordgo.Session)
	if !ok {
		return ErrTargetUnsupported
	}

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages
	session.AddHandler(b.ping)

	return nil
}
