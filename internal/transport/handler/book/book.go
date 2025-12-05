package book

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"nevermore/internal/dto"
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

// @Summary Создание новой книги
// @Description Создает новую книгу с возможностью загрузки файла
// @Tags books
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Название книги" minlength(1)
// @Param author formData string true "Автор книги" minlength(1)
// @Param description formData string false "Описание книги"
// @Param File formData file false "Файл книги (PDF, EPUB, etc.)"
// @Success 200 {object} map[string]string "Успешное создание книги"
// @Success 201 {object} map[string]string "Книга создана"
// @Failure 400 {object} map[string]string "Некорректные входные данные"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 413 {object} map[string]string "Превышен максимальный размер файла"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /books [post]
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

// @Summary Получение книг по автору
// @Description Возвращает список книг указанного автора по его ID
// @Tags books
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param author_id path integer true "ID автора" minimum(1)
// @Success 200 {object} map[string]interface{} "Успешный ответ с книгами автора"
// @Success 200 {object} map[string]interface{} "Список книг автора"
// @Failure 400 {object} map[string]string "Некорректный ID автора или параметр отсутствует"
// @Failure 404 {object} map[string]interface{} "Книги не найдены для указанного автора"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /book/by-author/{author_id} [get]
func (h *Handler) GetByAuthor(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Получаем ID автора из query параметра
	authorIDStr := c.Param("author_id")
	if authorIDStr == "" {
		c.JSON(400, gin.H{"error": "author_id path parameter is required"})
		return
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "author_id must be a valid integer"})
		return
	}

	// Вызываем сервис для получения книг автора
	books, err := h.srv.Book().GetByAuthor(ctx, authorID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(books) == 0 {
		c.JSON(404, gin.H{"message": "No books found for this author", "books": []dto.GetBookRequest{}})
		return
	}

	c.JSON(200, gin.H{"books": books})
}

// @Summary Поиск книг по названию
// @Description Регистронезависимый поиск книг по названию с поддержкой пагинации
// @Tags books
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Param limit query integer false "Количество результатов на странице (по умолчанию 50, максимум 100)" default(50) minimum(1) maximum(100)
// @Param offset query integer false "Смещение для пагинации (по умолчанию 0)" default(0) minimum(0)
// @Success 200 {object} map[string]interface{} "Результаты поиска"
// @Success 200 {object} map[string]interface{} "Успешный ответ с книгами"
// @Failure 400 {object} map[string]string "Отсутствует поисковый запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /book/search [get]
// SearchByTitle ищет книги по названию (регистронезависимый поиск)
func (h *Handler) SearchByTitle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Получаем поисковый запрос из query параметра
	searchQuery := strings.TrimSpace(c.Query("q"))
	if searchQuery == "" {
		c.JSON(400, gin.H{"error": "Search query parameter 'q' is required"})
		return
	}

	// Опциональные параметры пагинации
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50 // дефолтное значение
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Вызываем сервис для поиска книг
	books, err := h.srv.Book().SearchByTitle(ctx, searchQuery, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"books": books,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
			"query":  searchQuery,
		},
	}

	if len(books) == 0 {
		response["message"] = "No books found matching your search"
	}

	c.JSON(200, response)
}
