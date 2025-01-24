package service

import (
	"context"

	"concurrency/app/compute"
	"go.uber.org/zap"
)

type Service struct {
	engine  Engine
	compute Compute
	log     *zap.Logger
}

type Engine interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value string)
	Delete(ctx context.Context, key string)
}

type Compute interface {
	Parse(ctx context.Context, query string) (compute.Query, error)
}

func New(engine Engine, compute Compute, log *zap.Logger) *Service {
	return &Service{
		engine:  engine,
		compute: compute,
		log:     log,
	}
}

func (s *Service) Handler(ctx context.Context, queryStr string) string {
	query, err := s.compute.Parse(ctx, queryStr)
	if err != nil {
		s.log.Error("failed to parse query", zap.String("query", queryStr), zap.Error(err))
		return err.Error()
	}
	switch query.Command {
	case compute.SET:
		s.engine.Set(ctx, query.Args[0], query.Args[1])
		return ""
	case compute.GET:
		return s.engine.Get(ctx, query.Args[0])
	case compute.DEL:
		s.engine.Delete(ctx, query.Args[0])
		return ""
	}
	return "command unknown"
}
