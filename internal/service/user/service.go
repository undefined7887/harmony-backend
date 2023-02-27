package userservice

import (
	"context"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
)

type Service struct {
	repository userdomain.Repository
}

func NewService(
	repository userdomain.Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*userdomain.User, error) {
	user, err := s.repository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, userdomain.ErrUserNotFound()
	}

	return user, nil
}
