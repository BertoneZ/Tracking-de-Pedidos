package utils

import "errors"

var (
	ErrEmptyOrder       = errors.New("el pedido debe contener al menos un producto")
	ErrOrderNotFound	 = errors.New("pedido no encontrado")
	ErrOrderNotAvailable = errors.New("este pedido ya no está disponible")
	ErrInvalidAddress  = errors.New("dirección de destino inválida")
	ErrInvalidQuantity = errors.New("cantidad inválida para un producto")
	ErrProductNotFound = errors.New("producto no encontrado")
	ErrDeliveryNotFinished = errors.New("no se puede iniciar un nuevo pedido sin finalizar el actual")
	ErrInternal         = errors.New("error interno del servidor")
	ErrUnauthorizedAction = errors.New("no tienes permisos para realizar esta acción")
	ErrInvalidState = errors.New("la acción no es válida para el estado actual del pedido")
)
