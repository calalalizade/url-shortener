package shortener

import (
	"net/http"

	"github.com/calalalizade/url-shortener/internal/apperror"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	baseUrl string
}

func NewHandler(s *Service, baseUrl string) *Handler {
	return &Handler{
		service: s,
		baseUrl: baseUrl,
	}
}

func (h *Handler) Shorten(c *gin.Context) error {
	var req CreateUrlRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		return &apperror.AppError{
			Type:    apperror.Validation,
			Message: "invalid request body",
		}
	}

	url, err := h.service.ShortenUrl(req.Url)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, ToUrlResponseDTO(url, h.baseUrl))
	return nil
}

func (h *Handler) Resolve(c *gin.Context) error {
	code := c.Param("code")

	url, err := h.service.Resolve(c.Request.Context(), code)
	if err != nil {
		return err
	}

	c.Redirect(http.StatusMovedPermanently, url)

	return nil
}

func (h *Handler) GetStats(c *gin.Context) error {
	code := c.Param("code")

	u, err := h.service.GetStats(code)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, ToUrlResponseDTO(u, h.baseUrl))
	return nil
}
