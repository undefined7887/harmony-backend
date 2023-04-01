package calldomain

import (
	"time"
)

const (
	StatusRequest  = "request"
	StatusAccepted = "accepted"
	StatusDeclined = "declined"
	StatusFinished = "finished"
)

type Call struct {
	ID string `bson:"_id"`

	UserID string `bson:"user_id"`
	PeerID string `bson:"peer_id"`

	Status string `bson:"status"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
