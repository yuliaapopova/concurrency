package engine

import (
	"context"
	"hash/fnv"

	"go.uber.org/zap"
)

type Engine struct {
	log        *zap.Logger
	partitions []*HashTable
}

func NewEngine(log *zap.Logger, partitionCount int) *Engine {
	engine := &Engine{
		log:        log,
		partitions: make([]*HashTable, partitionCount),
	}

	for i := 0; i < partitionCount; i++ {
		engine.partitions[i] = NewHashTable()
	}

	return engine
}

func (e *Engine) Get(ctx context.Context, key string) string {
	hash := e.hash(key)
	partition := e.partitions[hash]
	return partition.Get(key)
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	hash := e.hash(key)
	partition := e.partitions[hash]
	partition.Set(key, value)
}

func (e *Engine) Delete(ctx context.Context, key string) {
	hash := e.hash(key)
	partition := e.partitions[hash]
	partition.Delete(key)
}

func (e *Engine) hash(key string) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return int(hash.Sum32()) % len(e.partitions)
}
