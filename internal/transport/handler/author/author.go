package author

import (
	"context"
	"encoding/json"
	"nevermore/internal/model/author"
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

// @Summary Get author
// @Description Get selected author information
// @Tags authors
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} user.User "Author information"
// @Failure 404 {object} string "Author not found"
// @Failure 500 {object} string "Internal server error"
// @Router /author/get/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	authorIdStr := c.Param("id")
	if authorIdStr == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	author, err := h.srv.Author().Get(ctx, authorId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, author)
}

// @Summary Get authors list
// @Description Get list of authors
// @Tags authors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.AuthorGetResponse "List of authors"
// @Failure 500 {object} string "Internal server error"
// @Router /author/list [get]
func (h *Handler) GetAuthorsList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	authors, err := h.srv.Author().GetAuthorsList(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, authors)
}

// @Summary Update author information
// @Description Update selected author information with optional photo upload
// @Tags authors
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param author formData string true "Author data in JSON format"
// @Param photo formData file false "photo"
// @Success 200 {object} string "Author updated successfully"
// @Failure 404 {object} string "Author not found"
// @Failure 400 {object} string "Bad request - invalid data"
// @Failure 500 {object} string "Internal server error"
// @Router /author/update/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	authorIdStr := c.Param("id")
	if authorIdStr == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	// Парсим multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); // 10 MB limit
	err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Получаем данные автора из формы
	authorJSON := c.PostForm("author")
	if authorJSON == "" {
		c.JSON(400, gin.H{"error": "Author data is required"})
		return
	}

	var authorData author.Author
	if err := json.Unmarshal([]byte(authorJSON), &authorData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid author data"})
		return
	}

	// Получаем файл фото

	err := h.srv.Author().Update(ctx, authorData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Author updated successfully"})
}

// @Summary Delete author
// @Description Delete selected author
// @Tags authors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Author ID"
// @Success 200 {object} string "Author deleted successfully"
// @Failure 404 {object} string "Author not found"
// @Failure 500 {object} string "Internal server error"
// @Router /author/delete/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	authorIdStr := c.Param("id")
	if authorIdStr == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	authorId, err := strconv.Atoi(authorIdStr)

	err = h.srv.Author().Delete(ctx, authorId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Author deleted successfully"})
}
