package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"tracking/internal/service"
	_"tracking/internal/dto"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=driver customer"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), body.Email, body.Password, body.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear usuario"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
// Login godoc
// @Summary Iniciar sesión
// @Description Devuelve un token JWT si las credenciales son válidas
// @Tags Auth
// @Accept json
// @Produce json
//@Param credentials body dto.LoginRequest true "Credenciales de usuario"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "credenciales incorrectas"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}