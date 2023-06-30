package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddress  string
	DbAddress      string
	AccrualAddress string
	SecretKey      string
}

var User string

func ParseFlagOs() *ServerConfig {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.DbAddress, "d", "postgres://localhost:5432", "server address")
	flag.StringVar(&cfg.SecretKey, "k", "qwerty12345aszx", "key for hashing")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "address accrual")
	flag.Parse()

	if envServer := os.Getenv("RUN_ADDRESS"); envServer != "" {
		cfg.ServerAddress = envServer
	}
	if envDb := os.Getenv("DATABASE_URI"); envDb != "" {
		cfg.DbAddress = envDb
	}
	if secretKey := os.Getenv("HASH_KEY"); secretKey != "" {
		cfg.DbAddress = secretKey
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		cfg.AccrualAddress = envAccrual
	}
	User = os.Getenv("USER")

	return cfg
}
