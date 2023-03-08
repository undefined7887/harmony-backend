package userservice

import (
	"context"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
)

type Service struct {
	userRepository userdomain.Repository
}

func NewService(
	userRepository userdomain.Repository,
) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) GetUser(ctx context.Context, id string) (*userdomain.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, userdomain.ErrUserNotFound()
	}

	return user, nil
}

func (s *Service) GetUserByNickname(ctx context.Context, nickname string) (*userdomain.User, error) {
	user, err := s.userRepository.GetByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, userdomain.ErrUserNotFound()
	}

	return user, nil
}
