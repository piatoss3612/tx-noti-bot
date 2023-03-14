package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/piatoss3612/tx-notification/internal/app/rest"
	hdr "github.com/piatoss3612/tx-notification/internal/handler/auth"
	"github.com/piatoss3612/tx-notification/internal/logger"
	"github.com/piatoss3612/tx-notification/internal/repository/user/mongo"
	rc "github.com/piatoss3612/tx-notification/internal/routes/auth"
	"golang.org/x/exp/slog"

	_ "github.com/joho/godotenv/autoload"
)

var (
	NAME = "authentication"
	PORT = "3000"
)

func main() {
	logger.SetStructuredLogger(NAME, os.Stdout)

	slog.Info("Starting setup for running application...")

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

	slog.Info("Successfully connected to MongoDB")

	hdr := hdr.New(repo)

	rc := rc.New(hdr)

	app, err := rest.New(NAME, PORT, rc)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info(fmt.Sprintf("Running application on port %s", PORT))

	shutdown, err := app.Setup().Open()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	slog.Info("Application is now running!")

	<-shutdown

	slog.Info("Shutdown application...")
}
