package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ServerConfig struct {
	ServerAddress  string
	DBAddress      string
	AccrualAddress string
	SecretKey      string
	Poll           int
}

var User string

func ParseFlagOs() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.DBAddress, "d", "postgres://localhost:5432", "server address")
	flag.StringVar(&cfg.SecretKey, "k", "qwerty12345aszx", "key for hashing")
	flag.StringVar(&cfg.AccrualAddress, "r", "http://localhost:8081", "address accrual")
	flag.IntVar(&cfg.Poll, "p", 2, "Poll interval")
	flag.Parse()

	if envServer := os.Getenv("RUN_ADDRESS"); envServer != "" {
		cfg.ServerAddress = envServer
	}
	if envDB := os.Getenv("DATABASE_URI"); envDB != "" {
		cfg.DBAddress = envDB
	}
	if secretKey := os.Getenv("HASH_KEY"); secretKey != "" {
		cfg.DBAddress = secretKey
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		cfg.AccrualAddress = envAccrual
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		if intPoll, err := strconv.ParseInt(envPoll, 10, 32); err == nil {
			cfg.Poll = int(intPoll)
		} else {
			return nil, fmt.Errorf("cannot parse POLL_INTERVAL: %w", err)
		}
	}
	User = os.Getenv("USER")

	return cfg, nil
}
