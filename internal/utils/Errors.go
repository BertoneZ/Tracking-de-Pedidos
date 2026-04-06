package utils

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	ErrEmptyOrder          = errors.New("el pedido debe contener al menos un producto")
	ErrOrderNotFound       = errors.New("pedido no encontrado")
	ErrOrderNotAvailable   = errors.New("este pedido ya no está disponible")
	ErrInvalidAddress      = errors.New("dirección de destino inválida")
	ErrInvalidQuantity     = errors.New("cantidad inválida para un producto")
	ErrProductNotFound     = errors.New("producto no encontrado")
	ErrDeliveryNotFinished = errors.New("no se puede iniciar un nuevo pedido sin finalizar el actual")
	ErrInternal            = errors.New("error interno del servidor")
	ErrUnauthorizedAction  = errors.New("no tienes permisos para realizar esta acción")
	ErrInvalidState        = errors.New("la acción no es válida para el estado actual del pedido")
)

// ErrorResponse es la estructura estándar para todas las respuestas de error
type ErrorResponse struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"status_code"`
}

// AppError representa un error de la aplicación con contexto adicional
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    map[string]string
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError crea un nuevo error de aplicación
func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]string),
		Err:        err,
	}
}

// ValidationError crea un error de validación con detalles de campos
func ValidationError(details map[string]string) *AppError {
	appErr := &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Error en la validación de datos",
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
	return appErr
}

// ToErrorResponse convierte un AppError a ErrorResponse
func (e *AppError) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Details:    e.Details,
	}
}

// HandleValidationErrors convierte los errores de validación de Gin a un AppError
func HandleValidationErrors(err error) *AppError {
	messages := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			fieldName := fe.Field()
			tag := fe.Tag()

			var message string
			switch tag {
			case "required":
				message = fieldName + " es requerido"
			case "email":
				message = "Debe ser un email válido"
			case "min":
				message = "Mínimo " + fe.Param() + " caracteres"
			case "max":
				message = "Máximo " + fe.Param() + " caracteres"
			case "oneof":
				message = "Debe ser uno de: " + fe.Param()
			default:
				message = "Campo inválido"
			}
			messages[fieldName] = message
		}
		return ValidationError(messages)
	}

	// Fallback para errores que no son de validación
	return ValidationError(map[string]string{"error": "Error en la estructura del JSON"})
}
