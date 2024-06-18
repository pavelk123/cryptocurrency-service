package cryptocurr

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/pavelk123/cryptocurrency-service/config"
)

type repository interface {
	GetByTitle(ctx context.Context, title string) (*CryptoCurrency, error)
	Add(ctx context.Context, cc *CryptoCurrency) error
	List(ctx context.Context) ([]*CryptoCurrency, error)
	GetStats(ctx context.Context, model *CryptoCurrency) (*Stats, error)
}

type provider interface {
	GetData(ctx context.Context) ([]*CryptoCurrency, error)
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

func (s *Service) GetAll(ctx context.Context) ([]DTO, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("failed to list entities from repo:" + err.Error())

		return nil, err
	}

	dtos := make([]DTO, 0, len(list))

	for _, model := range list {
		statsModel, err := s.repo.GetStats(ctx, model)
		if err != nil {
			s.logger.Error("failed to get stats for model: " + model.Title + " :" + err.Error())

			return nil, err
		}

		dtos = append(dtos, *NewDTO(model, statsModel))
	}

	return dtos, nil
}
func (s *Service) GetByTitle(ctx context.Context, title string) (*DTO, error) {
	model, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		s.logger.Error("failed to get model by title: " + title + " :" + err.Error())

		return nil, err
	}

	statsModel, err := s.repo.GetStats(ctx, model)
	if err != nil {
		s.logger.Error("failed to get Stats for model:" + model.Title + " :" + err.Error())

		return nil, err
	}

	dto := NewDTO(model, statsModel)

	return dto, nil
}

func (s *Service) updateData(ctx context.Context) error {
	data, err := s.provider.GetData(ctx)
	if err != nil {
		s.logger.Error("update data failed:", err)

		return err
	}

	for _, entity := range data {
		err = s.repo.Add(ctx, entity)
		if err != nil {
			s.logger.Error("problem with inserting data from provider: %w", err)

			return err
		}
	}

	return nil
}

func (s *Service) RunBackgroundUpdate(ctx context.Context) {
	go func() {
		frequency := time.Minute * time.Duration(s.cfg.UpdateTimeInMinutes)

		ticker := time.NewTicker(frequency)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.logger.Info("start updating data with frequency " + strconv.Itoa(s.cfg.UpdateTimeInMinutes) + " minutes")
				err := s.updateData(ctx)

				switch {
				case err != nil:
					s.logger.Error("error updating data:%w", err)
				case err == nil:
					s.logger.Info("data was updated successfully")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
