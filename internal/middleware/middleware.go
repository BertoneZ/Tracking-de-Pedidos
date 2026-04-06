package middleware

import (
	"net/http"
	"strings"
	"tracking/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere token"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de Authorization inválido"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := auth.ValidateAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		if claims.UserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			c.Abort()
			return
		}

		if claims.Role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
