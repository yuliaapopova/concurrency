package compute

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const (
	Set = "SET"
	Get = "GET"
	Del = "DEL"
)

type Compute struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Compute {
	return &Compute{
		log: log,
	}
}

type Query struct {
	Command Command
	Args    []string
}

func parseCommand(command string) Command {
	switch command {
	case Set:
		return SET
	case Get:
		return GET
	case Del:
		return DEL
	default:
		return UNKNOWN
	}
}

func validateCommand(command Command, args []string) error {
	switch command {
	case SET:
		if len(args) != 2 {
			return fmt.Errorf("invalid arguments for command: %s", command.String())
		}
	case GET, DEL:
		if len(args) != 1 {
			return fmt.Errorf("invalid arguments for command: %v", command)
		}
	default:
		return fmt.Errorf("invalid arguments for command: %v", command)
	}
	return nil
}

func (c *Compute) Parse(ctx context.Context, query string) (Query, error) {
	args := strings.Fields(query)
	if len(args) == 0 {
		return Query{}, errors.New("no command specified")
	}
	command := parseCommand(args[0])
	err := validateCommand(command, args[1:])
	if err != nil {
		c.log.Debug("invalid command", zap.String("command", command.String()), zap.Error(err))
		return Query{}, err
	}
	return Query{Command: command, Args: args[1:]}, nil
}
