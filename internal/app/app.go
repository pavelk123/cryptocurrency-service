package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/pavelk123/cryptocurrency-service/config"
	"github.com/pavelk123/cryptocurrency-service/internal/cryptocurr"
)

type App struct {
	httpServer *http.Server
	cfg        *config.Config
	db         *sqlx.DB
}

func NewApp(cfg *config.Config, db *sqlx.DB) (*App, error) {
	return &App{
		cfg: cfg,
		db:  db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := cryptocurr.NewRepository(a.db)
	provider := cryptocurr.NewProvider(http.DefaultClient, a.cfg.ProviderAPIURL, a.cfg.ProviderAPIKey)
	service := cryptocurr.NewService(a.cfg, logger, repo, provider)

	router := gin.Default()
	handler := cryptocurr.NewHandler(logger, service)

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
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	service.RunBackgroundUpdate(ctx)

	<-ctx.Done()

	return a.httpServer.Shutdown(ctx)
}

func InitDBConn(cfgDB *config.DBConfig) (*sqlx.DB, error) {
	connString := "postgres://" + cfgDB.DatabaseUser + ":" + cfgDB.DatabasePassword + "@" +
		cfgDB.DatabaseHost + ":" + cfgDB.DatabasePort + "/" +
		cfgDB.DatabaseName + "?sslmode=disable"

	dbConn, err := sqlx.Connect("postgres", connString)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: ", err)
	}

	return dbConn, nil
}
