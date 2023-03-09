package main

import (
	"log"
	"os"

	"github.com/piatoss3612/tx-noti-bot/internal/app/rest"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
)

var (
	NAME = "Authentication"
	PORT = "3000"
)

func main() {
	logger.SetStructuredLogger(NAME, os.Stdout)

	app, err := rest.New(NAME, PORT, nil)
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
