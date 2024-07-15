package server

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

// Получение конфигурации агента.
func GetConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&config.StoreInterval, "i", 300, "interval in sec to store metrics")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "path to save metrics")
	flag.BoolVar(&config.Restore, "r", true, "retore previous metrics data")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&config.Key, "k", "", "Key for hashing")
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("used not declared arguments")
	}

	var envConfig Config
	err := env.Parse(&envConfig)
	if err != nil {
		log.Fatal(err)
	}

	if envConfig.Host != "" {
		config.Host = envConfig.Host
	}
	if envConfig.StoreInterval != 0 {
		config.StoreInterval = envConfig.StoreInterval
	}
	if envConfig.FileStoragePath != "" {
		config.FileStoragePath = envConfig.FileStoragePath
	}
	if !envConfig.Restore {
		config.Restore = envConfig.Restore
	}
	if envConfig.DatabaseDSN != "" {
		config.DatabaseDSN = envConfig.DatabaseDSN
	}
	if envConfig.Key != "" {
		config.Key = envConfig.Key
	}

	return config
}
