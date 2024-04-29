package geko

import (
	"encoding/json"
	"fmt"
	"github.com/pavelk123/cryptocurrency-service/internal/entity"
	"io"
	"time"

	"net/http"
)

type Provider struct {
	client *http.Client
	url    string
	key    string
}

func NewProvider(client *http.Client, url string, key string) *Provider {
	return &Provider{client: client, url: url, key: key}
}

func (p *Provider) GetData() ([]*entity.CryptoCurrency, error) {

	req, _ := http.NewRequest("GET", p.url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-cg-pro-api-key", p.key)

	res, _ := p.client.Do(req)

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
func mapDTOsToEntity(dtos []providerDTO) []*entity.CryptoCurrency {
	entities := make([]*entity.CryptoCurrency, 0, len(dtos))

	for _, dto := range dtos {
		entities = append(entities, &entity.CryptoCurrency{Title: dto.Symbol, Cost: dto.CurrentPrice, Inserted: time.Now()})
	}

	return entities
}
