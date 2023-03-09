package main

import (
	"log"
	"os"

	"github.com/piatoss3612/tx-noti-bot/internal/app/rest"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
	"github.com/piatoss3612/tx-noti-bot/internal/routes/auth"
)

var (
	NAME = "Authentication"
	PORT = "3000"
)

func main() {
	logger.SetStructuredLogger(NAME, os.Stdout)

	hdr, err := auth.New()
	if err != nil {
		log.Fatal(err)
	}

	app, err := rest.New(NAME, PORT, hdr)
	if err != nil {
		log.Fatal(err)
	}

	shutdown, err := app.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Close()

	<-shutdown
}
