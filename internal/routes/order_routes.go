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

func RegisterOrderRoutes(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	// 1. Setup Órdenes (Postgres)
	orderRepo := repository.NewOrderRepository(db, rdb) 
	productRepo := repository.NewProductRepository(db)
    orderSvc := service.NewOrderService(orderRepo, productRepo)

	// 2. Setup Ubicación (Redis)
	locRepo := repository.NewLocationRepository(rdb)
	locSvc := service.NewLocationService(locRepo, orderRepo)

	// 3. Crear el Handler con AMBOS servicios
	h := handler.NewOrderHandler(orderSvc, locSvc)

	orders := r.Group("/api/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/", middleware.RoleBlock("customer"), h.Create)
		
		orders.GET("/pending", middleware.RoleBlock("driver"), h.GetPending)
		orders.PATCH("/:id/accept", middleware.RoleBlock("driver"), h.Accept)
		orders.PATCH("/:id/complete", middleware.RoleBlock("driver"), h.Complete)
		orders.POST("/location", middleware.RoleBlock("driver"), h.UpdateLocation)
		
		orders.GET("/:id/location", h.GetOrderLocation)
		orders.GET("/history", h.GetHistory)
	}
}