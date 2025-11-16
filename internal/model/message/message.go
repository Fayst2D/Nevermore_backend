package model

import (
	"time"
)

type Message struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"` // message, system, notification
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	RoomID    string    `json:"room_id" db:"room_id"`
}

type ChatRequest struct {
	Content string `json:"content" binding:"required"`
	RoomID  string `json:"room_id"`
}

type ChatResponse struct {
	Message     *Message `json:"message"`
	OnlineUsers int      `json:"online_users"`
}

type ConnectionInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
}
