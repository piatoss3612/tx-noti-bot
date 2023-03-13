package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/piatoss3612/tx-noti-bot/internal/app/bot"
	hdr "github.com/piatoss3612/tx-noti-bot/internal/handler/bot"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
	"golang.org/x/exp/slog"
)

var BOT_NAME = "Tx-Noti-Bot"

func main() {
	logger.SetStructuredLogger(BOT_NAME, os.Stdout)

	handler := hdr.New()

	slog.Info("Starting setup for running discord bot...")

	app, err := bot.New(BOT_NAME, os.Getenv("DISCORD_TOKEN"), handler)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Open connection to discord server")

	shutdown, err := app.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = app.Close()
	}()

	slog.Info("Discord bot is now running!")

	<-shutdown

	slog.Info("Shutdown discord bot...")
}
