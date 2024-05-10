package cryptocurr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) GetByTitle(ctx context.Context, title string) (*CryptoCurrency, error) {
	var model CryptoCurrency

	err := repo.db.QueryRowxContext(
		ctx,
		`
WITH
    recent_currencies AS (
		SELECT title, max(inserted) AS inserted
      	FROM cc.crypto_currency_journal
      	WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
      	GROUP BY title
    )

SELECT recent_currencies.title, recent_currencies.inserted,journal.cost
FROM recent_currencies
LEFT JOIN cc.crypto_currency_journal journal
    ON journal.inserted = recent_currencies.inserted
    AND journal.title=recent_currencies.title
WHERE ccj.title=$1 
`,
		title,
	).StructScan(&model)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotFound
	case err != nil:
		return nil, fmt.Errorf("exec: %w", err)
	}

	return &model, nil
}

func (repo *Repository) List(ctx context.Context) ([]*CryptoCurrency, error) {
	var models []*CryptoCurrency

	rows, err := repo.db.QueryxContext(
		ctx,
		`
WITH
    recent_currencies AS (
		SELECT title, max(inserted) AS inserted
      	FROM cc.crypto_currency_journal
      	WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
      	GROUP BY title
    )

SELECT recent_currencies.title, recent_currencies.inserted,journal.cost
FROM recent_currencies
LEFT JOIN cc.crypto_currency_journal journal
    ON journal.inserted = recent_currencies.inserted
    AND journal.title=recent_currencies.title
`,
	)
	if err != nil {
		return nil, fmt.Errorf("query:%w", err)
	}

	for rows.Next() {
		var model CryptoCurrency

		if err := rows.StructScan(&model); err != nil {
			return nil, fmt.Errorf("scan:%w", err)
		}

		models = append(models, &model)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close:%w", err)
	}

	return models, nil
}

func (repo *Repository) Add(ctx context.Context, model *CryptoCurrency) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO cc.crypto_currency_journal (title,cost,inserted)VALUES($1,$2,$3)",
		model.Title,
		model.Cost,
		model.Inserted,
	)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (repo *Repository) GetStats(ctx context.Context, model *CryptoCurrency) (*Stats, error) {
	var statsModel Stats

	err := repo.db.QueryRowxContext(
		ctx,
		`
WITH
    daily_results AS (
        SELECT title, MAX(cost) AS max_cost_per_day, MIN(cost) AS min_cost_per_day
        FROM cc.crypto_currency_journal
        WHERE inserted::date = CURRENT_DATE
        GROUP BY title
    ),

    percent_change_results AS (
        SELECT title, ((MAX(cost) - MIN(cost)) / MIN(cost)) * 100 AS percent_change_per_hour
        FROM cc.crypto_currency_journal
        WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
        GROUP BY title
    ),
            
    list_titles AS (
    	SELECT DISTINCT title 
    	FROM cc.crypto_currency_journal
    )

SELECT daily_results.max_cost_per_day,
       daily_results.min_cost_per_day,
       percent_change_results.percent_change_per_hour
FROM list_titles
JOIN daily_results ON daily_results.title = list_titles.title
JOIN percent_change_results ON percent_change_results.title = list_titles.title
WHERE list_titles.title = $1;
	`,
		model.Title,
	).StructScan(&statsModel)

	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	return &statsModel, nil
}
