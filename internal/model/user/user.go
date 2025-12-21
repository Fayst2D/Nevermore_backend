package user

import (
	"time"
)

type User struct {
	Id                int        `db:"id" json:"id"`
	Name              string     `db:"name" json:"name"`
	IsEmailVerified   bool       `db:"email_verified" json:"email_verified"`
	VerificationToken *string    `db:"verification_token" json:"-"`
	PhoneNumber       string     `db:"phone_number" json:"phone_number"`
	Email             string     `db:"email" json:"email"`
	Role              string     `db:"role" json:"role"`
	Password          string     `db:"password" json:"password"`
	Photo             *string    `db:"photo" json:"photo"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deleted_at"`
}
