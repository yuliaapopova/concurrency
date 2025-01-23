package compute

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestParseGet(t *testing.T) {
	logger := zap.NewNop()
	service := New(logger)

	ctx := context.Background()

	args := map[string]struct {
		queryStr string
		expect   Query
		err      error
	}{
		"empty queryStr": {
			queryStr: "",
			expect:   Query{},
			err:      errors.New("no command specified"),
		},
		"unknown command": {
			queryStr: "UPDATE 1",
			expect:   Query{},
			err:      errors.New("invalid arguments for command: UNKNOWN"),
		},
		"fail argument for command GET": {
			queryStr: "GET key value",
			expect:   Query{},
			err:      errors.New("invalid arguments for command: GET"),
		},
		"fail argument for command SET": {
			queryStr: "SET key",
			expect:   Query{},
			err:      errors.New("invalid arguments for command: SET"),
		},
		"fail argument for command DEL": {
			queryStr: "DEL key value",
			expect:   Query{},
			err:      errors.New("invalid arguments for command: DEL"),
		},
		"set queryStr": {
			queryStr: "SET key value",
			expect:   Query{Command: SET, Args: []string{"key", "value"}},
			err:      nil,
		},
		"get queryStr": {
			queryStr: "GET key",
			expect:   Query{Command: GET, Args: []string{"key"}},
			err:      nil,
		},
		"del queryStr": {
			queryStr: "DEL key",
			expect:   Query{Command: DEL, Args: []string{"key"}},
			err:      nil,
		},
	}

	for name, arg := range args {
		t.Run(name, func(t *testing.T) {
			res, err := service.Parse(ctx, arg.queryStr)
			assert.Equal(t, arg.expect, res)
			assert.Equal(t, arg.err, err)
		})
	}
}
