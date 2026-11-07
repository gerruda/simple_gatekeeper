package ports

import (
	"context"
)

type Keyboard struct {
	ButtonText string
	ButtonUrl  string
}

type Moderator interface {
	BunUser(ctx context.Context, chatID, userID int64) error
	UnbunUser(ctx context.Context, chatID, userID int64) error
}

type Greater interface {
	SendGreetingMsg(ctx context.Context, chatID, userID int64, text string) error
}

type Messenger interface {
	SndMsg(ctx context.Context, chatID int64, msg string) error
}

type AdminChecker interface {
	IsAdmin(ctx context.Context, chatID, userID int64) (bool, error)
}
