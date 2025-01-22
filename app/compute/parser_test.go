package compute

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestParseGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	engine := NewMockEngine(ctrl)
	logger := zap.NewNop()
	service := New(logger, engine)

	ctx := context.Background()

	args := map[string]struct {
		queryStr string
		expect   string
		query    Query
	}{
		"empty queryStr": {
			queryStr: "",
			expect:   "no command specified",
		},
		"unknown command": {
			queryStr: "UPDATE 1",
			expect:   "invalid arguments for command: UNKNOWN",
		},
		"fail argument for command GET": {
			queryStr: "GET key value",
			expect:   "invalid arguments for command: GET",
		},
		"fail argument for command SET": {
			queryStr: "SET key",
			expect:   "invalid arguments for command: SET",
		},
		"fail argument for command DEL": {
			queryStr: "DEL key value",
			expect:   "invalid arguments for command: DEL",
		},
		"set queryStr": {
			queryStr: "SET key value",
			query:    Query{Command: SET, Args: []string{"key", "value"}},
			expect:   "ok",
		},
		"get queryStr": {
			queryStr: "GET key",
			query:    Query{Command: GET, Args: []string{"key"}},
			expect:   "ok",
		},
		"del queryStr": {
			queryStr: "DEL key",
			query:    Query{Command: DEL, Args: []string{"key"}},
			expect:   "ok",
		},
	}

	for name, arg := range args {
		t.Run(name, func(t *testing.T) {
			if len(arg.query.Args) != 0 {
				engine.EXPECT().QueryHandler(ctx, arg.query).Return(arg.expect)
			}
			res := service.Parse(ctx, arg.queryStr)
			assert.Equal(t, arg.expect, res)
		})
	}
}
