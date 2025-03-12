package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable_Get(t *testing.T) {
	table := NewHashTable()

	argument := "key"
	res := table.Get(argument)
	assert.Equal(t, "not found", res)

	args := []string{"key", "value"}
	table.Set(args[0], args[1])
	res = table.Get(argument)
	assert.Equal(t, "value", res)
}

func TestHashTable_Set(t *testing.T) {
	repo := NewHashTable()

	args := []string{"key", "value"}
	repo.Set(args[0], args[1])
	assert.Equal(t, repo.data["key"], "value")
}

func TestHashTable_Del(t *testing.T) {
	repo := NewHashTable()

	argument := "key"
	repo.Set(argument, "value")
	repo.Delete(argument)
	assert.Equal(t, repo.data["key"], "")
}
