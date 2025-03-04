package engine

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type Engine struct {
	log  *zap.Logger
	data map[string]string
	mx   sync.RWMutex
}

func NewEngine(log *zap.Logger) *Engine {
	return &Engine{
		log:  log,
		data: map[string]string{},
	}
}

func (s *Engine) Get(ctx context.Context, key string) string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if data, ok := s.data[key]; ok {
		return data
	} else {
		return errors.New("not found").Error()
	}
}

func (s *Engine) Set(ctx context.Context, key, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[key] = value
}

func (s *Engine) Delete(ctx context.Context, key string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, key)
}
