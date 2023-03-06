package userdomain

import (
	"time"
)

type User struct {
	ID string `bson:"_id"`

	Email    string `bson:"email"`
	Photo    string `bson:"photo"`
	Nickname string `bson:"nickname"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

func (p *User) DTO() *UserDTO {
	return &UserDTO{
		ID:       p.ID,
		Photo:    p.Photo,
		Nickname: p.Nickname,
	}
}
