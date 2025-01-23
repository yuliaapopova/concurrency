package service

import (
	"context"
	"errors"
	"testing"

	"concurrency/app/compute"
	"concurrency/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestStorage_QueryHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	parser := mocks.NewMockCompute(ctrl)
	repo := mocks.NewMockEngine(ctrl)
	logger := zap.NewNop()
	service := New(repo, parser, logger)

	ctx := context.Background()

	testData := map[string]struct {
		queryStr string
		query    compute.Query
		err      error
		expected string
	}{
		"empty queryStr": {
			queryStr: "",
			query:    compute.Query{},
			err:      errors.New("no command specified"),
			expected: "no command specified",
		},
		"unknown command": {
			queryStr: "UPDATE 1",
			query:    compute.Query{},
			err:      errors.New("invalid arguments for command: UNKNOWN"),
			expected: "invalid arguments for command: UNKNOWN",
		},
		"fail argument for command GET": {
			queryStr: "GET key value",
			query:    compute.Query{},
			err:      errors.New("invalid arguments for command: GET"),
			expected: "invalid arguments for command: GET",
		},
		"fail argument for command SET": {
			queryStr: "SET key",
			query:    compute.Query{},
			err:      errors.New("invalid arguments for command: SET"),
			expected: "invalid arguments for command: SET",
		},
		"fail argument for command DEL": {
			queryStr: "DEL key value",
			query:    compute.Query{},
			err:      errors.New("invalid arguments for command: DEL"),
			expected: "invalid arguments for command: DEL",
		},
		"set queryStr": {
			queryStr: "SET key value",
			query:    compute.Query{Command: compute.SET, Args: []string{"key", "value"}},
			err:      nil,
			expected: "",
		},
		"get queryStr": {
			queryStr: "GET key",
			query:    compute.Query{Command: compute.GET, Args: []string{"key"}},
			err:      nil,
			expected: "value",
		},
		"del queryStr": {
			queryStr: "DEL key",
			query:    compute.Query{Command: compute.DEL, Args: []string{"key"}},
			err:      nil,
			expected: "",
		},
	}

	repo.EXPECT().Set(ctx, "key", "value").Return()
	repo.EXPECT().Get(ctx, "key").Return("value")
	repo.EXPECT().Delete(ctx, "key").Return()

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			parser.EXPECT().Parse(ctx, test.queryStr).Return(test.query, test.err)
			res := service.Handler(ctx, test.queryStr)
			assert.Equal(t, test.expected, res)
		})
	}
}
