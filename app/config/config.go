package config

import (
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultMessageSize       = "1MB"
	engineType               = "in_memory"
	defaultIdleTimeout       = time.Minute
	defaultAddress           = "127.0.0.1:9000"
	defaultLogLever          = "info"
	defaultLogOutput         = "/log/output.log"
	defaultFlushingBatchSize = 100
	defaultFlushingTimeout   = 10 * time.Second
	defaultMaxSegmentSize    = "10KB"
	defaultDataDirectory     = "/data/spider/wal"
)

type Config struct {
	Engine  *Engine  `yaml:"engine"`
	Network *Network `yaml:"network"`
	WAL     *WAL     `yaml:"wal"`
	Logging *Logging `yaml:"logging"`
}

type Engine struct {
	Type string `yaml:"type"`
}

type Network struct {
	Address        string        `yaml:"address"`
	MaxConn        int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type WAL struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func createEngine(engine *Engine) *Engine {
	if engine == nil {
		engine = &Engine{}
	}

	if engine.Type == "" {
		engine.Type = engineType
	}
	return engine
}

func createWAL(wal *WAL) *WAL {
	if wal == nil {
		return nil
	}

	if wal.FlushingBatchSize == 0 {
		wal.FlushingBatchSize = defaultFlushingBatchSize
	}

	if wal.FlushingBatchTimeout == 0 {
		wal.FlushingBatchTimeout = defaultFlushingTimeout
	}

	if wal.MaxSegmentSize == "" {
		wal.MaxSegmentSize = defaultMaxSegmentSize
	}

	if wal.DataDirectory == "" {
		wal.DataDirectory = defaultDataDirectory
	}

	return wal
}

func createNetwork(network *Network) *Network {
	if network == nil {
		network = &Network{}
	}

	if network.MaxMessageSize == "" {
		network.MaxMessageSize = defaultMessageSize
	}

	if network.IdleTimeout == 0 {
		network.IdleTimeout = defaultIdleTimeout
	}

	if network.Address == "" {
		network.Address = defaultAddress
	}

	return network
}

func createLogging(logging *Logging) *Logging {
	if logging == nil {
		logging = &Logging{}
	}

	if logging.Level == "" {
		logging.Level = defaultLogLever
	}

	if logging.Output == "" {
		logging.Output = defaultLogOutput
	}

	return logging
}

func BuildConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	config.Engine = createEngine(config.Engine)
	config.Network = createNetwork(config.Network)
	config.WAL = createWAL(config.WAL)
	config.Logging = createLogging(config.Logging)

	return config
}

func NewConfig(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("nil reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	cnf := BuildConfig(&config)
	return cnf, nil
}
