package middleware

import (
	"net/http"
	"strings"
	"tracking/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere token"})
			c.Abort() 
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		
		token, err := auth.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user_id", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))

		c.Next() 
	}
}