package crypto_currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"net/http"
)

type providerDTO struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
}

type Provider struct {
	client *http.Client
	url    string
	key    string
}

func NewProvider(client *http.Client, url string, key string) *Provider {
	return &Provider{client: client, url: url, key: key}
}

func (p *Provider) GetData(ctx context.Context) ([]*CryptoCurrency, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-cg-pro-api-key", p.key)

	res, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Do request: %w", err)
	}

	defer res.Body.Close()

	dtos := []providerDTO{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Providing data error: %w", err)
	}

	err = json.Unmarshal(body, &dtos)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error: %w", err)
	}

	return mapDTOsToEntity(dtos), nil
}
func mapDTOsToEntity(dtos []providerDTO) []*CryptoCurrency {
	entities := make([]*CryptoCurrency, 0, len(dtos))

	for _, dto := range dtos {
		entities = append(entities, &CryptoCurrency{Title: dto.Symbol, Cost: dto.CurrentPrice, Inserted: time.Now()})
	}

	return entities
}
