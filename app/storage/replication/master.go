package replication

import (
	"context"
	"fmt"
	"os"

	"concurrency/app/storage/filesystem"
	"go.uber.org/zap"
)

type TCPServer interface {
	HandleQueries(ctx context.Context, handler func(context.Context, []byte) []byte)
}

type Master struct {
	server  TCPServer
	walDir  string
	logger  *zap.Logger
	segment *filesystem.Segment
}

func NewMaster(server TCPServer, walDir string, logger *zap.Logger) *Master {
	return &Master{
		server: server,
		walDir: walDir,
		logger: logger,
	}
}

func (m *Master) IsMaster() bool {
	return true
}

func (m *Master) Start(ctx context.Context) {
	m.logger.Debug("starting master replication")
	m.server.HandleQueries(ctx, func(ctx context.Context, data []byte) []byte {
		if ctx.Err() != nil {
			return nil
		}

		var request Request
		if err := DecodeRequest(ctx, &request, data); err != nil {
			m.logger.Error("failed to decode request", zap.Error(err))
			return nil
		}

		response := m.sync(request)
		responseData, err := EncodeResponse(ctx, &response)
		if err != nil {
			m.logger.Error("failed to encode response", zap.Error(err))
			return nil
		}

		return responseData
	})
}

func (m *Master) sync(request Request) Response {
	var response Response

	segmentNameNext, err := filesystem.SegmentNameNext(m.walDir, request.LastSegmentName)
	if err != nil {
		m.logger.Error("failed to sync segment name", zap.Error(err))
		return response
	}

	if segmentNameNext == "" {
		response.Success = true
		return response
	}

	path, err := filesystem.Path(m.walDir)
	if err != nil {
		m.logger.Error("failed to sync segment path", zap.Error(err))
		return response
	}
	filename := fmt.Sprintf("%s/%s", path, segmentNameNext)
	data, err := os.ReadFile(filename)
	if err != nil {
		m.logger.Error("failed to sync segment", zap.Error(err))
		return response
	}

	response.Success = true
	response.SegmentName = segmentNameNext
	response.SegmentData = data
	return response
}
