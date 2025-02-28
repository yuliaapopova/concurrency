package storage

import (
	"math"
	"sync/atomic"
)

type IDGenerator struct {
	ID atomic.Uint64
}

func NewIDGenerator(prevID uint64) *IDGenerator {
	generator := &IDGenerator{}
	generator.ID.Store(prevID)
	return generator
}

func (g *IDGenerator) NextID() uint64 {
	g.ID.CompareAndSwap(math.MaxUint64, 0)
	return g.ID.Add(1)
}
