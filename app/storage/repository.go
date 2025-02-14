package storage

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type Storage struct {
	log  *zap.Logger
	data map[string]string
	mx   sync.RWMutex
}

func NewStorage(log *zap.Logger) *Storage {
	return &Storage{
		log:  log,
		data: map[string]string{},
	}
}

func (s *Storage) Get(ctx context.Context, key string) string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if data, ok := s.data[key]; ok {
		return data
	} else {
		return errors.New("not found").Error()
	}
}

func (s *Storage) Set(ctx context.Context, key, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[key] = value
}

func (s *Storage) Delete(ctx context.Context, key string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, key)
}
