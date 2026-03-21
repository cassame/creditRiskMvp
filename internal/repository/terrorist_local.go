package repository

import (
	"context"
	"sync"
)

type LocalTerroristStore struct {
	mu  sync.RWMutex
	set map[string]struct{}
}

func NewLocalTerroristStore() *LocalTerroristStore {
	return &LocalTerroristStore{
		set: make(map[string]struct{}),
	}
}

func (s *LocalTerroristStore) IsTerrorist(ctx context.Context, passport string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, found := s.set[passport]
	return found, nil
}

func (s *LocalTerroristStore) UpdateList(ctx context.Context, passports []string) error {
	newSet := make(map[string]struct{})
	for _, p := range passports {
		newSet[p] = struct{}{}
	}
	s.mu.Lock()
	s.set = newSet
	s.mu.Unlock()
	return nil
}
