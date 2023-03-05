package chatdomain

import "time"

type MessageDTO struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PeerID    string     `json:"peer_id"`
	PeerType  string     `json:"peer_type"`
	Text      string     `json:"text"`
	Read      bool       `json:"read"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type TypingDTO struct {
	UserID   string `json:"user_id"`
	PeerID   string `json:"peer_id"`
	PeerType string `json:"peer_type"`
	Typing   bool   `json:"typing"`
}
