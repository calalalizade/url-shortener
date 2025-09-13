package shortener

import "time"

type CreateUrlRequestDTO struct {
	Url string `json:"url" binding:"required"`
}

type UrlResponseDTO struct {
	Code           string    `json:"code"`
	Original       string    `json:"original"`
	ClickCount     int       `json:"clickCount"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ToUrlResponseDTO(s Url) UrlResponseDTO {
	return UrlResponseDTO{
		Code:           s.Code,
		Original:       s.Original,
		ClickCount:     s.ClickCount,
		ExpirationDate: s.ExpirationDate,
		CreatedAt:      s.CreatedAt,
	}
}
