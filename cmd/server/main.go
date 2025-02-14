package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"concurrency/app/config"
	"concurrency/app/service"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &config.Config{}

	data, err := os.ReadFile(ConfigFileName)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)
	cfg, err = config.NewConfig(reader)
	if err != nil {
		log.Fatal(err)
	}

	service.Start(ctx, cfg)
}
