package memory

import (
	"context"
	"github.com/gerruda/simple_gatekeeper/internal/ports"
	"sync"
)

var (
	_ ports.PendingRepository = (*Memory)(nil)
)

type key struct{ chatID, userID int64 }

type Memory struct {
	mutex sync.Map
}

func NewMemAdapter() *Memory {
	return &Memory{}
}

func (s *Memory) Add(ctx context.Context, chatID, userID int64) error {
	s.mutex.Store(key{chatID, userID}, 0)

	return nil
}

func (s *Memory) Remove(ctx context.Context, chatID, userID int64) error {
	s.mutex.Delete(key{chatID, userID})

	return nil
}

func (s *Memory) Exists(ctx context.Context, chatID, userID int64) (bool, error) {
	_, ok := s.mutex.Load(key{chatID, userID})

	return ok, nil
}
