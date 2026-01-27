package handler

import (
	"net/http"
	"tracking/internal/dto"
	_"tracking/internal/domain"
	"tracking/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
	locationSvc *service.LocationService
}

// Constructor que usás en routes.go
func NewOrderHandler(oSvc *service.OrderService, locationSvc *service.LocationService) *OrderHandler {
	return &OrderHandler{svc: oSvc, locationSvc: locationSvc}
}
// CreateOrder godoc
// @Summary Crear un nuevo pedido
// @Description Registra un pedido con origen y destino
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "Datos del pedido"
// @Success 201 {object} domain.Order
// @Router /orders [post]
// @Security BearerAuth
func (h *OrderHandler) Create(c *gin.Context) {
	var body dto.CreateOrderRequest

	// 1. Validar el JSON de entrada usando el DTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de orden inválidos"})
		return
	}

	// 2. Obtener el ID del cliente desde el Middleware (JWT)
	// c.MustGet trae el valor que guardamos con c.Set() en el middleware
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar al usuario"})
		return
	}

	// 3. Llamar al servicio para crear la orden
	order, err := h.svc.CreateOrder(c.Request.Context(), customerID.(string), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el pedido"})
		return
	}

	// 4. Devolver la orden creada
	c.JSON(http.StatusCreated, order)
}
// GetPending godoc
// @Summary Listar pedidos pendientes
// @Description Obtiene todos los pedidos con estado 'PENDING' disponibles para ser aceptados
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Router /orders/pending [get]
func (h *OrderHandler) GetPending(c *gin.Context) {
	orders, err := h.svc.GetPendingOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener pedidos"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, []interface{}{}) // Devolvemos lista vacía si no hay nada
		return
	}

	c.JSON(http.StatusOK, orders)
}
// AcceptOrder godoc
// @Summary Aceptar un pedido (Driver)
// @Description Cambia el estado del pedido a ASSIGNED
// @Tags Orders
// @Security BearerAuth
// @Param id path string true "ID del pedido"
// @Success 200 {object} map[string]string
// @Router /orders/{id}/accept [patch]
func (h *OrderHandler) Accept(c *gin.Context) {
	orderID := c.Param("id") // Saca el ID de la URL: /api/orders/:id/accept
	driverID := c.MustGet("user_id").(string) // Lo sacamos del Token (Middleware)

	err := h.svc.AcceptOrder(c.Request.Context(), orderID, driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pedido aceptado con éxito"})
}
// UpdateLocation godoc
// @Summary Actualizar GPS (Driver)
// @Description Guarda la ubicación actual en Redis
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Param location body dto.UpdateLocationRequest true "Coordenadas"
// @Success 200 {object} map[string]string
// @Router /orders/location [post]
func (h *OrderHandler) UpdateLocation(c *gin.Context) {
	var body dto.UpdateLocationRequest

	// 1. Validar que el JSON tenga lat y lng
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coordenadas requeridas"})
		return
	}

	// 2. Obtener el ID del Driver desde el Contexto (inyectado por el AuthMiddleware)
	driverID := c.MustGet("user_id").(string)

	// 3. Llamar al servicio de ubicación para guardar en Redis
	err := h.locationSvc.UpdateLocation(c.Request.Context(), driverID, body.Lat, body.Lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar ubicación en tiempo real"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ubicación actualizada correctamente"})
}
// GetOrderLocation godoc
// @Summary Consultar ubicación de un pedido (Cliente)
// @Description Obtiene la última posición registrada en Redis del driver asignado a la orden
// @Tags Orders
// @Security BearerAuth
// @Param id path string true "ID del pedido"
// @Produce json
// @Success 200 {object} map[string]float64
// @Router /orders/{id}/location [get]
func (h *OrderHandler) GetOrderLocation(c *gin.Context) {
	orderID := c.Param("id")

	// 1. Buscar la orden en Postgres para saber quién es el driver
	// (Podés usar un método del orderSvc para esto)
	order, err := h.svc.GetOrderById(c.Request.Context(), orderID)
	if err != nil || order.DriverID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "El pedido no tiene un repartidor asignado"})
		return
	}

	// 2. Buscar la ubicación actual del driver en Redis
	location, err := h.locationSvc.GetLocation(c.Request.Context(), order.DriverID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Devolver latitud y longitud al cliente
	c.JSON(http.StatusOK, gin.H{
		"lat": location.Latitude,
		"lng": location.Longitude,
	})
}
// Complete godoc
// @Summary Finalizar entrega (Driver)
// @Description Cambia el estado a 'DELIVERED' en Postgres y elimina la ubicación de Redis
// @Tags Orders
// @Security BearerAuth
// @Param id path string true "ID del pedido"
// @Success 200 {object} map[string]string
// @Router /orders/{id}/complete [patch]
func (h *OrderHandler) Complete(c *gin.Context) {
	orderID := c.Param("id")
	driverID := c.MustGet("user_id").(string)

	err := h.svc.CompleteOrder(c.Request.Context(), orderID, driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "¡Pedido entregado con éxito!"})
}
// GetHistory godoc
// @Summary Ver historial de pedidos
// @Description Trae todos los pedidos DELIVERED del usuario
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Router /orders/history [get]
func (h *OrderHandler) GetHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	orders, err := h.svc.GetUserHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener historial"})
		return
	}

	c.JSON(http.StatusOK, orders)
}