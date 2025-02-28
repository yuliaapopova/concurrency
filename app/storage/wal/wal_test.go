package wal

import (
	"context"
	"testing"
	"time"

	//"concurrency/app/compute"
	"concurrency/app/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewWAL(t *testing.T) {
	logger := zap.NewNop()

	cfg := config.WAL{
		FlushingBatchSize:    10,
		FlushingBatchTimeout: 50 * time.Millisecond,
		MaxSegmentSize:       "10KB",
		DataDirectory:        "testdata",
	}

	wal, err := NewWAL(&cfg, logger)
	assert.NoError(t, err)
	assert.NotNil(t, wal)
}

func TestWAL_Set(t *testing.T) {
	logger := zap.NewNop()

	cfg := config.WAL{
		FlushingBatchSize:    10,
		FlushingBatchTimeout: 50 * time.Millisecond,
		MaxSegmentSize:       "10KB",
		DataDirectory:        "testdata",
	}

	wal, err := NewWAL(&cfg, logger)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal.Start(ctx)

	err = wal.Set(context.WithValue(ctx, "ID", uint64(1)), "key1", "value1")
	assert.NoError(t, err)
	err = wal.Set(context.WithValue(ctx, "ID", uint64(2)), "key2", "value2")
	assert.NoError(t, err)
}
