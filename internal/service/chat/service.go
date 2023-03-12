package chatservice

import (
	"context"
	"github.com/undefined7887/harmony-backend/internal/domain"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"time"
)

type Service struct {
	userRepository    userdomain.Repository
	messageRepository chatdomain.MessageRepository
	chatRepository    chatdomain.ChatRepository

	centrifugoClient *centrifugo.Client
}

func NewService(
	userRepository userdomain.Repository,
	messageRepository chatdomain.MessageRepository,
	chatRepository chatdomain.ChatRepository,
	centrifugoClient *centrifugo.Client,
) *Service {
	return &Service{
		userRepository:    userRepository,
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
		centrifugoClient:  centrifugoClient,
	}
}

func (s *Service) CreateMessage(ctx context.Context, userID, peerID, peerType, text string) (*chatdomain.Message, error) {
	if err := s.checkPeer(ctx, peerID, peerType); err != nil {
		return nil, err
	}

	now := time.Now()

	message := &chatdomain.Message{
		ID:          domain.ID(),
		UserID:      userID,
		PeerID:      peerID,
		PeerType:    peerType,
		ChatID:      chatdomain.ChatID(userID, peerID, peerType),
		Text:        text,
		ReadUserIDs: []string{}, // Required to be not nil
		CreatedAt:   now,
	}

	if _, err := s.messageRepository.Create(ctx, message); err != nil {
		return nil, err
	}

	switch peerType {
	case chatdomain.PeerTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelMessageCreated, message.PeerID)

		if _, err := s.centrifugoClient.Publish(ctx, channel, message.DTO()); err != nil {
			return nil, err
		}

	default:
		return nil, domain.ErrNotImplemented()
	}

	return message, nil
}

func (s *Service) GetMessage(ctx context.Context, userID, id string) (*chatdomain.Message, error) {
	message, err := s.messageRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, chatdomain.ErrMessageNotFound()
	}

	switch message.PeerType {
	case chatdomain.PeerTypeUser:
		chatID := chatdomain.ChatID(userID, message.PeerID, message.PeerType)

		// Checking that user is participant of this chat
		if chatID != message.ChatID {
			return nil, domain.ErrForbidden()
		}

	case chatdomain.PeerTypeGroup:
		return nil, domain.ErrNotImplemented()
	}

	return message, nil
}

func (s *Service) ListMessages(
	ctx context.Context,
	userID, peerID, peerType string,
	offset, limit int64,
) ([]chatdomain.Message, error) {
	chatID := chatdomain.ChatID(userID, peerID, peerType)

	messages, err := s.messageRepository.List(ctx, chatID, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, chatdomain.ErrMessageNotFound()
	}

	return messages, nil
}

func (s *Service) UpdateMessage(ctx context.Context, userID, id, text string) (*chatdomain.Message, error) {
	updatedMessage, err := s.messageRepository.UpdateText(ctx, id, userID, text)
	if err != nil {
		return nil, err
	}

	if updatedMessage == nil {
		return nil, chatdomain.ErrMessageNotFound()
	}

	switch updatedMessage.PeerType {
	case chatdomain.PeerTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelMessageUpdated, updatedMessage.PeerID)

		if _, err := s.centrifugoClient.Publish(ctx, channel, updatedMessage.DTO()); err != nil {
			return nil, err
		}

	default:
		return nil, domain.ErrNotImplemented()
	}

	return updatedMessage, nil
}

func (s *Service) ListChats(
	ctx context.Context,
	userID, peerType string,
	offset, limit int64,
) ([]chatdomain.Chat, error) {
	chats, err := s.chatRepository.List(ctx, userID, peerType, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(chats) == 0 {
		return nil, chatdomain.ErrChatsNotFound()
	}

	return chats, nil
}

func (s *Service) UpdateChatRead(ctx context.Context, userID, peerID, peerType string) error {
	chatID := chatdomain.ChatID(userID, peerID, peerType)

	count, err := s.chatRepository.UpdateRead(ctx, userID, chatID)
	if err != nil {
		return err
	}

	if count == 0 {
		return chatdomain.ErrMessageNotFound()
	}

	switch peerType {
	case chatdomain.PeerTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelRead, peerID)

		_, err := s.centrifugoClient.Publish(ctx, channel, chatdomain.ReadDTO{
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

func (s *Service) UpdateChatTyping(ctx context.Context, userID, peerID, peerType string, typing bool) error {
	if err := s.checkPeer(ctx, peerID, peerType); err != nil {
		return err
	}

	switch peerType {
	case chatdomain.PeerTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelTyping, peerID)

		_, err := s.centrifugoClient.Publish(ctx, channel, chatdomain.TypingDTO{
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
	case chatdomain.PeerTypeUser:
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
