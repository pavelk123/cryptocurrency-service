package cryptocurr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pavelk123/cryptocurrency-service/config"
	"io"
	"net/http"
	"time"
)

type providerDTO struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
}

type Provider struct {
	client *http.Client
	cfg    *config.ProviderConfig
}

func NewProvider(client *http.Client, cfg *config.ProviderConfig) *Provider {
	return &Provider{client: client, cfg: cfg}
}

func (p *Provider) GetData(ctx context.Context) ([]*CryptoCurrency, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-cg-pro-api-key", p.cfg.Key)

	res, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer res.Body.Close()

	dtos := []providerDTO{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("providing data error: %w", err)
	}

	if err = json.Unmarshal(body, &dtos); err != nil {
		return nil, fmt.Errorf("json.Unmarshal error: %w", err)
	}

	return p.mapDTOsToEntity(dtos), nil
}
func (p *Provider) mapDTOsToEntity(dtos []providerDTO) []*CryptoCurrency {
	entities := make([]*CryptoCurrency, 0, len(dtos))

	for _, dto := range dtos {
		entities = append(entities, &CryptoCurrency{Title: dto.Symbol, Cost: dto.CurrentPrice, Inserted: time.Now()})
	}

	return entities
}
