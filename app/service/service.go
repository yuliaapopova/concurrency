package service

import (
	"context"

	"concurrency/app/compute"
	"concurrency/app/config"
	"concurrency/app/network"
	storage "concurrency/app/storage"
	"concurrency/app/storage/engine"
	"concurrency/app/storage/wal"
	"go.uber.org/zap"
)

type Service struct {
	config  *config.Config
	storage Storage
	compute Compute
	log     *zap.Logger
}

type Storage interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value string)
	Delete(ctx context.Context, key string)
}

type Compute interface {
	Parse(ctx context.Context, query string) (compute.Query, error)
}

func New(config *config.Config, storage Storage, compute Compute, log *zap.Logger) *Service {
	return &Service{
		config:  config,
		storage: storage,
		compute: compute,
		log:     log,
	}
}

func Start(ctx context.Context, config *config.Config) {
	logger, _ := zap.NewProduction()

	queryParser := compute.New(logger)
	eng := engine.NewEngine(logger)
	writeLog, err := wal.NewWAL(config.WAL, logger)
	if err != nil {
		logger.Fatal("Failed to create wal", zap.Error(err))
	}
	db, err := storage.NewStorage(logger, eng, writeLog)
	if err != nil {
		logger.Fatal("Failed to create storage", zap.Error(err))
	}

	s := New(config, db, queryParser, logger)
	net, err := network.NewTCPServer(ctx, config.Network, logger)
	if err != nil {
		logger.Fatal("Failed to create tcp server", zap.Error(err))
	}

	if writeLog != nil {
		go writeLog.Start(ctx)
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
		s.storage.Set(ctx, query.Args[0], query.Args[1])
		return "[ok]"
	case compute.GET:
		return s.storage.Get(ctx, query.Args[0])
	case compute.DEL:
		s.storage.Delete(ctx, query.Args[0])
		return "[ok]"
	}
	return "command unknown"
}
