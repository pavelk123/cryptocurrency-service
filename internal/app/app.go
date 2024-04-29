package app

import (
	"context"
	"fmt"
	"github.com/pavelk123/cryptocurrency-service/internal/delivery/rest"
	"github.com/pavelk123/cryptocurrency-service/internal/provider/geko"
	"github.com/pavelk123/cryptocurrency-service/internal/repository/pg"
	"github.com/pavelk123/cryptocurrency-service/internal/service"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pavelk123/cryptocurrency-service/config"

	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	httpServer *http.Server
	cfg        *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	return &App{
		cfg: cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := InitDbConn(a.cfg.DB)
	if err != nil {
		return fmt.Errorf("Faild to init db: %w", err)
	}

	repo := pg.NewRepository(db, a.cfg.DB)
	provider := geko.NewProvider(&http.Client{}, a.cfg.ProviderApiUrl, a.cfg.ProviderApiKey)
	service := service.NewService(a.cfg, logger, repo, provider)

	router := gin.Default()
	handler := rest.NewHandler(logger, service)

	group := router.Group("api/v1/rates")
	{
		group.GET("/", handler.GetAll)
		group.GET("/:title", handler.GetByTitle)
	}
	a.httpServer = &http.Server{
		Addr:         a.cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  time.Duration(a.cfg.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(a.cfg.WriteTimeoutInSeconds) * time.Second,
	}

	go func() {
		logger.Info("Server was started:" + a.cfg.ServerAddress)

		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %w", err)
		}
	}()

	service.RunBackgroundUpdate(ctx)

	<-ctx.Done()
	return a.httpServer.Shutdown(ctx)
}

func InitDbConn(cfgDb *config.DbConfig) (*sqlx.DB, error) {
	connString := "postgres://" + cfgDb.DatabaseUser + ":" + cfgDb.DatabasePassword + "@" +
		cfgDb.DatabaseHost + ":" + cfgDb.DatabasePort + "/" +
		cfgDb.DatabaseName + "?sslmode=disable"

	dbConn, err := sqlx.Connect("postgres", connString)

	if err != nil {
		return nil, fmt.Errorf("Failde to connect to db", err)
	}

	return dbConn, nil
}
