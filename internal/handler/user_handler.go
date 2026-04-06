package handler

import (
	"net/http"
	_ "tracking/internal/domain"
	_ "tracking/internal/dto"
	"tracking/internal/service"
	"tracking/internal/utils"

	"github.com/gin-gonic/gin"
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
		validationErr := utils.HandleValidationErrors(err)
		c.JSON(validationErr.StatusCode, validationErr.ToErrorResponse())
		return
	}

	user, err := h.svc.Register(c.Request.Context(), body.Email, body.Password, body.FullName, body.Role)
	if err != nil {
		appErr := utils.NewAppError("REGISTRATION_ERROR", "Error al crear usuario", 500, err)
		c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
		return
	}

	c.JSON(201, user)
}

// Login godoc
// @Summary Iniciar sesión
// @Description Devuelve un token JWT si las credenciales son válidas
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Credenciales de usuario"
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

	user, accessToken, refreshToken, err := h.svc.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// RefreshToken godoc
// @Summary Renovar tokens
// @Description Devuelve un nuevo access token y refresh token válidos
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token requerido"})
		return
	}

	_, accessToken, refreshToken, err := h.svc.Refresh(c.Request.Context(), body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token inválido o expirado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}

// BootstrapAdmin godoc
// @Summary Crear primer admin (bootstrap)
// @Description Crea el primer admin del sistema. Requiere el secret de bootstrap y falla si ya existe un admin.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body dto.BootstrapAdminRequest true "Datos para crear el admin inicial"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /auth/bootstrap-admin [post]
func (h *UserHandler) BootstrapAdmin(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Secret   string `json:"secret" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	user, err := h.svc.BootstrapAdmin(c.Request.Context(), body.Email, body.Password, body.FullName, body.Secret)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers godoc
// @Summary Listar usuarios
// @Description Devuelve todos los usuarios del sistema. Solo ADMIN. Soporta filtros por role y active.
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param role query string false "Filtrar por role (driver, customer, admin)"
// @Param active query boolean false "Filtrar por estado activo (true/false)"
// @Success 200 {array} domain.User
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	activeStr := c.Query("active")

	var active *bool
	if activeStr != "" {
		activeVal := activeStr == "true"
		active = &activeVal
	}

	users, err := h.svc.ListUsers(c.Request.Context(), role, active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar usuarios"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeactivateUser godoc
// @Summary Desactivar usuario
// @Description Realiza soft delete de un usuario (is_active=false). Solo ADMIN.
// @Tags Admin
// @Security BearerAuth
// @Param id path string true "ID del usuario"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/users/{id}/deactivate [patch]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	actorUserID := c.GetString("user_id")
	if actorUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id requerido"})
		return
	}

	err := h.svc.DeactivateUser(c.Request.Context(), actorUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuario desactivado correctamente"})
}
