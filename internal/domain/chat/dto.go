package chatdomain

import (
	"time"

	"github.com/undefined7887/harmony-backend/internal/domain"
)

type ChatDTO struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name,omitempty"`
	Message     MessageDTO `json:"message"`
	UnreadCount int64      `json:"unread_count"`
}

func MapChatDTO(chat Chat) ChatDTO {
	return ChatDTO{
		ID:          chat.ID,
		Type:        chat.Type,
		Name:        chat.Name,
		Message:     MapMessageDTO(chat.Message),
		UnreadCount: chat.UnreadCount,
	}
}

type MessageDTO struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	PeerID      string    `json:"peer_id"`
	PeerType    string    `json:"peer_type"`
	Text        string    `json:"text"`
	Edited      bool      `json:"edited"`
	Attachments []string  `json:"attachments,omitempty"`
	ReadUserIDs []string  `json:"read_user_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MapMessageDTO(message Message) MessageDTO {
	return MessageDTO{
		ID:          message.ID,
		UserID:      message.UserID,
		PeerID:      message.PeerID,
		PeerType:    message.PeerType,
		Text:        message.Text,
		Edited:      message.Edited,
		Attachments: message.Attachments,
		ReadUserIDs: message.ReadUserIDs,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	}
}

type PeerParams struct {
	PeerID   string `uri:"peer_id" binding:"id"`
	PeerType string `uri:"peer_type" binding:"oneof=user group"`
}

// ---

type CreateMessageRequestBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

type CreateMessageResponse struct {
	MessageID string `json:"message_id"`
}

type NewMessageNotification struct {
	MessageDTO
}

// ---

type GetMessageResponse struct {
	MessageDTO
}

// ---

type ListMessagesResponse struct {
	Items []MessageDTO `json:"items"`
}

// ---

type UpdateMessageRequestBody struct {
	Text string `json:"text" binding:"min=1,max=1000"`
}

type UpdateMessageNotification struct {
	MessageDTO
}

// ---

type ListChatsRequestQuery struct {
	PeerType string `form:"peer_type" binding:"omitempty,oneof=user group"`

	domain.PaginationQuery
}

type ListChatsResponse struct {
	Items []ChatDTO `json:"items"`
}

// ---

type UpdateChatReadNotification struct {
	UserID   string `json:"user_id"`
	PeerID   string `json:"peer_id"`
	PeerType string `json:"peer_type"`
}

// ---

type UpdateChatTypingRequestBody struct {
	Typing bool `json:"typing"`
}

type UpdateChatTypingNotification struct {
	UserID   string `json:"user_id"`
	PeerID   string `json:"peer_id"`
	PeerType string `json:"peer_type"`

	Typing bool `json:"typing"`
}
