package chatservice

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/undefined7887/harmony-backend/internal/domain"
	chatdomain "github.com/undefined7887/harmony-backend/internal/domain/chat"
	userdomain "github.com/undefined7887/harmony-backend/internal/domain/user"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	"github.com/undefined7887/harmony-backend/internal/repository"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"github.com/undefined7887/harmony-backend/internal/util"
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

func (s *Service) CreateMessage(ctx context.Context, userID, peerID, peerType, text string) (chatdomain.MessageDTO, error) {
	if err := s.checkPeer(ctx, peerID, peerType); err != nil {
		return chatdomain.MessageDTO{}, err
	}

	now := time.Now()

	message := chatdomain.Message{
		ID:          domain.ID(),
		UserID:      userID,
		PeerID:      peerID,
		PeerType:    peerType,
		ChatID:      chatdomain.ChatID(userID, peerID, peerType),
		Text:        text,
		ReadUserIDs: []string{}, // Required to be not nil
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if _, err := s.messageRepository.Create(ctx, &message); err != nil {
		return chatdomain.MessageDTO{}, err
	}

	switch peerType {
	case chatdomain.PeerTypeUser:
		s.centrifugoPublish(
			ctx,
			chatdomain.ChannelMessageNew(userID),
			chatdomain.NewMessageNotification{
				MessageDTO: chatdomain.MapMessageDTO(message),
			},
		)

		s.centrifugoPublish(
			ctx,
			chatdomain.ChannelMessageNew(peerID),
			chatdomain.NewMessageNotification{
				MessageDTO: chatdomain.MapMessageDTO(message),
			},
		)

	default:
		return chatdomain.MessageDTO{}, domain.ErrNotImplemented()
	}

	return chatdomain.MapMessageDTO(message), nil
}

func (s *Service) GetMessage(ctx context.Context, userID, id string) (chatdomain.MessageDTO, error) {
	message, err := s.messageRepository.Get(ctx, id)
	if repository.IsNoDocumentsErr(err) {
		return chatdomain.MessageDTO{}, chatdomain.ErrMessageNotFound()
	}

	if err != nil {
		return chatdomain.MessageDTO{}, err
	}

	switch message.PeerType {
	case chatdomain.PeerTypeUser:
		chatID := chatdomain.ChatID(userID, message.PeerID, message.PeerType)

		// Checking that user is participant of this chat
		if chatID != message.ChatID {
			return chatdomain.MessageDTO{}, domain.ErrForbidden()
		}

	case chatdomain.PeerTypeGroup:
		return chatdomain.MessageDTO{}, domain.ErrNotImplemented()
	}

	return chatdomain.MapMessageDTO(message), nil
}

func (s *Service) ListMessages(
	ctx context.Context,
	userID, peerID, peerType string,
	offset, limit int64,
) ([]chatdomain.MessageDTO, error) {
	chatID := chatdomain.ChatID(userID, peerID, peerType)

	messages, err := s.messageRepository.List(ctx, chatID, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, chatdomain.ErrMessageNotFound()
	}

	return util.Map(messages, chatdomain.MapMessageDTO), nil
}

func (s *Service) UpdateMessage(ctx context.Context, userID, id, text string) (chatdomain.MessageDTO, error) {
	updatedMessage, err := s.messageRepository.UpdateText(ctx, id, userID, text)
	if repository.IsNoDocumentsErr(err) {
		return chatdomain.MessageDTO{}, chatdomain.ErrMessageNotFound()
	}

	if err != nil {
		return chatdomain.MessageDTO{}, err
	}

	switch updatedMessage.PeerType {
	case chatdomain.PeerTypeUser:
		s.centrifugoPublish(
			ctx,
			chatdomain.ChannelMessageNew(updatedMessage.UserID),
			chatdomain.NewMessageNotification{
				MessageDTO: chatdomain.MapMessageDTO(updatedMessage),
			},
		)

		s.centrifugoPublish(
			ctx,
			chatdomain.ChannelMessageUpdates(updatedMessage.PeerID),
			chatdomain.UpdateMessageNotification{
				MessageDTO: chatdomain.MapMessageDTO(updatedMessage),
			},
		)

	default:
		return chatdomain.MessageDTO{}, domain.ErrNotImplemented()
	}

	return chatdomain.MapMessageDTO(updatedMessage), nil
}

func (s *Service) ListChats(
	ctx context.Context,
	userID, peerType string,
	offset, limit int64,
) ([]chatdomain.ChatDTO, error) {
	chats, err := s.chatRepository.List(ctx, userID, peerType, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(chats) == 0 {
		return nil, chatdomain.ErrChatsNotFound()
	}

	return util.Map(chats, chatdomain.MapChatDTO), nil
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
		data := chatdomain.UpdateChatReadNotification{
			UserID:   userID,
			PeerID:   peerID,
			PeerType: peerType,
		}

		// Publishing for current user
		s.centrifugoPublish(ctx, chatdomain.ChannelReadUpdates(userID), data)

		// Publishing for peer
		s.centrifugoPublish(ctx, chatdomain.ChannelReadUpdates(peerID), data)

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
		s.centrifugoPublish(
			ctx,
			chatdomain.ChannelTypingUpdates(peerID),
			chatdomain.UpdateChatTypingNotification{
				UserID:   userID,
				PeerID:   peerID,
				PeerType: peerType,
				Typing:   typing,
			},
		)

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
