package chat

import (
	"context"
	"nevermore/internal/model/message"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	CreateMessage(ctx context.Context, tx *sqlx.Tx, message *model.Message) error
	GetMessagesByRoom(ctx context.Context, roomID string, limit int) ([]*model.Message, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repo {
	return &repo{db: db}
}

func (s *repo) CreateMessage(ctx context.Context, tx *sqlx.Tx, message *model.Message) error {
	query := `
		INSERT INTO messages (id, user_id, username, content, type, timestamp, room_id)
		VALUES (:id, :user_id, :username, :content, :type, :timestamp, :room_id)
	`

	_, err := tx.NamedExecContext(ctx, query, message)
	return err
}

func (s *repo) GetMessagesByRoom(ctx context.Context, roomID string, limit int) ([]*model.Message, error) {
	var messages []*model.Message

	query := `
		SELECT id, user_id, username, content, type, timestamp, room_id
		FROM messages 
		WHERE room_id = $1 
		ORDER BY timestamp DESC 
		LIMIT $2
	`

	err := s.db.SelectContext(ctx, &messages, query, roomID, limit)
	if err != nil {
		return nil, err
	}

	// Реверсируем порядок для хронологического отображения
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
