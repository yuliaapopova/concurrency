package wal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAppendLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewNop()
	segment := NewMockSegment(ctrl)
	lm := NewLogManager(segment, logger)

	logs := []Log{{LSN: 1, Command: 1, Args: []string{"1", "2"}, Status: make(chan error, 1)}, {LSN: 2, Command: 1, Args: []string{"3", "4"}, Status: make(chan error, 1)}}

	var buffer bytes.Buffer
	for _, l := range logs {
		err := l.Encode(&buffer)
		require.Nil(t, err)
	}

	segment.EXPECT().Write(buffer.Bytes()).Return(nil)

	lm.AppendLogs(logs)

	for _, l := range logs {
		err := <-l.Status
		require.Nil(t, err)
	}
}

func TestLoadData(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := zap.NewNop()
	segment := NewMockSegment(ctrl)
	lm := NewLogManager(segment, logger)

	logs := []Log{{LSN: 1, Command: 1, Args: []string{"1", "2"}, Status: make(chan error, 1)}, {LSN: 2, Command: 1, Args: []string{"3", "4"}, Status: make(chan error, 1)}}

	var buffer bytes.Buffer
	for _, l := range logs {
		err := l.Encode(&buffer)
		require.Nil(t, err)
	}

	segment.EXPECT().LoadData().Return([][]byte{buffer.Bytes()}, nil)

	res, err := lm.Load()
	require.Nil(t, err)

	for i := range res {
		assert.Equal(t, res[i].LSN, logs[i].LSN)
		assert.Equal(t, res[i].Command, logs[i].Command)
		assert.Equal(t, res[i].Args, logs[i].Args)
	}
}
