package shortener

import (
	"strings"

	"github.com/calalalizade/url-shortener/internal/apperror"
	"github.com/calalalizade/url-shortener/internal/db"
)

type Service struct {
	repo       *Repository
	maxRetries int
}

func NewService(r *Repository) *Service {
	return &Service{
		repo:       r,
		maxRetries: 5,
	}
}

func (s *Service) ShortenUrl(url string) (Url, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return Url{}, &apperror.AppError{
			Type:    apperror.Validation,
			Message: "url cannot be empty",
		}
	}

	for i := 0; i < s.maxRetries; i++ {
		code, err := GenerateCode()
		if err != nil {
			return Url{}, err
		}

		u, err := s.repo.Create(url, code)
		if err == nil {
			return u, nil
		}

		if db.IsDuplicateKeyError(err) {
			continue
		}

		return Url{}, err
	}

	return Url{}, &apperror.AppError{
		Type:    apperror.Internal,
		Message: "failed to generate a unique short URL",
	}
}

func (s *Service) GetOriginalUrl(code string) (Url, error) {
	url, err := s.repo.GetByCode(code)
	if err != nil {
		return Url{}, &apperror.AppError{
			Type:    apperror.NotFound,
			Message: "resource not found",
		}
	}

	return url, nil
}
