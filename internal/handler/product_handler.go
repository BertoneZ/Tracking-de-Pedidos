package handler

import (
	"net/http"
	"tracking/internal/service"
	"github.com/gin-gonic/gin"
	"tracking/internal/dto"
	
)


type ProductHandler struct {
	svc service.ProductServiceInterface
}

func NewProductHandler(svc service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// GetProducts godoc
// @Summary Obtener el menú de productos
// @Tags Products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Router /products [get]
func (h *ProductHandler) GetProductsHandler(c *gin.Context) {
	products, err := h.svc.GetProductsService(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}
	if products == nil {
        products = []dto.ProductResponse{}
    }
	c.JSON(http.StatusOK, products)
}