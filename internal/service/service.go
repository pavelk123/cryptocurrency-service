package service

import (
	"context"
	"fmt"
	"github.com/pavelk123/cryptocurrency-service/config"
	"github.com/pavelk123/cryptocurrency-service/internal/entity"
	"log/slog"
	"strconv"
	"time"
)

type repository interface {
	GetByTitle(ctx context.Context, title string) (*entity.CryptoCurrency, error)
	Add(ctx context.Context, cc *entity.CryptoCurrency) error
	List(ctx context.Context) ([]*entity.CryptoCurrency, error)
	GetStats(ctx context.Context, model *entity.CryptoCurrency) (*entity.Stats, error)
}

type provider interface {
	GetData() ([]*entity.CryptoCurrency, error)
}

type Service struct {
	cfg      *config.Config
	logger   *slog.Logger
	repo     repository
	provider provider
}

func NewService(cfg *config.Config, logger *slog.Logger, repo repository, provider provider) *Service {
	return &Service{
		cfg:      cfg,
		logger:   logger,
		repo:     repo,
		provider: provider,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*entity.DTO, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*entity.DTO, 0, len(list))

	for _, model := range list {
		statsModel, err := s.repo.GetStats(ctx, model)
		if err != nil {
			s.logger.Error("Failed to get stats for models:" + err.Error())
			return nil, fmt.Errorf("Error get stats: %w", err)
		}
		dtos = append(dtos, entity.NewDTO(model, statsModel))
	}

	return dtos, nil
}
func (s *Service) GetByTitle(ctx context.Context, title string) (*entity.DTO, error) {
	model, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		s.logger.Error(err.Error())
	}

	statsModel, err := s.repo.GetStats(ctx, model)
	if err != nil {
		s.logger.Error(err.Error())
	}

	dto := entity.NewDTO(model, statsModel)
	return dto, nil
}

func (s *Service) updateData(ctx context.Context) error {
	data, err := s.provider.GetData()
	if err != nil {
		s.logger.Error("Update data failed:", err)
		return fmt.Errorf("Update data failed: %w", err)
	}

	for _, entity := range data {
		err = s.repo.Add(ctx, entity)
		if err != nil {
			s.logger.Error("Problem with inserting data from provider: %w", err)
			return err
		}
	}

	return nil
}

func (s *Service) RunBackgroundUpdate(ctx context.Context) {
	go func() {
		duration := time.Minute * time.Duration(s.cfg.UpdateTimeInMinutes)

		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.logger.Info("Start updating data with durartion " + strconv.Itoa(s.cfg.UpdateTimeInMinutes) + " minutes")
				err := s.updateData(ctx)
				if err != nil {
					s.logger.Error("Error updating data:%w", err)
				} else {
					s.logger.Info("Data was updated successfully")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
