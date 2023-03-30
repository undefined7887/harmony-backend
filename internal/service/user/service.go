package userservice

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	"github.com/undefined7887/harmony-backend/internal/repository"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
)

const (
	backgroundUpdateStatusesTimeout = time.Second * 30
)

type Service struct {
	logger           *zap.Logger
	userRepository   userdomain.Repository
	centrifugoClient *centrifugo.Client
}

func NewService(
	logger *zap.Logger,
	userRepository userdomain.Repository,
	centrifugoClient *centrifugo.Client,
) *Service {
	return &Service{
		logger:           logger,
		userRepository:   userRepository,
		centrifugoClient: centrifugoClient,
	}
}

func NewServiceRunner(lifecycle fx.Lifecycle, logger *zap.Logger, service *Service) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting background statuses updater")

			go service.BackgroundUpdateStatuses(context.Background())

			return nil
		},
		OnStop: nil,
	})
}

func (s *Service) GetUser(ctx context.Context, id string) (userdomain.UserDTO, error) {
	user, err := s.userRepository.Get(ctx, id)
	if repository.IsNoDocumentsErr(err) {
		return userdomain.UserDTO{}, userdomain.ErrUserNotFound()
	}

	if err != nil {
		return userdomain.UserDTO{}, err
	}

	return userdomain.MapUserDTO(user), nil
}

func (s *Service) SearchUser(ctx context.Context, nickname string) (userdomain.UserDTO, error) {
	user, err := s.userRepository.GetByNickname(ctx, nickname)
	if repository.IsNoDocumentsErr(err) {
		return userdomain.UserDTO{}, userdomain.ErrUserNotFound()
	}

	if err != nil {
		return userdomain.UserDTO{}, err
	}

	return userdomain.MapUserDTO(user), nil
}

func (s *Service) UpdateUserStatus(ctx context.Context, userID, status string, onlyOffline bool) error {
	user, err := s.userRepository.UpdateStatus(ctx, userID, status, onlyOffline)
	if repository.IsNoDocumentsErr(err) {
		return userdomain.ErrUserNotFound()
	}

	if err != nil {
		return err
	}

	s.centrifugoPublish(
		ctx,
		userdomain.ChannelUser(userID),
		userdomain.UpdateUserNotification{
			UserDTO: userdomain.MapUserDTO(user),
		},
	)

	return nil
}

func (s *Service) BackgroundUpdateStatuses(ctx context.Context) {
	ticker := time.NewTicker(backgroundUpdateStatusesTimeout)

	for {
		<-ticker.C

		ctx, cancel := context.WithTimeout(ctx, backgroundUpdateStatusesTimeout)

		var count int64

		err := s.userRepository.UpdateOutdatedStatuses(ctx, func(users []userdomain.User) {
			// We need to rewrite this with command pipelines
			// For more information: https://centrifugal.dev/docs/server/server_api#command-pipelining
			for _, user := range users {
				s.centrifugoPublish(
					ctx,
					userdomain.ChannelUser(user.ID),
					userdomain.UpdateUserNotification{
						UserDTO: userdomain.MapUserDTO(user),
					},
				)
			}

			count += int64(len(users))
		})

		if err != nil {
			s.logger.Warn("background update statuses error", zap.Error(err))
		} else {
			s.logger.Info("background update statuses successful", zap.Int64("updated_count", count))
		}

		cancel()
	}
}

func (s *Service) centrifugoPublish(ctx context.Context, channel string, data any) {
	if _, err := s.centrifugoClient.Publish(ctx, channel, data); err != nil {
		zaplog.
			UnpackLogger(ctx).
			Warn(
				"centrifugo publish error",
				zap.String("channel", channel),
				zap.Error(err),
			)
	}
}
