package dto

type SaveAuthorRequest struct {
	UserId   int `json:"user_id"`
	AuthorId int `json:"author_id"`
}
