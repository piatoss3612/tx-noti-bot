package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	s, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Open(); err != nil {
		log.Fatal(err)
	}

	_, err = s.ChannelMessageSend(os.Getenv("TEST_CHANNEL_ID"), "ping")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Close(); err != nil {
		log.Fatal(err)
	}
}
