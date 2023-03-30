package calldomain

import (
	"encoding/json"
	"time"
)

const (
	StatusRequest  = "request"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
)

type Call struct {
	ID string `bson:"_id"`

	UserID string `bson:"user_id"`
	PeerID string `bson:"peer_id"`

	Status string `bson:"status"`

	// WebRTC data
	WebRTC CallWebRTC `bson:"web_rtc"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CallWebRTC struct {
	Offer  json.RawMessage `bson:"offer,omitempty"`
	Answer json.RawMessage `bson:"answer,omitempty"`
}
