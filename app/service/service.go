package service

import (
	"context"
	"sync"

	"concurrency/app/compute"
	"concurrency/app/config"
	"concurrency/app/network"
	storage "concurrency/app/storage"
	"concurrency/app/storage/engine"
	"concurrency/app/storage/replication"
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
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
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
	eng := engine.NewEngine(logger, config.Engine.PartitionsNumber)
	WAL, err := wal.NewWAL(config.WAL, logger)
	if err != nil {
		logger.Fatal("Failed to create wal", zap.Error(err))
	}

	net, err := network.NewTCPServer(ctx, config.Network, logger)
	if err != nil {
		logger.Fatal("Failed to create tcp server", zap.Error(err))
	}

	if WAL != nil {
		go WAL.Start(ctx)
	}

	replica, err := replication.NewReplication(ctx, config.Replication, config.WAL, logger)
	if err != nil {
		logger.Fatal("Failed to create replication", zap.Error(err))
	}

	var isReplicaSlave bool
	var stream chan []wal.Log
	if replica != nil {
		if replica.Master != nil {
			stream = nil
		}
		if replica.Slave != nil {
			isReplicaSlave = true
			stream = replica.Slave.ReplicationStream()
		}
	}

	db, err := storage.NewStorage(logger, eng, WAL, isReplicaSlave, stream)
	if err != nil {
		logger.Fatal("Failed to create storage", zap.Error(err))
	}

	s := New(config, db, queryParser, logger)

	wg := &sync.WaitGroup{}
	if WAL != nil {
		wg.Add(1)
		if replica != nil && replica.Slave != nil {
			go func() {
				defer wg.Done()
				replica.Slave.Start(ctx)
			}()
		} else {
			go func() {
				defer func() {
					wg.Done()
				}()

				WAL.Start(ctx)
			}()
		}

		if replica != nil && replica.Master != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				replica.Master.Start(ctx)
			}()
		}
	}

	net.HandleQueries(ctx, func(ctx context.Context, bytes []byte) []byte {
		response := s.Handler(ctx, string(bytes))
		return []byte(response)
	})
	wg.Wait()
}

func (s *Service) Handler(ctx context.Context, queryStr string) string {
	query, err := s.compute.Parse(ctx, queryStr)
	if err != nil {
		s.log.Error("failed to parse query", zap.String("query", queryStr), zap.Error(err))
		return err.Error()
	}
	switch query.Command {
	case compute.SET:
		err = s.storage.Set(ctx, query.Args[0], query.Args[1])
		if err != nil {
			return err.Error()
		}
		return "[ok]"
	case compute.GET:
		return s.storage.Get(ctx, query.Args[0])
	case compute.DEL:
		err = s.storage.Delete(ctx, query.Args[0])
		if err != nil {
			return err.Error()
		}
		return "[ok]"
	}
	return "command unknown"
}
