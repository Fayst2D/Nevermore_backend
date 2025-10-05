package author

import "time"

type Author struct {
	Id        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Biography string     `db:"biography" json:"biography"`
	Photo     *string    `db:"photo_url" json:"photo_url"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
