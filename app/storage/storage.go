package storage

import (
	"context"
	"errors"
	"fmt"

	"concurrency/app/compute"
	"concurrency/app/storage/wal"
	"go.uber.org/zap"
)

type Engine interface {
	Get(ctx context.Context, key string) string
	Set(ctx context.Context, key string, value string)
	Delete(ctx context.Context, key string)
}

type WAL interface {
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	Recover() ([]wal.Log, error)
}

type Storage struct {
	logger *zap.Logger

	engine      Engine
	wal         WAL
	generatorID *IDGenerator

	isReplicaSlave bool
	stream         chan []wal.Log
}

func NewStorage(logger *zap.Logger, en Engine, wal WAL, isReplicaSlave bool, stream chan []wal.Log) (*Storage, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	if en == nil {
		return nil, errors.New("engine required")
	}

	storage := &Storage{
		logger:         logger,
		engine:         en,
		isReplicaSlave: isReplicaSlave,
		stream:         stream,
	}

	var lastLSN uint64
	if wal != nil {
		storage.wal = wal
		logs, err := wal.Recover()
		if err != nil {
			return nil, err
		}
		lastLSN = storage.applyData(logs)
	}

	storage.generatorID = NewIDGenerator(lastLSN)

	if stream != nil {
		go func() {
			for logs := range stream {
				storage.applyData(logs)
			}
		}()
	}

	return storage, nil
}

func (s *Storage) Get(ctx context.Context, key string) string {
	return s.engine.Get(ctx, key)
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	if s.isReplicaSlave {
		return errors.New("command SET is not available in slave")
	}
	ID := s.generatorID.NextID()
	nexCtx := context.WithValue(ctx, "ID", ID)
	if s.wal != nil {
		s.wal.Set(nexCtx, key, value)
		fmt.Println("wal finished")
	}
	s.engine.Set(nexCtx, key, value)
	return nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if s.isReplicaSlave {
		return errors.New("command DEL is not available in slave")
	}
	ID := s.generatorID.NextID()
	nexCtx := context.WithValue(ctx, "ID", ID)
	if s.wal != nil {
		s.wal.Delete(nexCtx, key)
	}
	s.engine.Delete(nexCtx, key)
	return nil
}

func (s *Storage) applyData(logs []wal.Log) uint64 {
	var lastLSN uint64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := context.WithValue(context.Background(), "ID", log.LSN)

		switch log.Command {
		case compute.SET.Int():
			s.engine.Set(ctx, log.Args[0], log.Args[1])
		case compute.DEL.Int():
			s.engine.Delete(ctx, log.Args[0])
		}
	}

	return lastLSN
}
