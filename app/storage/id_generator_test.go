package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDGenerator(t *testing.T) {
	generator := NewIDGenerator(0)
	assert.Equal(t, uint64(0), generator.ID.Load())
}

func TestNewIDGenerator_WithPrevID(t *testing.T) {
	generator := NewIDGenerator(10)
	assert.Equal(t, uint64(10), generator.ID.Load())
}

func TestIDGenerator_NextID(t *testing.T) {
	generator := NewIDGenerator(10)
	generator.NextID()
	assert.Equal(t, uint64(11), generator.ID.Load())
}
