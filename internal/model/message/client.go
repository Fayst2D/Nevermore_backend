package chat

type Client struct {
	UserID   int
	Username string
	Send     chan Message
}
