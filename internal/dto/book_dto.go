package dto

type CreateBookRequest struct {
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description,omitempty"`
	Author      string  `db:"author" json:"author"`
	UploadedBy  int     `db:"uploaded_by" json:"uploaded_by"`
	FileUrl     string  `db:"file_url"`
}
