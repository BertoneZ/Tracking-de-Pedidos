package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"tracking/internal/dto"
	"tracking/internal/service"
	"tracking/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc service.OrderServiceInterface
	locationSvc *service.LocationService
}

// Constructor que usás en routes.go
func NewOrderHandler(oSvc service.OrderServiceInterface, locationSvc *service.LocationService) *OrderHandler {
	return &OrderHandler{svc: oSvc, locationSvc: locationSvc}
}
// Create godoc
// @Summary Crear un nuevo pedido
// @Description Toma la dirección del cliente, busca las coordenadas y guarda el pedido
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "Datos del pedido"
// @Success 201 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
    var req dto.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    userID := c.MustGet("user_id").(string)

    id, err := h.svc.CreateOrder(c.Request.Context(), req, userID)
    if err != nil {
       if errors.Is(err, utils.ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al crear pedido"})
		return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Pedido creado con éxito",
        "id":      id,
    })
}
// GetPending godoc
// @Summary Listar pedidos pendientes
// @Description Obtiene todos los pedidos con estado 'PENDING' disponibles para ser aceptados
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Router /orders/pending [get]
func (h *OrderHandler) GetPending(c *gin.Context) {
	orders, err := h.svc.GetPendingOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener pedidos"})
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
	orderID := c.Param("id") 
	driverID := c.MustGet("user_id").(string) 

	err := h.svc.AcceptOrder(c.Request.Context(), orderID, driverID)
	if err != nil {
		if errors.Is(err, utils.ErrOrderNotAvailable) || errors.Is(err, utils.ErrOrderNotFound){
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
		slog.Error("error en accept", "err", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la solicitud"})
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
    userID := c.MustGet("user_id").(string) // ID del usuario logueado

    // 1. Buscamos la orden
    order, err := h.svc.GetOrderById(c.Request.Context(), orderID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
        return
    }

    // 2. SEGURIDAD: Solo el cliente que creó la orden puede trackearla
    if order.CustomerID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para trackear este pedido"})
        return
    }

    // 3. Verificamos si ya tiene un repartidor
    if order.DriverID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El pedido aún no tiene un repartidor asignado"})
        return
    }

    // 4. Buscamos la ubicación en Redis
    location, err := h.locationSvc.GetLocation(c.Request.Context(), order.DriverID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ubicación no disponible en tiempo real"})
        return
    }

    // 5. LIMPIEZA: Devolvemos solo lat y lng, nada de GeoHash o Dist
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
// @Success 200 {array} dto.OrderResponse
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