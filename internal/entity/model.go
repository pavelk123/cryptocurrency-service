package entity

import "time"

// CryptoCurrency is db entity
type CryptoCurrency struct {
	Title    string    `db:"title"`
	Cost     float64   `db:"cost"`
	Inserted time.Time ` db:"inserted"`
}

type Stats struct {
	MaxCostPerDay         float64 `db:"max_cost_per_day"`
	MinCostPerDay         float64 ` db:"min_cost_per_day"`
	ChangePerHourPercents float64 ` db:"percent_change_per_hour"`
}

// DTO is Data Transfer object
type DTO struct {
	Title                 string    `json:"title"`
	Cost                  float64   `json:"cost" `
	LastUpdate            time.Time `json:"last_update"`
	MaxCostPerDay         float64   `json:"max_cost_per_day" `
	MinCostPerDay         float64   `json:"min_cost_per_day" `
	ChangePerHourPercents float64   `json:"change_per_hour_percents" `
}

func NewDTO(model *CryptoCurrency, statsModel *Stats) (dto *DTO) {
	return &DTO{
		Title:                 model.Title,
		Cost:                  model.Cost,
		LastUpdate:            model.Inserted,
		MinCostPerDay:         statsModel.MinCostPerDay,
		MaxCostPerDay:         statsModel.MaxCostPerDay,
		ChangePerHourPercents: statsModel.ChangePerHourPercents,
	}
}
