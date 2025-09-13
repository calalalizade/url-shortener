package db

import (
	"errors"

	"github.com/lib/pq"
)

// IsDuplicateKeyError checks for Postgres unique constraint violation
func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
