package dto

import "time"

type ChatMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type ChatMessageResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type OnlineUsersResponse struct {
	Count int      `json:"count"`
	Users []string `json:"users"`
}
