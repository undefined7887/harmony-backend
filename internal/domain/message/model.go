package messagedomain

import "time"

type Chat struct {
	ID string `bson:"_id"`

	Name  string   `bson:"name,omitempty"` // Reserved for future use
	Users []string `bson:"users"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Message struct {
	ID string `bson:"_id"`

	Text   string `bson:"text"`
	Edited bool   `bson:"edited"`

	Attachments []string `bson:"attachments,omitempty"` // Reserved for future use

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}
