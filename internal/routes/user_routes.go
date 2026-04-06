package routes

import (
	"tracking/internal/handler"
	"tracking/internal/middleware"
	"tracking/internal/repository"
	"tracking/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	repo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(rdb)
	svc := service.NewUserService(repo, tokenRepo)
	h := handler.NewUserHandler(svc)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/bootstrap-admin", h.BootstrapAdmin)
	}

	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleBlock("admin"))
	{
		admin.GET("/users", h.ListUsers)
		admin.PATCH("/users/:id/deactivate", h.DeactivateUser)
		admin.PATCH("/drivers/:id/deactivate", h.DeactivateUser)
	}
}
