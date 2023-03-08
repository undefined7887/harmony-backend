package chatdomain

import (
	"github.com/undefined7887/harmony-backend/internal/domain"
	"time"
)

const (
	ChatTypeUser  = "user"
	ChatTypeGroup = "group"
)

type Chat struct {
	ID string `bson:"_id"`

	Type string `bson:"type"`
	Name string `bson:"name,omitempty"`

	// Information about messages in chat, optional
	Message *ChatMessage `bson:"message,omitempty"`

	CreatedAt *time.Time `bson:"created_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

type ChatMessage struct {
	Last Message `bson:"last"`

	// Total count of unread messages
	UnreadCount int64 `bson:"unread_count"`
}

type Message struct {
	ID string `bson:"_id"`

	UserID string `bson:"user_id"`
	PeerID string `bson:"peer_id"`

	// For group messages - group id
	// For private messages - combination of user ids
	ChatID   string `bson:"chat_id"`
	ChatType string `bson:"chat_type"`

	Text        string   `bson:"text"`
	Attachments []string `bson:"attachments,omitempty"`

	// Users, who read this message
	UserReadIDs []string `bson:"user_read_ids"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

func ChatID(userID, peerID, chatType string) string {
	var chatID string

	switch chatType {
	case ChatTypeUser:
		chatID = domain.CombineIDs(userID, peerID)

	case ChatTypeGroup:
		chatID = peerID
	}

	return chatID
}
