package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"concurrency/app/common"
	"concurrency/app/config"
	"go.uber.org/zap"
)

type TCPServer struct {
	listener net.Listener
	sem      chan struct{}

	bufferSize int
	maxConn    int

	logger *zap.Logger
}

func NewTCPServer(ctx context.Context, cfg *config.Network, logger *zap.Logger) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("error creating tcp listener: %w", err)
	}

	bufferSize, err := common.ParseMessageSize(cfg.MaxMessageSize)
	if err != nil {
		return nil, fmt.Errorf("error parsing message size: %w", err)
	}

	server := &TCPServer{
		listener:   listener,
		logger:     logger,
		bufferSize: bufferSize,
		maxConn:    cfg.MaxConn,
	}

	if cfg.MaxConn != 0 {
		server.sem = make(chan struct{}, cfg.MaxConn)
	}

	server.logger.Info("created tcp listener", zap.String("address", cfg.Address))

	return server, nil
}

func (s *TCPServer) HandleQueries(ctx context.Context, handler func(context.Context, []byte) []byte) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				s.logger.Error("error accepting connection", zap.Error(err))
				continue
			}
			s.sem <- struct{}{}
			go func() {
				defer func() {
					<-s.sem
				}()
				s.handleConnection(ctx, conn, handler)
			}()
		}
	}()

	<-ctx.Done()
	s.listener.Close()
	wg.Wait()
}

func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn, handler func(context.Context, []byte) []byte) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("recover panic", zap.Any("error", err))
		}

		err := conn.Close()
		if err != nil {
			s.logger.Error("error closing connection", zap.Error(err))
		}
	}()

	request := make([]byte, s.bufferSize)

	for {
		count, err := conn.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Error("error reading request", zap.Error(err))
			break
		}
		if count >= s.bufferSize {
			s.logger.Error("message too large", zap.Int("count", count))
			break
		}

		response := handler(ctx, request[:count])
		_, err = conn.Write(response)
		if err != nil {
			return
		}
	}
}
