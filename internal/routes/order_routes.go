package routes

import (
	"tracking/internal/handler"
	"tracking/internal/repository"
	"tracking/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9" 
	"tracking/internal/middleware"
)

// Agregá el tercer parámetro rdb aquí:
func RegisterOrderRoutes(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	// 1. Setup Órdenes (Postgres)
	orderRepo := repository.NewOrderRepository(db, rdb) 
    orderSvc := service.NewOrderService(orderRepo)

	// 2. Setup Ubicación (Redis)
	locRepo := repository.NewLocationRepository(rdb)
	locSvc := service.NewLocationService(locRepo)

	// 3. Crear el Handler con AMBOS servicios
	h := handler.NewOrderHandler(orderSvc, locSvc)

	orders := r.Group("/api/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/", h.Create)
		orders.GET("/pending", h.GetPending)
		orders.PATCH("/:id/accept", h.Accept)
		orders.POST("/location", h.UpdateLocation) // El endpoint para el GPS
		// Dentro de RegisterOrderRoutes
		orders.GET("/:id/location", h.GetOrderLocation)
		// Dentro de RegisterOrderRoutes
		orders.PATCH("/:id/complete", h.Complete)
		// Dentro del grupo de órdenes
		orders.GET("/history", h.GetHistory)

	}
}