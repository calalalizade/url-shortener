package shortener

import (
	"math/rand/v2"
	"strconv"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ShortenUrl(url string) (string, error) {
	ext := rand.IntN(999999)
	str := strconv.Itoa(ext)
	return str, nil
}
