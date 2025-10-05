package dto

type UserGetResponse struct {
	Name        string  `db:"name" json:"name"`
	PhoneNumber string  `db:"phone_number" json:"phone_number"`
	Email       string  `db:"email" json:"email"`
	Role        string  `db:"role" json:"role"`
	Photo       *string `db:"photo" json:"photo"`
}

type UpdateUserRequest struct {
	Name        string  `db:"name" json:"name"`
	PhoneNumber string  `db:"phone_number" json:"phone_number"`
	Password    string  `db:"password" json:"password"`
	Photo       *string `db:"photo" json:"photo"`
}
