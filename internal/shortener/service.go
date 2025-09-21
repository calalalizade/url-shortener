package shortener

import (
	"context"
	"fmt"
	"time"

	"github.com/calalalizade/url-shortener/internal/apperror"
	"github.com/calalalizade/url-shortener/internal/common"
	"github.com/calalalizade/url-shortener/internal/db"
	"github.com/calalalizade/url-shortener/internal/platform"
)

type Service struct {
	repo        *Repository
	cache       common.Cache
	cacheConfig platform.CacheConfig
	maxRetries  int
}

func NewService(r *Repository, c common.Cache, cacheConfig platform.CacheConfig) *Service {
	return &Service{
		repo:        r,
		cache:       c,
		cacheConfig: cacheConfig,
		maxRetries:  5,
	}
}

// POST /shorten
func (s *Service) ShortenUrl(ctx context.Context, url string) (Url, error) {
	valUrl, err := validateURL(url)
	if err != nil {
		return Url{}, err
	}

	// check cache first
	if urlObj, found := s.tryGetFromCache(ctx, valUrl); found {
		return urlObj, nil
	}

	// check db
	if urlObj, err := s.repo.GetByOriginalUrl(valUrl); err == nil {
		s.cacheUrl(ctx, urlObj)
		return urlObj, nil
	}

	return s.createShortenedUrl(ctx, valUrl)
}

func (s *Service) tryGetFromCache(ctx context.Context, originalUrl string) (Url, bool) {
	if s.cache == nil || !s.cacheConfig.Enabled {
		return Url{}, false
	}

	cacheKey := "url:" + originalUrl
	cachedCode, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return Url{}, false
	}

	urlObj, err := s.repo.GetByCode(cachedCode)
	if err != nil {
		return Url{}, false
	}

	return urlObj, true
}

func (s *Service) cacheUrl(ctx context.Context, urlObj Url) {
	if s.cache == nil || !s.cacheConfig.Enabled {
		return
	}

	ttl := s.getCacheTTL(urlObj)
	if ttl <= 0 {
		return
	}

	s.cache.Set(ctx, "url:"+urlObj.Original, urlObj.Code, ttl)
	s.cache.Set(ctx, "code:"+urlObj.Code, urlObj.Original, ttl)
}

func (s *Service) createShortenedUrl(ctx context.Context, originalUrl string) (Url, error) {
	code := GenerateCodeFromURL(originalUrl)

	for i := 0; i < s.maxRetries; i++ {
		urlObj, err := s.repo.Create(originalUrl, code)
		if err == nil {
			s.cacheUrl(ctx, urlObj)
			return urlObj, nil
		}

		if !db.IsDuplicateKeyError(err) {
			return Url{}, err
		}

		// Handle collision
		existingUrl, err := s.repo.GetByCode(code)
		if err == nil && existingUrl.Original == originalUrl {
			s.cacheUrl(ctx, existingUrl)
			return existingUrl, nil
		}

		// generate new code with salt
		code = GenerateCodeFromURL(originalUrl + fmt.Sprintf("_salt_%d", i))
	}

	return Url{}, &apperror.AppError{
		Type:    apperror.Internal,
		Message: "failed to generate a unique short URL",
	}
}

// GET /:code
func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	if s.cache != nil && s.cacheConfig.Enabled {
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

	if time.Now().After(url.ExpirationDate) {
		return "", &apperror.AppError{
			Type:    apperror.NotFound,
			Message: "URL has expired",
		}
	}

	if s.cache != nil && s.cacheConfig.Enabled {
		ttl := s.getCacheTTL(url)

		if ttl > 0 {
			s.cache.Set(ctx, "code:"+code, url.Original, ttl)

			s.cache.Set(ctx, "url:"+url.Original, code, ttl)
		}
	}

	go func() {
		bgCtx := context.Background()
		s.repo.IncrementClickCount(bgCtx, code)
	}()

	return url.Original, nil
}

// GET /:code/stats
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

func (s *Service) getCacheTTL(url Url) time.Duration {
	if !s.cacheConfig.Enabled {
		return 0
	}

	var baseTTL time.Duration

	// Based on click count -> determine tier
	if url.ClickCount >= s.cacheConfig.HotThreshold {
		baseTTL = s.cacheConfig.HotURLTTL
	} else if url.ClickCount >= s.cacheConfig.WarmThreshold {
		baseTTL = s.cacheConfig.WarmURLTTL
	} else {
		baseTTL = s.cacheConfig.ColdURLTTL
	}

	timeUntilExpiry := time.Until(url.ExpirationDate)
	if timeUntilExpiry < baseTTL {
		baseTTL = timeUntilExpiry
	}

	if baseTTL > s.cacheConfig.MaxTTL {
		baseTTL = s.cacheConfig.MaxTTL
	}

	if baseTTL < s.cacheConfig.MinTTL {
		return 0
	}

	return baseTTL
}
