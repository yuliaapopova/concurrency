package config

import (
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultMessageSize = "1MB"
	engineType         = "in_memory"
	defaultIdleTimeout = time.Minute
	defaultAddress     = "127.0.0.1:9000"
	defaultLogLever    = "info"
	defaultLogOutput   = "/log/output.log"
)

type Config struct {
	Engine  *Engine  `yaml:"engine"`
	Network *Network `yaml:"network"`
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
