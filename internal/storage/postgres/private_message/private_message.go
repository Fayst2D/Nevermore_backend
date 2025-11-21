package private_message

import (
	"context"
	"fmt"
	"nevermore/internal/model/message"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	CreatePrivateMessage(ctx context.Context, message chat.PrivateMessage) error
	GetPrivateMessages(ctx context.Context, user1ID, user2ID int, limit int) ([]chat.PrivateMessage, error)
	GetConversations(ctx context.Context, userID int) ([]chat.PrivateConversation, error)
	MarkMessagesAsRead(ctx context.Context, messageIDs []int) error
	GetUnreadCount(ctx context.Context, userID int) (int, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repo {
	return &repo{db: db}
}

func (s *repo) CreatePrivateMessage(ctx context.Context, message chat.PrivateMessage) error {
	query := `
		INSERT INTO private_messages (sender_id, receiver_id, sender_name, receiver_name, content, is_read, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := s.db.QueryRowContext(ctx, query,
		message.SenderID,
		message.ReceiverID,
		message.SenderName,
		message.ReceiverName,
		message.Content,
		message.IsRead,
		time.Now(),
	).Scan(&message.ID)

	return err
}

func (s *repo) GetPrivateMessages(ctx context.Context, user1ID, user2ID int, limit int) ([]chat.PrivateMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, sender_name, receiver_name, content, is_read, created_at
		FROM private_messages 
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC 
		LIMIT $3
	`

	rows, err := s.db.QueryContext(ctx, query, user1ID, user2ID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []chat.PrivateMessage
	for rows.Next() {
		var msg chat.PrivateMessage
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.SenderName,
			&msg.ReceiverName,
			&msg.Content,
			&msg.IsRead,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Реверсируем порядок для хронологического отображения
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (s *repo) GetConversations(ctx context.Context, userID int) ([]chat.PrivateConversation, error) {
	query := `
		WITH LastMessages AS (
			SELECT 
				CASE 
					WHEN sender_id = $1 THEN receiver_id 
					ELSE sender_id 
				END as other_user_id,
				CASE 
					WHEN sender_id = $1 THEN receiver_name 
					ELSE sender_name 
				END as other_username,
				content as last_message,
				created_at as last_message_at,
				is_read,
				id as message_id
			FROM private_messages 
			WHERE sender_id = $1 OR receiver_id = $1
		),
		RankedMessages AS (
			SELECT *,
				ROW_NUMBER() OVER (PARTITION BY other_user_id ORDER BY last_message_at DESC) as rn
			FROM LastMessages
		),
		UnreadCounts AS (
			SELECT 
				CASE 
					WHEN sender_id = $1 THEN receiver_id 
					ELSE sender_id 
				END as other_user_id,
				COUNT(*) as unread_count
			FROM private_messages 
			WHERE receiver_id = $1 AND is_read = false
			GROUP BY other_user_id
		)
		SELECT 
			rm.other_user_id,
			rm.other_username,
			rm.last_message,
			rm.last_message_at,
			COALESCE(uc.unread_count, 0) as unread_count
		FROM RankedMessages rm
		LEFT JOIN UnreadCounts uc ON rm.other_user_id = uc.other_user_id
		WHERE rm.rn = 1
		ORDER BY rm.last_message_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []chat.PrivateConversation
	for rows.Next() {
		var conv chat.PrivateConversation
		err := rows.Scan(
			&conv.UserID,
			&conv.Username,
			&conv.LastMessage,
			&conv.LastMessageAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (s *repo) MarkMessagesAsRead(ctx context.Context, messageIDs []int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	// Создаем плейсхолдеры для IN запроса
	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE private_messages 
		SET is_read = true 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *repo) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM private_messages 
		WHERE receiver_id = $1 AND is_read = false
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
