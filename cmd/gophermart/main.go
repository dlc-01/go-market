package main

import (
	"context"
	"log"

	"github.com/dlc/go-market/internal/accrual"
	"github.com/dlc/go-market/internal/app"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/storage"
)

func main() {
	cfg, err := config.ParseFlagOs()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.InitLogger(); err != nil {
		log.Fatal(err)
	}

	if err := storage.InitDBStorage(context.Background(), cfg); err != nil {
		logger.Fatalf("error while  init storage: %s", err)
	}

	go accrual.ListenOtherServ(context.Background(), cfg)

	app.Run(cfg)

	storage.Close(context.Background())
}
