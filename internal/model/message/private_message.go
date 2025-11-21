package chat

import "time"

type PrivateMessage struct {
	ID           int       `db:"id" json:"id"`
	SenderID     int       `db:"sender_id" json:"sender_id"`
	ReceiverID   int       `db:"receiver_id" json:"receiver_id"`
	SenderName   string    `db:"sender_name" json:"sender_name"`
	ReceiverName string    `db:"receiver_name" json:"receiver_name"`
	Content      string    `db:"content" json:"content"`
	IsRead       bool      `db:"is_read" json:"is_read"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type PrivateConversation struct {
	UserID        int       `json:"user_id"`
	Username      string    `json:"username"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`
}
