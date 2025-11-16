package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"nevermore/internal/dto"
	"nevermore/internal/model/message"
	"nevermore/internal/storage"

	"github.com/gammazero/workerpool"
)

type Service interface {
	CreateMessage(ctx context.Context, req *dto.CreateMessageRequest) error
	GetChatHistory(ctx context.Context, roomID string, limit int) ([]*model.Message, error)
	GetStats(ctx context.Context) (*dto.ChatStats, error)
	RegisterConnection(conn *dto.Connection)
	UnregisterConnection(conn *dto.Connection)
	BroadcastMessage(message *model.Message) error
	Run()
}

type service struct {
	st          storage.Storage
	connections map[string]map[string]*dto.Connection // roomID -> userID -> Connection
	broadcast   chan []byte
	register    chan *dto.Connection
	unregister  chan *dto.Connection
	mu          sync.RWMutex
	workerPool  *workerpool.WorkerPool
}

func New(st storage.Storage, wp *workerpool.WorkerPool) Service {
	srv := &service{
		st:          st,
		connections: make(map[string]map[string]*dto.Connection),
		broadcast:   make(chan []byte, 100),
		register:    make(chan *dto.Connection, 10),
		unregister:  make(chan *dto.Connection, 10),
		workerPool:  wp,
	}

	// Запускаем обработчик соединений в отдельной горутине
	go srv.Run()
	return srv
}

func (s *service) Run() {
	for {
		select {
		case connection := <-s.register:
			s.handleRegister(connection)

		case connection := <-s.unregister:
			s.handleUnregister(connection)

		case message := <-s.broadcast:
			s.handleBroadcast(message)
		}
	}
}

func (s *service) handleRegister(connection *dto.Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.connections[connection.RoomID]; !exists {
		s.connections[connection.RoomID] = make(map[string]*dto.Connection)
	}

	s.connections[connection.RoomID][connection.UserID] = connection

	// Асинхронно отправляем системное сообщение о новом пользователе
	s.workerPool.Submit(func() {
		s.broadcastSystemMessage(
			connection.RoomID,
			fmt.Sprintf("Пользователь %s присоединился к чату", connection.Username),
		)
	})

	log.Printf("Пользователь %s подключен к комнате %s. Онлайн: %d",
		connection.Username, connection.RoomID, len(s.connections[connection.RoomID]))
}

func (s *service) handleUnregister(connection *dto.Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if roomConnections, exists := s.connections[connection.RoomID]; exists {
		if _, userExists := roomConnections[connection.UserID]; userExists {
			delete(roomConnections, connection.UserID)
			close(connection.Send)

			// Удаляем комнату если она пустая
			if len(roomConnections) == 0 {
				delete(s.connections, connection.RoomID)
			}
		}
	}

	// Асинхронно отправляем системное сообщение о выходе пользователя
	s.workerPool.Submit(func() {
		s.broadcastSystemMessage(
			connection.RoomID,
			fmt.Sprintf("Пользователь %s покинул чат", connection.Username),
		)
	})

	log.Printf("Пользователь %s отключен от комнаты %s", connection.Username, connection.RoomID)
}

func (s *service) handleBroadcast(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// В этой базовой реализации broadcast отправляется во все комнаты
	// Можно доработать для отправки в конкретную комнату
	for roomID, roomConnections := range s.connections {
		for userID, connection := range roomConnections {
			select {
			case connection.Send <- message:
				// Сообщение отправлено
			default:
				// Если канал заполнен, закрываем соединение
				close(connection.Send)
				delete(roomConnections, userID)

				if len(roomConnections) == 0 {
					delete(s.connections, roomID)
				}
			}
		}
	}
}

func (s *service) RegisterConnection(conn *dto.Connection) {
	s.register <- conn
}

func (s *service) UnregisterConnection(conn *dto.Connection) {
	s.unregister <- conn
}

func (s *service) CreateMessage(ctx context.Context, req *dto.CreateMessageRequest) error {
	tx, err := s.st.DB().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("ChatService:CreateMessage err -> %s", err.Error())
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Rollback error: %v", err)
		}
	}()

	// Сохраняем сообщение в БД
	message := &model.Message{
		ID:        generateMessageID(),
		UserID:    req.UserID,
		Username:  req.Username,
		Content:   req.Content,
		Type:      req.Type,
		Timestamp: time.Now(),
		RoomID:    req.RoomID,
	}

	if err := s.st.DB().Chat().CreateMessage(ctx, tx, message); err != nil {
		return fmt.Errorf("ChatService:CreateMessage err -> %s", err.Error())
	}

	// Асинхронно broadcast сообщение
	s.workerPool.Submit(func() {
		if err := s.BroadcastMessage(message); err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	})

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ChatService:CreateMessage err -> %s", err.Error())
	}

	return nil
}

func (s *service) BroadcastMessage(message *model.Message) error {
	response := model.ChatResponse{
		Message:     message,
		OnlineUsers: s.getOnlineUsersCountByRoom(message.RoomID),
	}

	wsMessage := dto.WSMessage{
		Type:    "chat_message",
		Payload: response,
	}

	messageBytes, err := json.Marshal(wsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	s.broadcast <- messageBytes
	return nil
}

func (s *service) broadcastSystemMessage(roomID, content string) {
	message := &model.Message{
		ID:        generateMessageID(),
		UserID:    "system",
		Username:  "System",
		Content:   content,
		Type:      "system",
		Timestamp: time.Now(),
		RoomID:    roomID,
	}

	if err := s.BroadcastMessage(message); err != nil {
		log.Printf("Error broadcasting system message: %v", err)
	}
}

func (s *service) GetChatHistory(ctx context.Context, roomID string, limit int) ([]*model.Message, error) {
	messages, err := s.st.DB().Chat().GetMessagesByRoom(ctx, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("ChatService:GetChatHistory err -> %s", err.Error())
	}
	return messages, nil
}

func (s *service) GetStats(ctx context.Context) (*dto.ChatStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &dto.ChatStats{
		OnlineUsers:      make([]string, 0),
		ActiveRooms:      make([]string, 0),
		OnlineUsersCount: 0,
	}

	for roomID, roomConnections := range s.connections {
		stats.ActiveRooms = append(stats.ActiveRooms, roomID)
		stats.OnlineUsersCount += len(roomConnections)

		for _, connection := range roomConnections {
			stats.OnlineUsers = append(stats.OnlineUsers, connection.Username)
		}
	}

	return stats, nil
}

func (s *service) getOnlineUsersCountByRoom(roomID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if roomConnections, exists := s.connections[roomID]; exists {
		return len(roomConnections)
	}
	return 0
}

func generateMessageID() string {
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}
