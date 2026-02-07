package routes

import (
	"tracking/internal/handler"
	"tracking/internal/repository"
	"tracking/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *gin.Engine, db *pgxpool.Pool) {
	
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}