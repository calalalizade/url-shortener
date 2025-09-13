package shortener

import "database/sql"

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
