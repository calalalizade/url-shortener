package shortener

import (
	"net/http"

	"github.com/calalalizade/url-shortener/internal/apperror"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/shorten", wrap(h.Shorten))
	r.GET("/:code", wrap(h.GetByCode))
}

func wrap(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		switch appErr.Type {
		case apperror.Validation:
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Message})
		case apperror.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
