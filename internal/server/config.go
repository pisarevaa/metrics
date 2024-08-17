package server

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/pisarevaa/metrics/internal/server/utils"
)

type Config struct {
	Host            string `env:"ADDRESS"           json:"address"`
	StoreInterval   int    `env:"STORE_INTERVAL"    json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE"           json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN"      json:"database_dsn"`
	Key             string `env:"KEY"               json:"key,omitempty"`
	CryptoKey       string `env:"CRYPTO_KEY"        json:"crypto_key"`
	Config          string `env:"CONFIG"            json:"config,omitempty"`
}

func getFromJSONFile(config *Config) error {
	var fileConfig Config
	data, err := os.ReadFile(config.Config)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &fileConfig); err != nil {
		return err
	}

	if config.Host == "" && fileConfig.Host != "" {
		config.Host = fileConfig.Host
	}
	if !config.Restore && !fileConfig.Restore {
		config.Restore = fileConfig.Restore
	}
	if config.StoreInterval == 0 && fileConfig.StoreInterval != 0 {
		config.StoreInterval = fileConfig.StoreInterval
	}
	if config.FileStoragePath == "" && fileConfig.FileStoragePath != "" {
		config.FileStoragePath = fileConfig.FileStoragePath
	}
	if config.DatabaseDSN == "" && fileConfig.DatabaseDSN != "" {
		config.DatabaseDSN = fileConfig.DatabaseDSN
	}
	if config.CryptoKey == "" && fileConfig.CryptoKey != "" {
		config.CryptoKey = fileConfig.CryptoKey
	}
	return nil
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
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "path to private key")
	flag.StringVar(&config.Config, "c", "server_env.json", "path to config JSON file")
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("used not declared arguments")
	}

	var envConfig Config
	err := env.Parse(&envConfig)
	if err != nil {
		log.Fatal(err)
	}

	if envConfig.Config != "" {
		config.Config = envConfig.Config
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
	if envConfig.CryptoKey != "" {
		config.CryptoKey = envConfig.CryptoKey
	}

	if config.Config != "" {
		err = getFromJSONFile(&config)
		if err != nil {
			log.Fatal(err)
		}
	}

	if config.CryptoKey != "" {
		err = utils.InitPrivateKey(config.CryptoKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	return config
}
