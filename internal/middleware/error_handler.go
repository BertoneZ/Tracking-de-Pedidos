package middleware

import (
	"net/http"
	"tracking/internal/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware es un middleware que captura panics y errores globalmente
// Proporciona una respuesta consistente para todos los errores de la aplicación
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
					Code:       "INTERNAL_ERROR",
					Message:    "Error interno del servidor",
					StatusCode: http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()

		c.Next()

	
	}
}

// HandleError es una función auxiliar para responder con errores de forma consistente
func HandleError(c *gin.Context, appErr *utils.AppError) {
	c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
}
