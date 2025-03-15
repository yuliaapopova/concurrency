package replication

import (
	"context"
	"sync"
	"testing"

	"concurrency/app/config"
	"concurrency/app/network"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMaster_Start(t *testing.T) {
	logger, _ := zap.NewProduction()
	walDir := "testdata"
	ctx := context.Background()
	masterAddr := "127.0.0.1:9020"

	server, err := network.NewTCPServer(ctx, &config.Network{Address: masterAddr, MaxMessageSize: "4KB", MaxConn: 2}, logger)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	master := NewMaster(server, walDir, logger)
	assert.NotNil(t, master)

	go master.Start(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		request := NewRequest("wal_0.wal")
		requestByte, err := EncodeRequest(request)
		assert.NoError(t, err)

		client, err := network.NewTcpClient(masterAddr)
		assert.NoError(t, err)

		responseBytes, err := client.Send(requestByte)
		assert.NoError(t, err)

		var response Response
		err = DecodeResponse(&response, responseBytes)
		assert.NoError(t, err)
		assert.Equal(t, response.SegmentName, "wal_1000.wal")
		assert.Equal(t, response.Success, true)
	}()

	wg.Wait()
}
