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
		// 1. Extraer el header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere token"})
			c.Abort() // Corta la ejecución aquí
			return
		}

		// 2. Limpiar el formato "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		
		// 3. Validar el token usando la lógica de jwt.go
		token, err := auth.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// 4. Inyectar los datos del usuario en el contexto para que los handlers los usen
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user_id", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))

		c.Next() // Permite que la petición siga al siguiente Handler
	}
}