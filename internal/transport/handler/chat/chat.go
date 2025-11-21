package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nevermore/internal/dto"
	"nevermore/internal/model/message"
	"nevermore/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	srv      service.Service
	upgrader websocket.Upgrader
}

func New(srv service.Service) *Handler {
	return &Handler{
		srv: srv,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // В продакшене нужно настроить properly
			},
		},
	}
}

// @Summary WebSocket соединение для чата
// @Description Установка WebSocket соединения для участия в чате
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Router /chat/ws [get]
func (h *Handler) WebSocketHandler(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Получаем информацию о пользователе из контекста (должна быть установлена middleware аутентификации)
	userID, exists := c.Get("userID")
	if !exists {
		conn.WriteMessage(websocket.CloseMessage, []byte("Unauthorized"))
		return
	}

	username, exists := c.Get("username")
	if !exists {
		conn.WriteMessage(websocket.CloseMessage, []byte("Username not found"))
		return
	}

	// Используем модель chat.Client вместо service.Client
	client := &chat.Client{
		UserID:   userID.(int),
		Username: username.(string),
		Send:     make(chan chat.Message, 256),
	}

	h.srv.Chat().AddClient(client)
	defer h.srv.Chat().RemoveClient(client)

	// Горутина для отправки сообщений клиенту
	go h.writePump(conn, client)

	// Горутина для чтения сообщений от клиента
	h.readPump(conn, client)
}

func (h *Handler) writePump(conn *websocket.Conn, client *chat.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				// Канал закрыт
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}

		case <-ticker.C:
			// Ping для поддержания соединения
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Handler) readPump(conn *websocket.Conn, client *chat.Client) {
	defer close(client.Send)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}

		var msgReq dto.ChatMessageRequest
		if err := json.Unmarshal(message, &msgReq); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Создаем сообщение для рассылки
		chatMessage := chat.Message{
			UserID:    client.UserID,
			Username:  client.Username,
			Content:   msgReq.Content,
			Type:      "message",
			CreatedAt: time.Now(),
		}

		// Отправляем сообщение всем клиентам
		if err := h.srv.Chat().BroadcastMessage(chatMessage); err != nil {
			log.Printf("Broadcast error: %v", err)
		}
	}
}

// @Summary Получить историю сообщений
// @Description Получить последние сообщения из чата
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param limit query int false "Количество сообщений" default(50)
// @Success 200 {array} dto.ChatMessageResponse "Список сообщений"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/messages [get]
func (h *Handler) GetMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	messages, err := h.srv.Chat().GetRecentMessages(c.Request.Context(), limit)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get messages: %v", err)})
		return
	}

	c.JSON(200, messages)
}

// @Summary Получить список онлайн пользователей
// @Description Получить список пользователей онлайн в чате
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.OnlineUsersResponse "Список онлайн пользователей"
// @Router /chat/online [get]
func (h *Handler) GetOnlineUsers(c *gin.Context) {
	users := h.srv.Chat().GetOnlineUsers()
	response := dto.OnlineUsersResponse{
		Count: len(users),
		Users: users,
	}

	c.JSON(200, response)
}

// Личные сообщения
// @Summary Отправить личное сообщение
// @Description Отправка личного сообщения другому пользователю
// @Tags private-chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PrivateMessageRequest true "Данные сообщения"
// @Success 200 {object} string "Сообщение отправлено"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/private/send [post]
func (h *Handler) SendPrivateMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req dto.PrivateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	// Создаем приватное сообщение
	privateMessage := chat.PrivateMessage{
		SenderID:     userID.(int),
		ReceiverID:   req.ReceiverID,
		SenderName:   username.(string),
		ReceiverName: "", // Можно получить из БД или кэша
		Content:      req.Content,
		IsRead:       false,
		CreatedAt:    time.Now(),
	}

	err := h.srv.Chat().SendPrivateMessage(c.Request.Context(), privateMessage)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to send message: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Private message sent successfully"})
}

// @Summary Получить личную переписку
// @Description Получить историю личных сообщений с пользователем
// @Tags private-chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Param limit query int false "Количество сообщений" default(50)
// @Success 200 {array} dto.PrivateMessageResponse "История сообщений"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/private/conversation/{user_id} [get]
func (h *Handler) GetPrivateConversation(c *gin.Context) {
	userID, _ := c.Get("userID")

	otherUserIDStr := c.Param("user_id")
	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	messages, err := h.srv.Chat().GetPrivateMessages(c.Request.Context(), userID.(int), otherUserID, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get conversation: %v", err)})
		return
	}

	c.JSON(200, messages)
}

// @Summary Пометить сообщения как прочитанные
// @Description Пометить личные сообщения как прочитанные
// @Tags private-chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.MarkAsReadRequest true "ID сообщений"
// @Success 200 {object} string "Сообщения помечены как прочитанные"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/private/mark-read [post]
func (h *Handler) MarkMessagesAsRead(c *gin.Context) {
	var req dto.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	err := h.srv.Chat().MarkMessagesAsRead(c.Request.Context(), req.MessageIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to mark messages as read: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Messages marked as read"})
}

// @Summary Получить количество непрочитанных сообщений
// @Description Получить общее количество непрочитанных личных сообщений
// @Tags private-chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int "Количество непрочитанных сообщений"
// @Failure 500 {object} string "Internal server error"
// @Router /chat/private/unread-count [get]
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userID")

	count, err := h.srv.Chat().GetUnreadCount(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get unread count: %v", err)})
		return
	}

	c.JSON(200, gin.H{"unread_count": count})
}
