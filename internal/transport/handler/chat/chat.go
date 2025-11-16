package chat

import (
	"context"
	"nevermore/internal/dto"
	"nevermore/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// @Summary Create message
// @Description Create new message in chat room
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMessageRequest true "Message data"
// @Success 200 {object} string "Message created successfully"
// @Failure 400 {object} string "Bad request - invalid data"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/message [post]
func (h *Handler) CreateMessage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Получаем ID пользователя из контекста (после аутентификации)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Получаем данные из JSON body
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Валидация обязательных полей
	if req.Content == "" {
		c.JSON(400, gin.H{"error": "Content is required"})
		return
	}
	if req.RoomID == "" {
		c.JSON(400, gin.H{"error": "RoomID is required"})
		return
	}

	// Устанавливаем пользователя из контекста
	req.UserID = userID.(string)

	// Если username не передан, используем userID или можно получить из БД
	if req.Username == "" {
		req.Username = "User_" + userID.(string)
	}

	// Устанавливаем тип сообщения по умолчанию
	if req.Type == "" {
		req.Type = "message"
	}

	err := h.srv.Chat().CreateMessage(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message created successfully"})
}

// @Summary Get chat history
// @Description Get message history for specific chat room
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id query string true "Room ID"
// @Param limit query int false "Limit number of messages (default: 50, max: 1000)" default(50) minimum(1) maximum(1000)
//
//	@Success 200 {object} object "Chat history" {object} struct {
//	    Messages []dto.Message "List of messages",
//	    RoomID   string       "Room ID",
//	    Count    int          "Number of messages"
//	}
//
// @Failure 400 {object} string "Bad request - missing room_id"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/history [get]
func (h *Handler) GetChatHistory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(400, gin.H{"error": "RoomID is required"})
		return
	}

	// Получаем лимит из query параметров (по умолчанию 50)
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	// Ограничиваем максимальный лимит
	if limit > 1000 {
		limit = 1000
	}

	messages, err := h.srv.Chat().GetChatHistory(ctx, roomID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"messages": messages,
		"room_id":  roomID,
		"count":    len(messages),
	})
}

// @Summary WebSocket connection
// @Description Establish WebSocket connection for real-time chat
// @Tags chat
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param username query string true "Username"
// @Param room_id query string true "Room ID"
//
//	@Success 200 {object} object "WebSocket connection established" {object} struct {
//	    Message string "Connection message",
//	    UserID  string "User ID",
//	    RoomID  string "Room ID"
//	}
//
// @Failure 400 {object} string "Bad request - missing required parameters"
// @Router /chat/ws [get]
func (h *Handler) WebSocketHandler(c *gin.Context) {
	// Получаем параметры из query string
	userID := c.Query("user_id")
	username := c.Query("username")
	roomID := c.Query("room_id")

	if userID == "" || username == "" || roomID == "" {
		c.JSON(400, gin.H{"error": "user_id, username and room_id are required"})
		return
	}

	// Здесь будет логика upgrade до WebSocket соединения
	// (используется отдельный WebSocket handler как в предыдущем примере)

	c.JSON(200, gin.H{
		"message": "WebSocket connection established",
		"user_id": userID,
		"room_id": roomID,
	})
}
