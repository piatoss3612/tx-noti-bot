package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/piatoss3612/tx-noti-bot/internal/app/bot"
	hdr "github.com/piatoss3612/tx-noti-bot/internal/handler/bot"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
)

var BOT_NAME = "Tx-Noti-Bot"

func main() {
	logger.SetStructuredLogger(BOT_NAME, os.Stdout)

	handler := hdr.New()

	app, err := bot.New(BOT_NAME, os.Getenv("DISCORD_TOKEN"), handler)
	if err != nil {
		log.Fatal(err)
	}

	shutdown, err := app.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = app.Close()
	}()

	<-shutdown
}
