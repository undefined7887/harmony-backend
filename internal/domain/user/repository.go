package userdomain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (bool, error)

	Get(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByNickname(ctx context.Context, nickname string) (*User, error)

	Exists(ctx context.Context, id string) (bool, error)
}
