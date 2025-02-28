package wal

import (
	"bytes"
	"encoding/gob"
)

type Log struct {
	LSN     uint64
	Command int
	Args    []string

	Status chan error
}

func NewLog(lsn uint64, command int, args []string) Log {
	return Log{
		LSN:     lsn,
		Command: command,
		Args:    args,

		Status: make(chan error, 1),
	}
}

func (l *Log) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*l)
}

func (l *Log) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(l)
}
