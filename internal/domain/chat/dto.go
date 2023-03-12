package chatdomain

import (
	"time"
)

type ChatDTO struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name,omitempty"`
	Message     MessageDTO `json:"message"`
	UnreadCount int64      `json:"unread_count"`
}

func (c *Chat) DTO() ChatDTO {
	return ChatDTO{
		ID:          c.ID,
		Type:        c.Type,
		Name:        c.Name,
		Message:     c.Message.DTO(),
		UnreadCount: c.UnreadCount,
	}
}

type MessageDTO struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	PeerID      string     `json:"peer_id"`
	PeerType    string     `json:"peer_type"`
	Text        string     `json:"text"`
	Attachments []string   `json:"attachments,omitempty"`
	ReadUserIDs []string   `json:"read_user_ids,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (m *Message) DTO() MessageDTO {
	return MessageDTO{
		ID:          m.ID,
		UserID:      m.UserID,
		PeerID:      m.PeerID,
		PeerType:    m.PeerType,
		Text:        m.Text,
		Attachments: m.Attachments,
		ReadUserIDs: m.ReadUserIDs,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type ReadDTO struct {
	UserID   string `json:"user_id"`
	PeerID   string `json:"peer_id"`
	PeerType string `json:"peer_type"`
}

type TypingDTO struct {
	UserID   string `json:"user_id"`
	PeerID   string `json:"peer_id"`
	PeerType string `json:"peer_type"`

	Typing bool `json:"typing"`
}
