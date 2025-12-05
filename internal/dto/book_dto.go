package dto

type CreateBookRequest struct {
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description,omitempty"`
	Author      string  `db:"author" json:"author"`
	UploadedBy  int     `db:"uploaded_by" json:"uploaded_by"`
	FileUrl     string  `db:"file_url"`
}

type GetBookRequest struct {
	Title           string  `db:"title" json:"title"`
	Description     *string `db:"description" json:"description,omitempty"`
	Cover_image_url *string `db:"cover_image_url" json:"cover_image_url,omitempty"`
	Author_id       string  `db:"author_id" json:"author"`
	UploadedBy      int     `db:"uploaded_by" json:"uploaded_by"`
	FileUrl         string  `db:"file_url"`
}
