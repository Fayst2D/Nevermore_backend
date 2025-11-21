package chat

import (
	"context"
	"fmt"
	"nevermore/internal/model/message"
	"nevermore/internal/storage"
	"sync"

	"github.com/gammazero/workerpool"
)

type Service interface {
	BroadcastMessage(message chat.Message) error
	GetRecentMessages(ctx context.Context, limit int) ([]chat.Message, error)
	AddClient(client *chat.Client)
	RemoveClient(client *chat.Client)
	GetOnlineUsers() []string

	//Личные сообщения
	SendPrivateMessage(ctx context.Context, message chat.PrivateMessage) error
	GetPrivateMessages(ctx context.Context, user1ID, user2ID int, limit int) ([]chat.PrivateMessage, error)
	GetConversations(ctx context.Context, userID int) ([]chat.PrivateConversation, error)
	MarkMessagesAsRead(ctx context.Context, messageIDs []int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

type service struct {
	st      storage.Storage
	wp      *workerpool.WorkerPool
	clients sync.Map
	mu      sync.RWMutex
}

func New(st storage.Storage, wp *workerpool.WorkerPool) Service {
	return &service{
		st: st,
		wp: wp,
	}
}

func (s *service) BroadcastMessage(message chat.Message) error {
	// Сохраняем сообщение в БД
	ctx := context.Background()
	err := s.st.DB().Chat().CreateMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Рассылаем сообщение всем подключенным клиентам
	s.clients.Range(func(key, value interface{}) bool {
		client := value.(*chat.Client)
		select {
		case client.Send <- message:
		default:
			// Если канал полный, пропускаем клиента
			fmt.Printf("Client %s channel is full, skipping\n", client.Username)
		}
		return true
	})

	return nil
}

func (s *service) GetRecentMessages(ctx context.Context, limit int) ([]chat.Message, error) {
	return s.st.DB().Chat().GetRecentMessages(ctx, limit)
}

func (s *service) AddClient(client *chat.Client) {
	s.clients.Store(client.UserID, client)

	// Отправляем системное сообщение о подключении
	systemMessage := chat.Message{
		UserID:   0,
		Username: "System",
		Content:  fmt.Sprintf("%s присоединился к чату", client.Username),
		Type:     "join",
	}

	s.wp.Submit(func() {
		s.BroadcastMessage(systemMessage)
	})
}

func (s *service) RemoveClient(client *chat.Client) {
	s.clients.Delete(client.UserID)

	// Отправляем системное сообщение об отключении
	systemMessage := chat.Message{
		UserID:   0,
		Username: "System",
		Content:  fmt.Sprintf("%s покинул чат", client.Username),
		Type:     "leave",
	}

	s.wp.Submit(func() {
		s.BroadcastMessage(systemMessage)
	})
}

func (s *service) GetOnlineUsers() []string {
	var users []string
	s.clients.Range(func(key, value interface{}) bool {
		client := value.(*chat.Client)
		users = append(users, client.Username)
		return true
	})
	return users
}

//Личные сообщения

func (s *service) SendPrivateMessage(ctx context.Context, message chat.PrivateMessage) error {
	// Сохраняем сообщение в БД
	err := s.st.DB().PrivateMessage().CreatePrivateMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to save private message: %w", err)
	}

	// Отправляем сообщение получателю, если он онлайн
	if receiver, ok := s.clients.Load(message.ReceiverID); ok {
		receiverClient := receiver.(*chat.Client)
		select {
		case receiverClient.Send <- chat.Message{
			ID:        message.ID,
			UserID:    message.SenderID,
			Username:  message.SenderName,
			Content:   fmt.Sprintf("[Private] %s", message.Content),
			Type:      "private",
			CreatedAt: message.CreatedAt,
		}:
		default:
			fmt.Printf("Receiver %s channel is full, skipping\n", message.ReceiverName)
		}
	}

	return nil
}

func (s *service) GetPrivateMessages(ctx context.Context, user1ID, user2ID int, limit int) ([]chat.PrivateMessage, error) {
	return s.st.DB().PrivateMessage().GetPrivateMessages(ctx, user1ID, user2ID, limit)
}

func (s *service) GetConversations(ctx context.Context, userID int) ([]chat.PrivateConversation, error) {
	return s.st.DB().PrivateMessage().GetConversations(ctx, userID)
}

func (s *service) MarkMessagesAsRead(ctx context.Context, messageIDs []int) error {
	return s.st.DB().PrivateMessage().MarkMessagesAsRead(ctx, messageIDs)
}

func (s *service) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	return s.st.DB().PrivateMessage().GetUnreadCount(ctx, userID)
}
