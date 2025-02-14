package network

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"concurrency/app/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewTCPServer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	//server, err := NewTCPServer(ctx, logger)

	testData := map[string]struct {
		cfg      *config.Network
		logger   *zap.Logger
		err      error
		response func() *TCPServer
	}{
		"without logger": {
			cfg: &config.Network{
				Address: "localhost:8082",
			},
			err: errors.New("logger is nil"),
		},
		"invalid buffer size": {
			cfg: &config.Network{
				Address:        "localhost:8083",
				MaxMessageSize: "2KB2",
			},
			logger: logger,
			err:    errors.New("error parsing message size: invalid max_message_size"),
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			_, err := NewTCPServer(ctx, test.cfg, test.logger)
			if err != nil {
				assert.Error(t, err, test.err)
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}
}

func TestTCPServer_HandleQueries(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	cfg := &config.Network{
		Address:        "localhost:8084",
		MaxMessageSize: "4KB",
	}
	server, err := NewTCPServer(ctx, cfg, logger)
	assert.NoError(t, err)
	defer server.listener.Close()

	handler := func(ctx context.Context, bytes []byte) []byte {
		return bytes
	}

	go func() {
		server.HandleQueries(ctx, handler)
	}()

	time.Sleep(1 * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", cfg.Address)
		defer conn.Close()
		assert.NoError(t, err)

		_, err = conn.Write([]byte("msg"))
		assert.NoError(t, err)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		assert.NoError(t, err)

		err = conn.Close()
		assert.NoError(t, err)
		assert.Equal(t, []byte("msg"), buf[:n])
	}()
}
