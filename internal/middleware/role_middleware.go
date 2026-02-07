package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RoleBlock recibe una lista de roles permitidos (variadic parameter)
func RoleBlock(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener el rol que inyectó el AuthMiddleware
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "rol no encontrado en la sesión"})
			c.Abort()
			return
		}

		// 2. Verificar si el rol del usuario está en la lista de permitidos
		roleStr := userRole.(string)
		isAllowed := false
		for _, role := range allowedRoles {
			if role == roleStr {
				isAllowed = true
				break
			}
		}

		// 3. Si no está permitido, cortar la ejecución
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "no tienes permisos para realizar esta acción"})
			c.Abort()
			return
		}

		c.Next()
	}
}