package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Storage struct {
	log  *zap.Logger
	data map[string]string
}

func NewStorage(log *zap.Logger) *Storage {
	return &Storage{
		log:  log,
		data: map[string]string{},
	}
}

func (s *Storage) Get(ctx context.Context, key string) string {
	if data, ok := s.data[key]; ok {
		return data
	} else {
		return errors.New("not found").Error()
	}
}

func (s *Storage) Set(ctx context.Context, key, value string) {
	s.data[key] = value
}

func (s *Storage) Delete(ctx context.Context, key string) {
	delete(s.data, key)
}
