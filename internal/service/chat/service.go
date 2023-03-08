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

func (s *Service) CreateMessage(ctx context.Context, userID, peerID, chatType, text string) (*chatdomain.Message, error) {
	if err := s.checkPeer(ctx, peerID, chatType); err != nil {
		return nil, err
	}

	now := time.Now()

	message := &chatdomain.Message{
		ID:          domain.ID(),
		UserID:      userID,
		PeerID:      peerID,
		ChatID:      chatdomain.ChatID(userID, peerID, chatType),
		ChatType:    chatType,
		Text:        text,
		UserReadIDs: []string{}, // Required to be not nil
		CreatedAt:   now,
	}

	if _, err := s.messageRepository.Create(ctx, message); err != nil {
		return nil, err
	}

	switch chatType {
	case chatdomain.ChatTypeUser:
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

	switch message.ChatType {
	case chatdomain.ChatTypeUser:
		chatID := chatdomain.ChatID(userID, message.PeerID, message.ChatType)

		// Checking that user is participant of this chat
		if chatID != message.ChatID {
			return nil, domain.ErrForbidden()
		}

	case chatdomain.ChatTypeGroup:
		return nil, domain.ErrNotImplemented()
	}

	return message, nil
}

func (s *Service) ListMessages(
	ctx context.Context,
	userID, peerID, chatType string,
	offset, limit int64,
) ([]chatdomain.Message, error) {
	chatID := chatdomain.ChatID(userID, peerID, chatType)

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

	switch updatedMessage.ChatType {
	case chatdomain.ChatTypeUser:
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
	userID, chatType string,
	offset, limit int64,
) ([]chatdomain.Chat, error) {
	chats, err := s.chatRepository.List(ctx, userID, chatType, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(chats) == 0 {
		return nil, chatdomain.ErrChatsNotFound()
	}

	return chats, nil
}

func (s *Service) UpdateChatRead(ctx context.Context, userID, peerID, chatType string) error {
	chatID := chatdomain.ChatID(userID, peerID, chatType)

	count, err := s.chatRepository.UpdateRead(ctx, userID, chatID)
	if err != nil {
		return err
	}

	if count == 0 {
		return chatdomain.ErrMessageNotFound()
	}

	switch chatType {
	case chatdomain.ChatTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelRead, peerID)

		_, err := s.centrifugoClient.Publish(ctx, channel, chatdomain.ReadDTO{
			UserID:   userID,
			PeerID:   peerID,
			CharID:   chatID,
			ChatType: chatType,
		})
		if err != nil {
			return err
		}

	default:
		return domain.ErrNotImplemented()
	}

	return nil
}

func (s *Service) UpdateChatTyping(ctx context.Context, userID, peerID, chatType string, typing bool) error {
	if err := s.checkPeer(ctx, peerID, chatType); err != nil {
		return err
	}

	chatID := chatdomain.ChatID(userID, peerID, chatType)

	switch chatType {
	case chatdomain.ChatTypeUser:
		channel := chatdomain.Channel(chatdomain.ChannelTyping, peerID)

		_, err := s.centrifugoClient.Publish(ctx, channel, chatdomain.TypingDTO{
			UserID:   userID,
			PeerID:   peerID,
			ChatID:   chatID,
			ChatType: chatType,
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

func (s *Service) checkPeer(ctx context.Context, peerID, chatType string) error {
	switch chatType {
	case chatdomain.ChatTypeUser:
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
