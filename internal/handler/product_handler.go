package handler

import (
	"errors"
	"net/http"
	"tracking/internal/dto"
	"tracking/internal/service"
	"tracking/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
// @Security BearerAuth
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

// CreateProduct godoc
// @Summary Crear producto
// @Description Crea un nuevo producto del catalogo. Solo ADMIN.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body dto.UpsertProductRequest true "Datos del producto"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) CreateProductHandler(c *gin.Context) {
	var req dto.UpsertProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	product, err := h.svc.CreateProductService(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Actualizar producto
// @Description Actualiza un producto existente por ID. Solo ADMIN.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID del producto"
// @Param product body dto.UpdateProductRequest true "Campos a actualizar del producto"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProductHandler(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	product, err := h.svc.UpdateProductService(c.Request.Context(), productID, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Eliminar producto
// @Description Elimina un producto por ID. Solo ADMIN.
// @Tags Products
// @Security BearerAuth
// @Param id path string true "ID del producto"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProductHandler(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	err := h.svc.DeleteProductService(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	c.Status(http.StatusNoContent)
}
