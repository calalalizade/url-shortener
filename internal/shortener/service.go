package shortener

import (
	"context"
	"strings"
	"time"

	"github.com/calalalizade/url-shortener/internal/apperror"
	"github.com/calalalizade/url-shortener/internal/common"
	"github.com/calalalizade/url-shortener/internal/db"
)

type Service struct {
	repo       *Repository
	cache      common.Cache
	maxRetries int
}

func NewService(r *Repository, c common.Cache) *Service {
	return &Service{
		repo:       r,
		cache:      c,
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

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	if s.cache != nil {
		if cachedUrl, err := s.cache.Get(ctx, "code:"+code); err == nil {
			go func() {
				bgCtx := context.Background()
				s.repo.IncrementClickCount(bgCtx, code)
			}()
			return cachedUrl, nil
		}
	}

	url, err := s.repo.GetByCode(code)
	if err != nil {
		return "", &apperror.AppError{
			Type:    apperror.NotFound,
			Message: "resource not found",
		}
	}

	if s.cache != nil {
		s.cache.Set(ctx, "code:"+code, url.Original, 1*time.Minute)
	}

	go func() {
		bgCtx := context.Background()
		s.repo.IncrementClickCount(bgCtx, code)
	}()

	return url.Original, nil
}

func (s *Service) GetStats(code string) (Url, error) {
	url, err := s.repo.GetStats(code)

	if err != nil {
		return Url{}, &apperror.AppError{
			Type:    apperror.NotFound,
			Message: "resource not found",
		}
	}

	return url, nil
}
