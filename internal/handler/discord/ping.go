package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (d *discordHandler) ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "pong")
		if err != nil {
			log.Println(err)
		}
	}
}
