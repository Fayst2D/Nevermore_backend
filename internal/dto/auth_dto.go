package dto

type RegisterRequest struct {
	Name        string `db:"name" json:"name"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
	Email       string `db:"email" json:"email"`
	Password    string `db:"password" json:"password"`
}

type LoginRequest struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}
