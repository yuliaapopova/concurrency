package wal

import (
	"context"
	"errors"
	"sync"
	"time"

	"concurrency/app/common"
	"concurrency/app/compute"
	"concurrency/app/config"
	"concurrency/app/storage/filesystem"
	"go.uber.org/zap"
)

type logManager interface {
	AppendLogs([]Log)
	Load() ([]Log, error)
}

type WAL struct {
	mx         sync.Mutex
	logManager logManager

	flushTimeout time.Duration
	batchSize    int

	batchChan chan []Log
	batch     []Log

	status chan error
}

func NewWAL(cfg *config.WAL, logger *zap.Logger) (*WAL, error) {
	if cfg == nil {
		return nil, nil
	}
	if logger == nil {
		return nil, errors.New("nil logger")
	}
	maxSegmentSize, err := common.ParseMessageSize(cfg.MaxSegmentSize)
	if err != nil {
		return nil, err
	}

	fs := filesystem.NewSegment(cfg.DataDirectory, maxSegmentSize)
	lm := NewLogManager(fs, logger)

	return &WAL{
		logManager:   lm,
		flushTimeout: cfg.FlushingBatchTimeout,
		batchSize:    cfg.FlushingBatchSize,
		batchChan:    make(chan []Log, 1),
	}, nil
}

func (w *WAL) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.flushTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			default:
			}

			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			case batch := <-w.batchChan:
				w.logManager.AppendLogs(batch)
				ticker.Reset(w.flushTimeout)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Set(ctx context.Context, key, value string) error {
	w.write(ctx, compute.SET.Int(), []string{key, value})
	return <-w.status
}

func (w *WAL) Delete(ctx context.Context, key string) error {
	w.write(ctx, compute.DEL.Int(), []string{key})
	return <-w.status
}

func (w *WAL) write(ctx context.Context, command int, args []string) {
	id := ctx.Value("ID").(uint64)
	record := NewLog(id, command, args)

	w.mx.Lock()
	w.batch = append(w.batch, record)
	if len(w.batch) >= w.batchSize {
		w.batchChan <- w.batch
		w.batch = nil
	}
	w.mx.Unlock()

	w.status = record.Status
}

func (w *WAL) flushBatch() {
	w.mx.Lock()
	defer w.mx.Unlock()
	if len(w.batch) != 0 {
		w.logManager.AppendLogs(w.batch)
		w.batch = nil
	}
}

func (w *WAL) Recover() ([]Log, error) {
	return w.logManager.Load()
}
