package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pavelk123/cryptocurrency-service/config"
	"github.com/pavelk123/cryptocurrency-service/internal/app"
	"log"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(ctx)
	if err != nil {
		log.Fatalf("Failed to parsing config %v", err)
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to init app %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Failed to run app %v", err)
	}
}
