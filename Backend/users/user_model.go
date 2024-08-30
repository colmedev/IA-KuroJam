package users

import "time"

type User struct {
	Id        int64     `db:"id" json:"id,omitempty"`
	ClerkId   string    `db:"clerk_id" json:"-"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Version   int       `db:"version" json:"-"`
}
