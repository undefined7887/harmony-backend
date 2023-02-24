package userdomain

import "time"

type User struct {
	ID string `bson:"id"`

	Email    string `bson:"email"`
	Photo    string `bson:"photo"`
	Nickname string `bson:"nickname"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
