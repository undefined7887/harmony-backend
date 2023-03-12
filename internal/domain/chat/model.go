package chatdomain

import (
	"github.com/undefined7887/harmony-backend/internal/domain"
	"time"
)

const (
	PeerTypeUser  = "user"
	PeerTypeGroup = "group"
)

type Chat struct {
	ID   string `bson:"_id"`
	Type string `bson:"type"`

	// For future use in groups
	Name string `bson:"name,omitempty"`

	// Last message from chat
	Message Message `bson:"message"`

	// Unread messages count in chat
	UnreadCount int64 `bson:"unread_count"`
}

type Message struct {
	ID     string `bson:"_id"`
	UserID string `bson:"user_id"`

	// Peer of message
	PeerID   string `bson:"peer_id"`
	PeerType string `bson:"peer_type"`

	// For users messages - combination of user ids
	// For group messages - group id
	// Internal use only
	ChatID string `bson:"chat_id"`

	Text        string   `bson:"text"`
	Attachments []string `bson:"attachments,omitempty"`

	// Users, who read this message
	ReadUserIDs []string `bson:"read_user_ids"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

func ChatID(userID, peerID, peerType string) string {
	var chatID string

	switch peerType {
	case PeerTypeUser:
		chatID = domain.CombineIDs(userID, peerID)

	case PeerTypeGroup:
		chatID = peerID
	}

	return chatID
}
