package callservice

import (
	"context"
	"encoding/json"
	"github.com/undefined7887/harmony-backend/internal/domain"
	calldomain "github.com/undefined7887/harmony-backend/internal/domain/call"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	"github.com/undefined7887/harmony-backend/internal/repository"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	userRepository userdomain.Repository
	callRepository calldomain.Repository

	centrifugoClient *centrifugo.Client
}

func NewService(
	userRepository userdomain.Repository,
	callRepository calldomain.Repository,
	centrifugoClient *centrifugo.Client,
) *Service {
	return &Service{
		userRepository:   userRepository,
		callRepository:   callRepository,
		centrifugoClient: centrifugoClient,
	}
}

func (s *Service) CreateCall(ctx context.Context, userID, peerID string) (string, error) {
	if err := s.checkPeer(ctx, userID); err != nil {
		return "", err
	}

	now := time.Now()

	call := calldomain.Call{
		ID:        domain.ID(),
		UserID:    userID,
		PeerID:    peerID,
		Status:    calldomain.StatusRequest,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.callRepository.Create(ctx, &call)
	if err != nil {
		return "", err
	}

	if !created {
		return "", calldomain.ErrCallAlreadyExists()
	}

	s.centrifugoPublish(
		ctx,
		calldomain.ChannelCallNew(peerID),
		calldomain.NewCallNotification{
			CallDTO: calldomain.MapCallDTO(call),
		},
	)

	return call.ID, nil
}

func (s *Service) GetCall(ctx context.Context, userID string) (calldomain.CallDTO, error) {
	call, err := s.callRepository.ReadLast(ctx, userID, calldomain.StatusRequest)
	if repository.IsNoDocumentsErr(err) {
		return calldomain.CallDTO{}, calldomain.ErrCallNotFound()
	}

	if err != nil {
		return calldomain.CallDTO{}, err
	}

	return calldomain.MapCallDTO(call), nil
}

func (s *Service) UpdateCallStatus(ctx context.Context, userID, id string, status string) error {
	var (
		matchPeerID           string
		matchPreviousStatuses []string
	)

	switch status {
	case calldomain.StatusAccepted:
		// Call can be accepted only peer
		matchPeerID = userID
		matchPreviousStatuses = []string{calldomain.StatusRequest}

	case calldomain.StatusDeclined:
		// Call can be declined by anyone
		matchPreviousStatuses = []string{calldomain.StatusRequest}

	case calldomain.StatusFinished:
		// Call can be finished by anyone
		matchPreviousStatuses = []string{calldomain.StatusRequest, calldomain.StatusAccepted}
	}

	call, err := s.callRepository.UpdateStatus(
		ctx,
		id,
		matchPeerID,
		matchPreviousStatuses,
		status,
	)
	if repository.IsNoDocumentsErr(err) {
		return calldomain.ErrCallNotFound()
	}

	if err != nil {
		return err
	}

	peerID := call.PeerID

	// Changing peer if peer is a current user
	if userID == call.PeerID {
		peerID = call.UserID
	}

	s.centrifugoPublish(
		ctx,
		calldomain.ChannelCallUpdates(peerID),
		calldomain.UpdateCallNotification{
			CallDTO: calldomain.MapCallDTO(call),
		},
	)

	return nil
}

func (s *Service) ProxyCallData(ctx context.Context, userID, id, name string, data json.RawMessage) error {
	call, err := s.callRepository.Read(ctx, id, calldomain.StatusAccepted)
	if repository.IsNoDocumentsErr(err) {
		return calldomain.ErrCallNotFound()
	}

	if err != nil {
		return err
	}

	peerID := call.PeerID

	// Changing peer if peer is a current user
	if userID == call.PeerID {
		peerID = call.UserID
	}

	s.centrifugoPublish(ctx,
		calldomain.ChannelCallData(peerID),
		calldomain.CallDataNotification{
			ID:   call.ID,
			Name: name,
			Data: data,
		},
	)

	return nil
}

func (s *Service) checkPeer(ctx context.Context, peerID string) error {
	exists, err := s.userRepository.Exists(ctx, peerID)
	if err != nil {
		return err
	}

	if !exists {
		return userdomain.ErrUserNotFound()
	}

	return nil
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
