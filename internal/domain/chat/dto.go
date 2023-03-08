package chatdomain

import "time"

type ChatDTO struct {
	ID string `json:"id"`

	Type string `json:"type"`
	Name string `json:"name,omitempty"`

	Message ChatMessageDTO `json:"message"`
}

type ChatMessageDTO struct {
	Last        MessageDTO `json:"last"`
	UnreadCount int64      `json:"unread_count"`
}

func (c *Chat) DTO() ChatDTO {
	return ChatDTO{
		ID:   c.ID,
		Type: c.Type,
		Name: c.Name,

		Message: ChatMessageDTO{
			Last:        c.Message.Last.DTO(),
			UnreadCount: c.Message.UnreadCount,
		},
	}
}

type MessageDTO struct {
	ID string `json:"_id"`

	UserID string `json:"user_id"`
	PeerID string `json:"peer_id"`

	ChatID   string `json:"chat_id"`
	ChatType string `json:"chat_type"`

	Text        string   `json:"text"`
	Attachments []string `json:"attachments,omitempty"`

	// Users, who read this message
	UserReadIDs []string `json:"user_read_ids,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m *Message) DTO() MessageDTO {
	return MessageDTO{
		ID:          m.ID,
		UserID:      m.UserID,
		PeerID:      m.PeerID,
		ChatID:      m.ChatID,
		ChatType:    m.ChatType,
		Text:        m.Text,
		Attachments: m.Attachments,
		UserReadIDs: m.UserReadIDs,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type ReadDTO struct {
	UserID string `json:"user_id"`
	PeerID string `json:"peer_id"`

	CharID   string `json:"char_id"`
	ChatType string `json:"chat_type"`
}

type TypingDTO struct {
	UserID string `json:"user_id"`
	PeerID string `json:"peer_id"`

	ChatID   string `json:"char_id"`
	ChatType string `json:"chat_type"`

	Typing bool `json:"typing"`
}
