package service

import (
	"context"
	"fmt"

	"concurrency/app/compute"
	"concurrency/app/config"
	"concurrency/app/network"
	"concurrency/app/storage"
	"go.uber.org/zap"
)

type Service struct {
	config  *config.Config
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

func New(config *config.Config, engine Engine, compute Compute, log *zap.Logger) *Service {
	return &Service{
		config:  config,
		engine:  engine,
		compute: compute,
		log:     log,
	}
}

func Start(ctx context.Context, config *config.Config) {
	logger, _ := zap.NewProduction()
	engine := storage.NewStorage(logger)
	queryParser := compute.New(logger)
	s := New(config, engine, queryParser, logger)
	net, err := network.NewTCPServer(ctx, config.Network, logger)
	if err != nil {
		fmt.Println(err)
	}

	net.HandleQueries(ctx, func(ctx context.Context, bytes []byte) []byte {
		response := s.Handler(ctx, string(bytes))
		return []byte(response)
	})
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
		return "[ok]"
	case compute.GET:
		return s.engine.Get(ctx, query.Args[0])
	case compute.DEL:
		s.engine.Delete(ctx, query.Args[0])
		return "[ok]"
	}
	return "command unknown"
}
