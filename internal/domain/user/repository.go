package userdomain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (bool, error)

	Get(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByNickname(ctx context.Context, nickname string) (User, error)

	Exists(ctx context.Context, id string) (bool, error)

	UpdatePhoto(ctx context.Context, id, photo string) (User, error)
	UpdateStatus(ctx context.Context, id, status string, onlyOffline bool) (User, error)

	// UpdateOutdatedStatuses sets status 'offline' to users who weren't updated for 1 minute
	UpdateOutdatedStatuses(ctx context.Context, cb func(users []User)) error
}
