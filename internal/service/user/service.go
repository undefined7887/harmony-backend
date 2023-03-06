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

func (s *Service) Read(ctx context.Context, id string) (*userdomain.User, error) {
	user, err := s.userRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, userdomain.ErrUserNotFound()
	}

	return user, nil
}

func (s *Service) ReadByNickname(ctx context.Context, nickname string) (*userdomain.User, error) {
	user, err := s.userRepository.ReadByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, userdomain.ErrUserNotFound()
	}

	return user, nil
}
