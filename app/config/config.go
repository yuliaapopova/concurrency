package config

import (
	"errors"
	"io"
	"strconv"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

const (
	defaultMessageSize    = "1MB"
	defaultMessageSizeInt = 1024 * 1024
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
		engine.Type = "in_memory"
	}
	return engine
}

func createNetwork(network *Network) *Network {
	if network == nil {
		network = &Network{}
	}

	if network.MaxMessageSize == "" {
		network.MaxMessageSize = "1MB"
	}

	if network.IdleTimeout == 0 {
		network.IdleTimeout = 5 * time.Minute
	}

	if network.MaxConn == 0 {
		network.MaxConn = 100
	}

	if network.Address == "" {
		network.Address = "127.0.0.1:9000"
	}

	return network
}

func createLogging(logging *Logging) *Logging {
	if logging == nil {
		logging = &Logging{}
	}

	if logging.Level == "" {
		logging.Level = "info"
	}

	if logging.Output == "" {
		logging.Output = "/log/output.log"
	}

	return logging
}

func (n *Network) ParseMessageSize() (int, error) {
	if n.MaxMessageSize == "0" || len(n.MaxMessageSize) < 2 || n.MaxMessageSize[0] == '0' {
		return 0, errors.New("invalid max_message_size")
	}

	r := []rune(n.MaxMessageSize)
	sizeR := []rune{}
	idx := 0

	for unicode.IsDigit(r[idx]) {
		sizeR = append(sizeR, r[idx])
		idx++
	}

	size, err := strconv.Atoi(string(sizeR))
	if err != nil {
		return 0, errors.New("invalid max_message_size")
	}
	switch n.MaxMessageSize[idx:] {
	case "B", "b", "":
		break
	case "KB", "kb", "Kb":
		size = size * 1024
	case "MB", "mb", "Mb":
		size = size * 1024 * 1024
	case "GB", "gb", "Gb":
		size = size * 1024 * 1024 * 1024
	default:
		return 0, errors.New("invalid max_message_size")
	}
	return size, nil
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
