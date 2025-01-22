package storage

import (
	"context"
	"testing"

	"concurrency/app/compute"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStorage_QueryHandler(t *testing.T) {
	repo := NewStorage(zap.NewNop())
	ctx := context.Background()

	testData := map[string]struct {
		query    compute.Query
		expected string
	}{
		"command GET not found": {
			query:    compute.Query{Command: compute.GET, Args: []string{"key"}},
			expected: "not found",
		},
		"command SET": {
			query:    compute.Query{Command: compute.SET, Args: []string{"key", "value"}},
			expected: "ok",
		},
		"command GET": {
			query:    compute.Query{Command: compute.GET, Args: []string{"key"}},
			expected: "value",
		},
		"command DEL": {
			query:    compute.Query{Command: compute.DEL, Args: []string{"key"}},
			expected: "ok",
		},
		"command GET after DEL": {
			query:    compute.Query{Command: compute.GET, Args: []string{"key"}},
			expected: "not found",
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, repo.QueryHandler(ctx, test.query))
		})
	}
}

func TestStorage_Get(t *testing.T) {
	repo := NewStorage(zap.NewNop())
	ctx := context.Background()

	argument := "key"
	res := repo.Get(ctx, argument)
	assert.Equal(t, "not found", res)

	args := []string{"key", "value"}
	_ = repo.Set(ctx, args)
	res = repo.Get(ctx, argument)
	assert.Equal(t, "value", res)
}

func TestStorage_Set(t *testing.T) {
	repo := NewStorage(zap.NewNop())
	ctx := context.Background()

	args := []string{"key", "value"}
	res := repo.Set(ctx, args)
	assert.Equal(t, "ok", res)
}

func TestStorage_Del(t *testing.T) {
	repo := NewStorage(zap.NewNop())
	ctx := context.Background()

	args := []string{"key", "value"}
	_ = repo.Set(ctx, args)
	res := repo.Delete(ctx, "key")
	assert.Equal(t, "ok", res)
}
