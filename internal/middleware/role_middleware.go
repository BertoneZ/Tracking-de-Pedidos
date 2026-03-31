package middleware

import (
	"net/http"
	"strings"

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

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "rol invalido en la sesión"})
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if strings.EqualFold(role, roleStr) {
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