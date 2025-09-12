package shortener

import "time"

type CreateShortenerRequest struct {
	Url string `json:"url" binding:"required"`
}

type ShortenerResponse struct {
	Code       string    `json:"code"`
	Original   string    `json:"original"`
	ClickCount int       `json:"clickCount"`
	CreatedAt  time.Time `json:"createdAt"`
}

func ToShortenerResponse(s ShortUrl) ShortenerResponse {
	return ShortenerResponse{
		Code:       s.Code,
		Original:   s.Original,
		ClickCount: s.ClickCount,
		CreatedAt:  s.CreatedAt,
	}
}
