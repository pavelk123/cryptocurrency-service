package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/pavelk123/cryptocurrency-service/config"
	"github.com/pavelk123/cryptocurrency-service/internal/app"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(ctx)
	if err != nil {
		log.Printf("failed to parsing config: %v", err)

		return
	}

	db, err := app.InitDBConn(&cfg.DB)
	if err != nil {
		log.Printf("faild to init db: %v", err)

		return
	}

	app, err := app.NewApp(cfg, db)
	if err != nil {
		log.Printf("failed to init app: %v", err)

		return
	}

	if err := app.Run(ctx); err != nil {
		log.Printf("failed to run app: %v", err)

		return
	}
}
