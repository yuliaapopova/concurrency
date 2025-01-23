package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"concurrency/app/compute"
	"concurrency/app/service"
	"concurrency/app/storage"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	engine := storage.NewStorage(logger)
	queryParser := compute.New(logger)
	ctx := context.Background()
	s := service.New(engine, queryParser, logger)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "exit" {
			logger.Debug("Exiting...")
			break
		}
		res := s.Handler(ctx, query)
		logger.Info(res)
	}
}
