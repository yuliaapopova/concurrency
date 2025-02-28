package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStorage_Get(t *testing.T) {
	repo := NewEngine(zap.NewNop())
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
	repo := NewEngine(zap.NewNop())
	ctx := context.Background()

	args := []string{"key", "value"}
	repo.Set(ctx, args[0], args[1])
	assert.Equal(t, repo.data["key"], "value")
}

func TestStorage_Del(t *testing.T) {
	repo := NewEngine(zap.NewNop())
	ctx := context.Background()

	argument := "key"
	repo.Set(ctx, argument, "value")
	repo.Delete(ctx, argument)
	assert.Equal(t, repo.data["key"], "")
}
