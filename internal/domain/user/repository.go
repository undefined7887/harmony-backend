package userdomain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (bool, error)

	Read(ctx context.Context, id string) (*User, error)
	ReadByEmail(ctx context.Context, email string) (*User, error)

	Exists(ctx context.Context, id string) (bool, error)
}
