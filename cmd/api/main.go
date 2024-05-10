package main

import (
	"context"
	"log"
	"log/slog"
	"os"

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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.New(ctx)
	if err != nil {
		logger.Error("failed to parsing config:", err.Error())

		return
	}

	db, err := app.InitDBConn(&cfg.DB)
	if err != nil {
		logger.Error("faild to init db:", err.Error())

		return
	}

	app, err := app.NewApp(cfg, db, logger)
	if err != nil {
		logger.Error("failed to init app:", err.Error())

		return
	}

	if err := app.Run(ctx); err != nil {
		logger.Error("failed to run app:", err.Error())

		return
	}
}
