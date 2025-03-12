package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStorage_Get(t *testing.T) {
	repo := NewEngine(zap.NewNop(), 8)
	ctx := context.Background()

	argument := "key"
	res := repo.Get(ctx, argument)
	assert.Equal(t, "not found", res)

	args := []string{"key", "value"}
	repo.Set(ctx, args[0], args[1])
	res = repo.Get(ctx, argument)
	assert.Equal(t, "value", res)
}

func TestStorage_Set(t *testing.T) {
	repo := NewEngine(zap.NewNop(), 8)
	ctx := context.Background()

	args := []string{"key", "value"}
	repo.Set(ctx, args[0], args[1])
	hash := repo.hash(args[0])
	assert.Equal(t, repo.partitions[hash].data[args[0]], "value")
}

func TestStorage_Del(t *testing.T) {
	repo := NewEngine(zap.NewNop(), 8)
	ctx := context.Background()

	argument := "key"
	repo.Set(ctx, argument, "value")
	repo.Delete(ctx, argument)
	hash := repo.hash(argument)
	assert.Equal(t, repo.partitions[hash].data["key"], "")
}
