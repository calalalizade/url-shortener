package shortener

import (
	"net/url"
	"strings"

	"github.com/calalalizade/url-shortener/internal/apperror"
)

func validateURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return "", &apperror.AppError{
			Type:    apperror.Validation,
			Message: "url cannot be empty",
		}
	}

	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}

	parsedURL, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", &apperror.AppError{
			Type:    apperror.Validation,
			Message: "invalid url format",
		}
	}

	if parsedURL.Host == "" {
		return "", &apperror.AppError{
			Type:    apperror.Validation,
			Message: "url must have a valid host",
		}
	}

	return raw, nil
}
