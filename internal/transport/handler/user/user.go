package user

import (
	"context"
	"encoding/json"
	"nevermore/internal/dto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"nevermore/internal/service"
)

const timeout = 15 * time.Second

type Handler struct {
	srv service.Service
}

func New(srv service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

// @Summary Get user profile
// @Description Get current user profile information
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} user.User "User profile"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /user/get [get]
func (h *Handler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	userId, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user, err := h.srv.User().Get(ctx, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

// @Summary Update user profile
// @Description Update current user profile information with optional photo upload
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param user formData string true "User data in JSON format"
// @Param photo formData file false "Profile photo"
// @Success 200 {object} string "User updated successfully"
// @Failure 400 {object} string "Bad request - invalid data"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /user/update [put]
func (h *Handler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	id, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// Парсим multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); // 10 MB limit
	err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Получаем данные пользователя из формы
	userJSON := c.PostForm("user")
	if userJSON == "" {
		c.JSON(400, gin.H{"error": "User data is required"})
		return
	}

	var userData dto.UpdateUserRequest
	if err := json.Unmarshal([]byte(userJSON), &userData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user data"})
		return
	}

	// Получаем файл фото

	err = h.srv.User().Update(ctx, id, userData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// @Summary Delete user account
// @Description Delete current user account
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} string "User deleted successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /user/delete [delete]
func (h *Handler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	userId, err := strconv.Atoi(userIDStr)

	err = h.srv.User().Delete(ctx, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
