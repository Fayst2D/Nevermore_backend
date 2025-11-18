package chat

import "time"

type Message struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Content   string    `db:"content" json:"content"`
	Type      string    `db:"type" json:"type"` // "message", "system", "join", "leave"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
