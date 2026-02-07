package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RoleBlock(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {		
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "rol no encontrado en la sesión"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		isAllowed := false
		for _, role := range allowedRoles {
			if role == roleStr {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "no tienes permisos para realizar esta acción"})
			c.Abort()
			return
		}

		c.Next()
	}
}