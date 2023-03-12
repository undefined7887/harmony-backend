package userdomain

import (
	"time"
)

const (
	StatusOnline  = "online"
	StatusAway    = "away"
	StatusSilence = "silence"
	StatusOffline = "offline"
)

var (
	UserPingInterval    = time.Second * 30
	UserOutdatedTimeout = time.Minute
)

type User struct {
	ID     string `bson:"_id"`
	Status string `bson:"status"`

	Email    string `bson:"email"`
	Photo    string `bson:"photo"`
	Nickname string `bson:"nickname"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
