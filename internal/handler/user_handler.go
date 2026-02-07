package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"tracking/internal/service"
	_"tracking/internal/dto"
	_"tracking/internal/domain"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}
// Register godoc
// @Summary      Registrar un nuevo usuario
// @Description  Crea un usuario en la DB con email, password, nombre y rol
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterRequest  true  "Datos del registro"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  map[string]string
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var body struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
        FullName string `json:"full_name" binding:"required"` 
        Role     string `json:"role" binding:"required,oneof=driver customer"`
    }

	if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	user, err := h.svc.Register(c.Request.Context(), body.Email, body.Password, body.FullName, body.Role)
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
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email y contraseña requeridos"})
        return
    }

    user, token, err := h.svc.Login(c.Request.Context(), body.Email, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  gin.H{
            "id":        user.ID,
            "email":     user.Email,
            "full_name": user.FullName, 
            "role":      user.Role,
        },
    })
}