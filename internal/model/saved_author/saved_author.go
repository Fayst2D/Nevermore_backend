package saved_author

import "time"

type SavedAuthor struct {
	UserId    int       `json:"user_id"`
	AuthorId  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
