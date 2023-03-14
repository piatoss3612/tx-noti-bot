package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/piatoss3612/tx-notification/internal/app/discord"
	hdr "github.com/piatoss3612/tx-notification/internal/handler/discord"
	"github.com/piatoss3612/tx-notification/internal/logger"
	"golang.org/x/exp/slog"
)

var BOT_NAME = "Tx-Notification"

func main() {
	logger.SetStructuredLogger(BOT_NAME, os.Stdout)

	handler := hdr.New()

	slog.Info("Starting setup for running discord bot...")

	bot, err := discord.New(BOT_NAME, os.Getenv("DISCORD_TOKEN"), handler)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Open connection to discord server")

	shutdown, err := bot.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = bot.Close()
	}()

	slog.Info("Discord bot is now running!")

	<-shutdown

	slog.Info("Shutdown discord bot...")
}
