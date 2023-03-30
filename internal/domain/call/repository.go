package calldomain

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, call *Call) (bool, error)

	Read(ctx context.Context, id, status string) (Call, error)
	ReadLast(ctx context.Context, userID, status string) (Call, error)

	UpdateStatus(ctx context.Context, userID, id, status string, webRTC CallWebRTC) (Call, error)
}
