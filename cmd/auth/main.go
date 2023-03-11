package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/piatoss3612/tx-noti-bot/internal/app/rest"
	hdr "github.com/piatoss3612/tx-noti-bot/internal/handler/auth"
	"github.com/piatoss3612/tx-noti-bot/internal/logger"
	"github.com/piatoss3612/tx-noti-bot/internal/repository/user/mongo"
	rc "github.com/piatoss3612/tx-noti-bot/internal/routes/auth"

	_ "github.com/joho/godotenv/autoload"
)

var (
	NAME = "authentication"
	PORT = "3000"
)

func main() {
	logger.SetStructuredLogger(NAME, os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo, err := mongo.New(ctx, NAME, os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = repo.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	hdr := hdr.New(repo)

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
