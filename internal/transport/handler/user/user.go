package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nevermore/internal/dto"
	"nevermore/pkg/logger"
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

// Get godoc
// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.User
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /user [get]
func (h *Handler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log := logger.Get()
	log.Info().Msg(c.Request.RequestURI)

	fmt.Println(c.Request.RequestURI)

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

// Update godoc
// @Summary Update user profile
// @Description Update current authenticated user's profile information and photo
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param user formData string true "User data in JSON format"
// @Param photo formData file false "Profile photo"
// @Success 200 {object} dto.MessageResponse "message: User updated successfully"
// @Failure 400 {object} dto.ErrorResponse "error: Bad request"
// @Failure 401 {object} dto.ErrorResponse "error: Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "error: Internal server error"
// @Router /user [put]
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

	userId, err := strconv.Atoi(userIDStr)

	if err := c.Request.ParseMultipartForm(10 << 20); // 10 MB limit
		err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form data"})
		return
	}

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

	file, header, err := c.Request.FormFile("photo")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(400, gin.H{"error": "Failed to get photo"})
		return
	}

	var photo dto.FileInfo
	if file != nil {
		defer file.Close()
		photo = dto.FileInfo{
			File:   file,
			Header: header,
		}
	}

	log := logger.Get()
	log.Info().Msg("CALL UPDATE")

	err = h.srv.User().Update(ctx, userId, userData, photo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// Delete godoc
// @Summary Delete user account
// @Description Delete current authenticated user's account (soft delete)
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse "message: User deleted successfully"
// @Failure 401 {object} dto.ErrorResponse "error: Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "error: Internal server error"
// @Router /user [delete]
func (h *Handler) Delete(c *gin.Context) {
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

	err = h.srv.User().Delete(ctx, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
