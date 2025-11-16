package dto

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Connection struct {
	UserID   string
	Username string
	RoomID   string
	Send     chan []byte
}

type CreateMessageRequest struct {
	UserID   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Content  string `json:"content" db:"content"`
	Type     string `json:"type" db:"type"`
	RoomID   string `json:"room_id" db:"room_id"`
}

type ChatStats struct {
	OnlineUsersCount int      `json:"online_users_count"`
	OnlineUsers      []string `json:"online_users"`
	ActiveRooms      []string `json:"active_rooms"`
}
