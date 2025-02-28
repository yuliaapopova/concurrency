package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testCfgData = `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:8080"
  max_connections: 10
  max_message_size: "4KB"
  idle_timeout: 10m
wal:
  flushing_batch_length: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "/data/spider/wal"
logging:
  level: "info"
  output: "/log/output.log"
`

func TestNewConfig(t *testing.T) {
	testData := map[string]struct {
		reader string
		cfg    *Config
	}{
		"empty config": {
			reader: "",
			cfg: &Config{
				Engine: &Engine{
					Type: "in_memory",
				},
				Network: &Network{
					Address:        "127.0.0.1:9000",
					MaxConn:        0,
					MaxMessageSize: "1MB",
					IdleTimeout:    time.Minute * 1,
				},
				Logging: &Logging{
					Level:  "info",
					Output: "/log/output.log",
				},
			},
		},
		"config": {
			reader: testCfgData,
			cfg: &Config{
				Engine: &Engine{
					Type: "in_memory",
				},
				Network: &Network{
					Address:        "127.0.0.1:8080",
					MaxConn:        10,
					MaxMessageSize: "4KB",
					IdleTimeout:    time.Minute * 10,
				},
				WAL: &WAL{
					FlushingBatchSize:    100,
					FlushingBatchTimeout: time.Millisecond * 10,
					MaxSegmentSize:       "10MB",
					DataDirectory:        "/data/spider/wal",
				},
				Logging: &Logging{
					Level:  "info",
					Output: "/log/output.log",
				},
			},
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			reader := strings.NewReader(test.reader)
			cfg, err := NewConfig(reader)
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, cfg)
		})
	}
}
