package routes

import (
	"tracking/internal/handler"
	"tracking/internal/middleware"
	"tracking/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, svc *service.ProductService) {
	h := handler.NewProductHandler(svc)
	productGroup := r.Group("/api/products")
	{
		// Catalogo visible solo para usuarios autenticados.
		productGroup.GET("",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("customer", "driver", "admin"),
			h.GetProductsHandler,
		)
		productGroup.GET("/",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("customer", "driver", "admin"),
			h.GetProductsHandler,
		)

		// Protegido: ABM completo solo para admin.
		productGroup.POST("/",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("admin"),
			h.CreateProductHandler,
		)
		productGroup.PUT("/:id",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("admin"),
			h.UpdateProductHandler,
		)
		productGroup.PATCH("/:id",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("admin"),
			h.UpdateProductHandler,
		)
		productGroup.DELETE("/:id",
			middleware.AuthMiddleware(),
			middleware.RoleBlock("admin"),
			h.DeleteProductHandler,
		)
	}
}
