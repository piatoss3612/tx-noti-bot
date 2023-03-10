package main

import (
	"log"
	"os"

	"github.com/piatoss3612/tx-noti-bot/internal/app/rest"
	hdr "github.com/piatoss3612/tx-noti-bot/internal/handler/auth"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
	rc "github.com/piatoss3612/tx-noti-bot/internal/routes/auth"
)

var (
	NAME = "Authentication"
	PORT = "3000"
)

func main() {
	logger.SetStructuredLogger(NAME, os.Stdout)

	hdr, err := hdr.New()
	if err != nil {
		log.Fatal(err)
	}

	rc := rc.New(hdr)

	app, err := rest.New(NAME, PORT, rc)
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
