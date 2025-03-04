package wal

import (
	"bytes"
	"sort"

	"go.uber.org/zap"
)

type Segment interface {
	Write([]byte) error
	LoadData() ([][]byte, error)
}

type LogManager struct {
	segment Segment
	logger  *zap.Logger
}

func NewLogManager(segment Segment, logger *zap.Logger) *LogManager {
	return &LogManager{
		segment: segment,
		logger:  logger,
	}
}

func (l *LogManager) AppendLogs(logs []Log) {
	var buffer bytes.Buffer
	for _, log := range logs {
		if err := log.Encode(&buffer); err != nil {
			l.logger.Error("encode log failed", zap.Error(err))
			l.acknowledgeWrite(logs, err)
			return
		}
	}

	err := l.segment.Write(buffer.Bytes())
	if err != nil {
		l.logger.Error("write logs failed", zap.Error(err))
	}

	l.acknowledgeWrite(logs, err)
}

func (l *LogManager) acknowledgeWrite(logs []Log, err error) {
	for id := range logs {
		logs[id].Status <- err
	}
}

func (l *LogManager) Load() ([]Log, error) {
	data, err := l.segment.LoadData()
	if err != nil {
		return nil, err
	}

	var logs []Log
	for _, logByte := range data {
		logs, err = l.readLogs(logs, logByte)
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].LSN < logs[j].LSN
	})
	l.logger.Debug("load logs", zap.Int("num", len(logs)))

	return logs, nil
}

func (l *LogManager) readLogs(logs []Log, data []byte) ([]Log, error) {
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var log Log
		if err := log.Decode(buffer); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}
