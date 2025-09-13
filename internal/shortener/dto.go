package shortener

import "time"

type CreateUrlRequestDTO struct {
	Url string `json:"url"`
}

type UrlResponseDTO struct {
	Code           string    `json:"code"`
	ShortUrl       string    `json:"shortUrl"`
	Original       string    `json:"original"`
	ClickCount     int       `json:"clickCount"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ToUrlResponseDTO(s Url, baseUrl string) UrlResponseDTO {
	return UrlResponseDTO{
		Code:           s.Code,
		ShortUrl:       baseUrl + "/" + s.Code,
		Original:       s.Original,
		ClickCount:     s.ClickCount,
		ExpirationDate: s.ExpirationDate,
		CreatedAt:      s.CreatedAt,
	}
}
