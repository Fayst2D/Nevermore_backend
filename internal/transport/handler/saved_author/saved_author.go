package saved_author

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

// @Summary Save author
// @Description Save author to user's saved list
// @Tags saved_authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SaveAuthorRequest true "Save author request"
// @Success 200 {object} string "Author saved successfully"
// @Failure 400 {object} string "Bad request - invalid data"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /saved-author/create [post]
func (h *Handler) Create(c *gin.Context) {
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

	var request dto.SaveAuthorRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	request.UserId = userId

	err = h.srv.SavedAuthor().Create(ctx, request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Author saved successfully"})
}

// @Summary Delete saved author
// @Description Remove author from user's saved list
// @Tags saved_authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param author_id query int true "Author ID"
// @Success 200 {object} string "Author removed successfully"
// @Failure 400 {object} string "Bad request - invalid author ID"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /saved-author/delete [delete]
func (h *Handler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	authorIDStr := c.Query("author_id")
	if authorIDStr == "" {
		c.JSON(400, gin.H{"error": "Author ID is required"})
		return
	}

	authorId, err := strconv.Atoi(authorIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid author ID"})
		return
	}

	err = h.srv.SavedAuthor().Delete(ctx, authorId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Author removed successfully"})
}

// @Summary Get saved authors list
// @Description Get list of user's saved authors
// @Tags saved_authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AuthorGetResponse "List of saved authors"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /saved-author/list [get]
func (h *Handler) GetList(c *gin.Context) {
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

	authors, err := h.srv.SavedAuthor().GetSavedAuthorsList(ctx, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, authors)
}
