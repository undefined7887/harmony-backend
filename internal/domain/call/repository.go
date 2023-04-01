package calldomain

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, call *Call) (bool, error)

	Read(ctx context.Context, id, status string) (Call, error)
	ReadLast(ctx context.Context, userID, status string) (Call, error)

	UpdateStatus(ctx context.Context, id, peerID string, previousStatuses []string, newStatus string) (Call, error)
}
