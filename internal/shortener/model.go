package shortener

import "time"

type Url struct {
	ID             int       `json:"id"`
	Code           string    `json:"code"`
	Original       string    `json:"original"`
	ClickCount     int       `json:"clickCount"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreatedAt      time.Time `json:"createdAt"`
}
