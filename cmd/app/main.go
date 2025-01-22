package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"concurrency/app/compute"
	"concurrency/app/storage"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	engine := storage.NewStorage(logger)
	queryParser := compute.New(logger, engine)
	ctx := context.Background()
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
		res := queryParser.Parse(ctx, query)
		logger.Info(res)
	}
}
