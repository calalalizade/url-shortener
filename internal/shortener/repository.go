package shortener

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(url, code string) (Url, error) {
	u := Url{
		Original: url,
		Code:     code,
	}

	err := r.db.QueryRow(
		"INSERT INTO urls (original, code) VALUES ($1, $2) RETURNING created_at, expiration_date",
		url, code,
	).Scan(&u.CreatedAt, &u.ExpirationDate)

	if err != nil {
		return Url{}, err
	}

	return u, nil
}

func (r *Repository) GetByCode(code string) (Url, error) {
	var u Url

	err := r.db.QueryRow(
		`SELECT id, original, code, created_at, expiration_date, click_count 
		 FROM urls WHERE code = $1 AND expiration_date > NOW()`,
		code,
	).Scan(&u.ID, &u.Original, &u.Code, &u.CreatedAt, &u.ExpirationDate, &u.ClickCount)

	return u, err
}

func (r *Repository) GetByOriginalUrl(originalUrl string) (Url, error) {
	var u Url
	err := r.db.QueryRow(
		`SELECT id, code, original, click_count, expiration_date, created_at 
         FROM urls 
         WHERE original = $1 AND expiration_date > NOW()
         ORDER BY created_at DESC
         LIMIT 1`,
		originalUrl,
	).Scan(&u.ID, &u.Code, &u.Original, &u.ClickCount, &u.ExpirationDate, &u.CreatedAt)

	return u, err
}

func (r *Repository) IncrementClickCount(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE urls SET click_count = click_count + 1 WHERE code = $1",
		code,
	)
	return err
}

func (r *Repository) GetStats(code string) (Url, error) {
	return r.GetByCode(code)
}
