package ports

import (
	"context"
)

type PendingRepository interface {
	Add(ctx context.Context, chatID, userID int64) error
	Remove(ctx context.Context, chatID, userID int64) error
	Exists(ctx context.Context, chatID, userID int64) (bool, error)
}
