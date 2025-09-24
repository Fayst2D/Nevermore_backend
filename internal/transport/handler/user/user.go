package user

import (
	"context"
	"encoding/json"
	"nevermore/internal/model/user"
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
	}

	user, err := h.srv.User().Get(ctx, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
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

	var userData user.User
	if err := json.Unmarshal([]byte(userJSON), &userData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user data"})
		return
	}

	// Получаем файл фото

	err := h.srv.User().Update(ctx, userData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// Delete удаляет пользователя
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
