package dto

type AuthorGetResponse struct {
	Name      string  `db:"name" json:"name"`
	Biography string  `db:"biography" json:"biography"`
	Photo     *string `db:"photo_url" json:"photo_url"`
}
