package book

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"nevermore/internal/dto"
	"strconv"
	"time"

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

func (h *Handler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	id, exists := c.Get("userID")
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

	// Получаем данные из отдельных полей формы
	title := c.PostForm("title")
	author := c.PostForm("author")
	description := c.PostForm("description")

	if title == "" {
		c.JSON(400, gin.H{"error": "Title is required"})
		return
	}
	if author == "" {
		c.JSON(400, gin.H{"error": "Author is required"})
		return
	}

	// Создаем запрос
	req := dto.CreateBookRequest{
		Title:       title,
		Author:      author,
		Description: &description, // используем указатель
	}

	// Преобразуем ID пользователя
	uploadedBy, err := strconv.Atoi(id.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	req.UploadedBy = uploadedBy

	// Получаем файл
	file, header, err := c.Request.FormFile("File") // "File" - имя поля из Postman
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(400, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}

	var fileInfo dto.FileInfo
	if file != nil {
		defer file.Close()
		fileInfo = dto.FileInfo{
			File:   file,
			Header: header,
		}
	}

	err = h.srv.Book().Create(ctx, &req, fileInfo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Book created successfully"})
}
