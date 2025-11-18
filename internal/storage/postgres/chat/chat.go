package chat

import (
	"context"
	"nevermore/internal/model/message"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	CreateMessage(ctx context.Context, message chat.Message) error
	GetRecentMessages(ctx context.Context, limit int) ([]chat.Message, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repo {
	return &repo{db: db}
}

func (s *repo) CreateMessage(ctx context.Context, message chat.Message) error {
	query := `
		INSERT INTO chat_messages (user_id, username, content, type, created_at) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.ExecContext(ctx, query,
		message.UserID,
		message.Username,
		message.Content,
		message.Type,
		time.Now(),
	)

	return err
}

func (s *repo) GetRecentMessages(ctx context.Context, limit int) ([]chat.Message, error) {
	query := `
		SELECT id, user_id, username, content, type, created_at 
		FROM chat_messages 
		ORDER BY created_at DESC 
		LIMIT $1
	`

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []chat.Message
	for rows.Next() {
		var msg chat.Message
		err := rows.Scan(
			&msg.ID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.Type,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
