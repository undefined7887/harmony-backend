package chatdomain

import (
	"github.com/undefined7887/harmony-backend/internal/util"
	"strings"
	"time"
)

const (
	MessagePeerTypeUser  = "user"
	MessagePeerTypeGroup = "group"
)

type Message struct {
	ID string `bson:"_id"`

	// Sender of message
	UserID string `bson:"user_id"`

	// Receiver of message
	PeerID   string `bson:"peer_id"`
	PeerType string `bson:"peer_type"`

	// Grouping factor
	// - For users: join(sort(user_from_id, user_peer_id))
	// - For groups: group_id
	PeerHash string `bson:"peer_hash"`

	Text string `bson:"text"`
	Read bool   `bson:"read"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}

func (m *Message) DTO() MessageDTO {
	return MessageDTO{
		ID:        m.ID,
		UserID:    m.UserID,
		PeerID:    m.PeerID,
		PeerType:  m.PeerType,
		Text:      m.Text,
		Read:      m.Read,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func CalculatePeerHash(userID, peerID, peerType string) string {
	var peerHash string

	switch peerType {
	case MessagePeerTypeUser:
		peerHash = strings.Join(util.SortSequence(userID, peerID), "_")

	case MessagePeerTypeGroup:
		peerHash = peerID
	}

	return peerHash
}
