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

func (r *Repository) GetByCode(code string) (Url, error) {
	var u Url

	err := r.db.QueryRow(
		`SELECT id, original, code, created_at, expiration_date, click_count 
		 FROM urls WHERE code = $1`,
		code,
	).Scan(&u.ID, &u.Original, &u.Code, &u.CreatedAt, &u.ExpirationDate, &u.ClickCount)

	return u, err
}

func (r *Repository) IncrementClickCount(code string) error {
	_, err := r.db.Exec(
		"UPDATE urls SET click_count = click_count + 1 WHERE code = $1",
		code,
	)
	return err
}

func (r *Repository) GetStats(code string) (Url, error) {
	return r.GetByCode(code)
}
