package shortener

import "time"

type ShortUrl struct {
	Code           string    `json:"code"`
	Original       string    `json:"original"`
	ClickCount     int       `json:"clickCount"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreatedAt      time.Time `json:"createdAt"`
}
