package main

import (
	"context"
	"github.com/dlc/go-market/internal/app"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/storage"
	"log"
)

func main() {
	cfg := config.ParseFlagOs()
	if err := logger.InitLogger(); err != nil {
		log.Fatal(err)
	}
	err := storage.InitDBStorage(context.Background(), cfg)
	if err != nil {
		logger.Fatalf("error while  init storage: %s", err)
	}
	app.Run(cfg)
}
