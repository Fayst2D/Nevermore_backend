package dto

import (
	"time"
)

type PrivateMessageRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

type PrivateMessageResponse struct {
	ID           int       `json:"id"`
	SenderID     int       `json:"sender_id"`
	ReceiverID   int       `json:"receiver_id"`
	SenderName   string    `json:"sender_name"`
	ReceiverName string    `json:"receiver_name"`
	Content      string    `json:"content"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    time.Time `json:"created_at"`
}

type MarkAsReadRequest struct {
	MessageIDs []int `json:"message_ids" binding:"required"`
}
