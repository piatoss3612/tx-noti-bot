package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/piatoss3612/tx-noti-bot/internal/app/bot"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
)

func main() {
	logger.SetStructuredLogger("tx-noti-bot", os.Stdout)

	b, err := bot.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	shutdown, err := b.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = b.Close()
	}()

	<-shutdown
}
