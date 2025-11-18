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
