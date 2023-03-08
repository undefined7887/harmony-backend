package chatdomain

import "context"

type MessageRepository interface {
	Create(ctx context.Context, message *Message) (bool, error)

	Get(ctx context.Context, id string) (*Message, error)
	List(ctx context.Context, chatID string, offset, limit int64) ([]Message, error)

	UpdateText(ctx context.Context, id, userID, text string) (*Message, error)
}

type ChatRepository interface {
	List(
		ctx context.Context,
		userID, chatType string,
		offset, limit int64,
	) ([]Chat, error)

	UpdateRead(ctx context.Context, userID, chatID string) (int64, error)
}
