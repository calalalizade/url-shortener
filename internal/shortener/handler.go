package shortener

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Shorten(c *gin.Context) error {
	var req CreateShortenerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	}

	code, err := h.service.ShortenUrl(req.Url)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, code)
	return nil
}
