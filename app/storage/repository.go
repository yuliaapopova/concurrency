package storage

import (
	"context"
	"errors"
	"sync"

	"concurrency/app/compute"
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

func (s *Storage) QueryHandler(ctx context.Context, query compute.Query) string {
	switch query.Command {
	case compute.SET:
		return s.Set(ctx, query.Args)
	case compute.GET:
		return s.Get(ctx, query.Args[0])
	case compute.DEL:
		return s.Delete(ctx, query.Args[0])
	}
	return "command unknown"
}

func (s *Storage) Get(ctx context.Context, key string) string {
	if data, ok := s.data[key]; ok {
		return data
	} else {
		return errors.New("not found").Error()
	}
}

func (s *Storage) Set(ctx context.Context, args []string) string {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[args[0]] = args[1]
	return "ok"
}

func (s *Storage) Delete(ctx context.Context, key string) string {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, key)
	return "ok"
}
