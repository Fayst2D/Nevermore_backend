package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nevermore/internal/dto"

	"nevermore/internal/service"
	_ "nevermore/pkg/auth"
)

type Handler struct {
	srv service.Service
}

func New(srv service.Service) *Handler {
	result := &Handler{
		srv: srv,
	}
	return result
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} dto.MessageResponse "Успешная регистрация"
// @Failure 400 {object} dto.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.srv.Auth().Register(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Учетные данные"
// @Success 200 {object} auth.Token "Успешная аутентификация"
// @Failure 400 {object} dto.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.srv.Auth().Login(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary Выход пользователя
// @Description Деактивирует токен пользователя
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.MessageResponse "Successfully logged out"
// @Failure 500 {object} dto.ErrorResponse "Failed to log out"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	email, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	err := h.srv.Auth().Logout(c.Request.Context(), email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Refresh godoc
// @Summary Обновление токена
// @Description Обновляет access токен с помощью refresh токена
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh токен"
// @Success 200 {object} Token "Новые access и refresh токены"
// @Failure 400 {object} dto.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.srv.Auth().Refresh(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, token)
}
