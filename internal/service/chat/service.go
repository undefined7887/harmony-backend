package chatservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/undefined7887/harmony-backend/internal/domain"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"time"
)

type Service struct {
	userRepository userdomain.Repository
	chatRepository chatdomain.Repository

	centrifugoClient *centrifugo.Client
}

func NewService(
	userRepository userdomain.Repository,
	chatRepository chatdomain.Repository,
	centrifugoClient *centrifugo.Client,
) *Service {
	return &Service{
		userRepository:   userRepository,
		chatRepository:   chatRepository,
		centrifugoClient: centrifugoClient,
	}
}

func (s *Service) CreateMessage(ctx context.Context, userID, peerID, peerType, text string) (*chatdomain.Message, error) {
	if err := s.checkPeer(ctx, peerID, peerType); err != nil {
		return nil, err
	}

	now := time.Now()

	message := &chatdomain.Message{
		ID:        uuid.NewString(),
		UserID:    userID,
		PeerID:    peerID,
		PeerType:  peerType,
		PeerHash:  chatdomain.CalculatePeerHash(userID, peerID, peerType),
		Text:      text,
		CreatedAt: now,
	}

	if _, err := s.chatRepository.Create(ctx, message); err != nil {
		return nil, err
	}

	switch peerType {
	case chatdomain.MessagePeerTypeUser:
		_, err := s.centrifugoClient.
			Publish(ctx, chatdomain.ChannelNewMessage(message.PeerID), message.DTO())
		if err != nil {
			return nil, err
		}

	default:
		return nil, domain.ErrNotImplemented()
	}

	return message, nil
}

func (s *Service) ListMessages(
	ctx context.Context,
	userID string,
	peerID, peerType string,
	offset, limit int64,
) ([]chatdomain.Message, error) {
	messages, err := s.chatRepository.List(
		ctx,
		chatdomain.CalculatePeerHash(userID, peerID, peerType),
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, chatdomain.ErrMessageNotFound()
	}

	return messages, nil
}

func (s *Service) ListRecentMessages(
	ctx context.Context,
	userID string,
	peerType string,
	offset, limit int64,
) ([]chatdomain.Message, error) {
	messages, err := s.chatRepository.ListRecent(ctx, userID, peerType, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, chatdomain.ErrMessageNotFound()
	}

	return messages, nil
}

func (s *Service) UpdateMessage(ctx context.Context, userID, peerID, peerType string, id, text string) (*chatdomain.Message, error) {
	peerHash := chatdomain.CalculatePeerHash(userID, peerID, peerType)

	updatedMessage, err := s.chatRepository.Update(ctx, userID, peerHash, id, text)
	if err != nil {
		return nil, err
	}

	if updatedMessage == nil {
		return nil, chatdomain.ErrMessageNotFound()
	}

	switch updatedMessage.PeerType {
	case chatdomain.MessagePeerTypeUser:
		_, err := s.centrifugoClient.
			Publish(ctx, chatdomain.ChannelUpdatedMessage(updatedMessage.PeerID), updatedMessage.DTO())
		if err != nil {
			return nil, err
		}

	default:
		return nil, domain.ErrNotImplemented()
	}

	return updatedMessage, nil
}

func (s *Service) ReadMessages(ctx context.Context, userID, peerID, peerType string) error {
	count, err := s.chatRepository.UpdateRead(ctx, userID, chatdomain.CalculatePeerHash(userID, peerID, peerType))
	if err != nil {
		return err
	}

	if count == 0 {
		return chatdomain.ErrMessageNotFound()
	}

	switch peerType {
	case chatdomain.MessagePeerTypeUser:
		_, err := s.centrifugoClient.
			Publish(ctx, chatdomain.ChannelReadMessage(peerID), chatdomain.ReadDTO{
				UserID:   userID,
				PeerID:   peerID,
				PeerType: peerType,
			})
		if err != nil {
			return err
		}

	default:
		return domain.ErrNotImplemented()
	}

	return nil
}

func (s *Service) UpdateTyping(ctx context.Context, userID, peerID, peerType string, typing bool) error {
	if err := s.checkPeer(ctx, peerID, peerType); err != nil {
		return err
	}

	switch peerType {
	case chatdomain.MessagePeerTypeUser:
		_, err := s.centrifugoClient.Publish(ctx, chatdomain.ChannelTyping(peerID), chatdomain.TypingDTO{
			UserID:   userID,
			PeerID:   peerID,
			PeerType: peerType,
			Typing:   typing,
		})
		if err != nil {
			return err
		}

	default:
		return domain.ErrNotImplemented()
	}

	return nil
}

func (s *Service) checkPeer(ctx context.Context, peerID, peerType string) error {
	switch peerType {
	case chatdomain.MessagePeerTypeUser:
		exists, err := s.userRepository.Exists(ctx, peerID)
		if err != nil {
			return err
		}

		if !exists {
			return userdomain.ErrUserNotFound()
		}

	default:
		return domain.ErrNotImplemented()
	}

	return nil
}
