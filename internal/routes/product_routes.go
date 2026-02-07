package routes
import (
	"github.com/gin-gonic/gin"
	"tracking/internal/handler"
	"tracking/internal/service"
)


func RegisterProductRoutes(r *gin.Engine, svc *service.ProductService) {
	h := handler.NewProductHandler(svc)
	r.GET("/api/products", h.GetProductsHandler)
}