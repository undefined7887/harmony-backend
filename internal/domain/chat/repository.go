package chatdomain

import "context"

type Repository interface {
	Create(ctx context.Context, message *Message) (bool, error)

	// List returns messages from chat
	List(
		ctx context.Context,
		peerHash string,
		offset, limit int64,
	) ([]Message, error)

	// ListRecent returns most recent message in each chat for user
	ListRecent(
		ctx context.Context,
		userID, peerType string,
		offset, limit int64,
	) ([]Message, error)

	// Update updates single message contents
	Update(ctx context.Context, userID, peerHash, id, text string) (*Message, error)

	// UpdateRead special function to set read status to true for all messages in chat
	UpdateRead(ctx context.Context, userID, peerHash string) (int64, error)
}
