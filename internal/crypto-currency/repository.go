package crypto_currency

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
SELECT ccj.title, ccj.inserted,inst.cost
FROM cc.crypto_currency_journal inst
RIGHT JOIN
     (SELECT title, max(inserted) AS inserted
      FROM cc.crypto_currency_journal
      WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
      GROUP BY title
     ) AS ccj
    ON inst.inserted = ccj.inserted
    AND inst.title=ccj.title
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
SELECT ccj.title, ccj.inserted,inst.cost
FROM cc.crypto_currency_journal inst
RIGHT JOIN
     (SELECT title, max(inserted) AS inserted
      FROM cc.crypto_currency_journal
      WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
      GROUP BY title
     ) AS ccj
    ON inst.inserted = ccj.inserted
    AND inst.title=ccj.title
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
SELECT daily_results.max_cost_per_day,
       daily_results.min_cost_per_day,
       percent_change_results.percent_change_per_hour
FROM (SELECT title FROM cc.crypto_currency_journal GROUP BY title) ccj
INNER JOIN
    (SELECT title, MAX(cost) AS max_cost_per_day, MIN(cost) AS min_cost_per_day
     FROM cc.crypto_currency_journal
     WHERE inserted::date = CURRENT_DATE
     GROUP BY title) AS daily_results ON daily_results.title = ccj.title
INNER JOIN
    (SELECT title, ((MAX(cost) - MIN(cost)) / MIN(cost)) * 100 AS percent_change_per_hour
     FROM cc.crypto_currency_journal
     WHERE inserted >= CURRENT_TIMESTAMP - INTERVAL '1 hour'
     GROUP BY title) AS percent_change_results ON percent_change_results.title = ccj.title
WHERE ccj.title =$1;
	`,
		model.Title,
	).StructScan(&statsModel)

	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	return &statsModel, nil
}
