package replication

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"concurrency/app/storage/filesystem"
	"concurrency/app/storage/wal"
	"go.uber.org/zap"
)

type TCPClient interface {
	Send(request []byte) ([]byte, error)
	Close()
}

type Slave struct {
	client          TCPClient
	stream          chan []wal.Log
	syncInterval    time.Duration
	walDir          string
	lastSegmentName string
	logger          *zap.Logger
}

func NewSlave(client TCPClient, walDir string, syncInterval time.Duration, logger *zap.Logger) (*Slave, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	if client == nil {
		return nil, errors.New("client required")
	}

	lastSegmentName, err := filesystem.SegmentLastName(walDir)
	if err != nil {
		return nil, err
	}

	return &Slave{
		client:          client,
		stream:          make(chan []wal.Log),
		syncInterval:    syncInterval,
		walDir:          walDir,
		lastSegmentName: lastSegmentName,
		logger:          logger,
	}, nil
}

func (s *Slave) Start(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer func() {
		ticker.Stop()
		s.client.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Warn("context canceled")
				return
			default:
			}

			select {
			case <-ctx.Done():
				s.logger.Warn("context canceled")
				return
			case <-ticker.C:
				s.sync()
			}
		}
	}()
}

func (s *Slave) isMaster() bool {
	return false
}

func (s *Slave) ReplicationStream() chan []wal.Log {
	return s.stream
}

func (s *Slave) sync() {
	request := NewRequest(s.lastSegmentName)
	data, err := EncodeRequest(request)
	if err != nil {
		s.logger.Error("encode request", zap.Error(err))
		return
	}

	responseByte, err := s.client.Send(data)
	if err != nil {
		s.logger.Error("send request", zap.Error(err))
		return
	}

	var resp Response
	if err = DecodeResponse(&resp, responseByte); err != nil {
		s.logger.Error("decode response", zap.Error(err))
		return
	}

	err = s.saveSegment(resp.SegmentName, resp.SegmentData)
	if err != nil {
		s.logger.Error("save segment", zap.Error(err))
		return
	}

	err = s.writeData(resp.SegmentData)
	if err != nil {
		s.logger.Error("write data", zap.Error(err))
		return
	}
}

func (s *Slave) saveSegment(segmentName string, segmentData []byte) error {
	file, err := filesystem.CreateFile(s.walDir, segmentName)
	if err != nil {
		return err
	}

	return filesystem.WriteFile(file, segmentData)
}

func (s *Slave) writeData(data []byte) error {
	var logs []wal.Log
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var log wal.Log
		if err := log.Decode(buffer); err != nil {
			return fmt.Errorf("decode log: %w", err)
		}

		logs = append(logs, log)
	}
	s.stream <- logs
	return nil
}
