package agent

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/pisarevaa/metrics/internal/agent/utils"
)

type Config struct {
	Host           string `env:"ADDRESS"         json:"address"`
	PollInterval   int    `env:"REPORT_INTERVAL" json:"poll_interval"`
	ReportInterval int    `env:"POLL_INTERVAL"   json:"report_interval"`
	Key            string `env:"KEY"             json:"key,omitempty"`
	RateLimit      int    `env:"RATE_LIMIT"      json:"rate_limit,omitempty"`
	CryptoKey      string `env:"CRYPTO_KEY"      json:"crypto_key"`
	Config         string `env:"CONFIG"          json:"config,omitempty"`
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
	if config.PollInterval == 0 && fileConfig.PollInterval != 0 {
		config.PollInterval = fileConfig.PollInterval
	}
	if config.ReportInterval == 0 && fileConfig.ReportInterval != 0 {
		config.ReportInterval = fileConfig.ReportInterval
	}
	if config.CryptoKey == "" && fileConfig.CryptoKey != "" {
		config.CryptoKey = fileConfig.CryptoKey
	}
	return nil
}

// Получение конфигурации агента.
func GetConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "server host")
	flag.IntVar(&config.PollInterval, "p", 2, "frequency of sending metrics to the server")
	flag.IntVar(&config.ReportInterval, "r", 10, "frequency of polling metrics from the runtime package")
	flag.StringVar(&config.Key, "k", "", "Key for hashing")
	flag.IntVar(&config.RateLimit, "l", 20, "Rate limit to send HTTP requests")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "path to public key")
	flag.StringVar(&config.Config, "c", "agent_env.json", "path to config JSON file")
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
	if envConfig.PollInterval != 0 {
		config.PollInterval = envConfig.PollInterval
	}
	if envConfig.ReportInterval != 0 {
		config.ReportInterval = envConfig.ReportInterval
	}
	if envConfig.Key != "" {
		config.Key = envConfig.Key
	}
	if envConfig.RateLimit != 0 {
		config.RateLimit = envConfig.RateLimit
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
		err = utils.InitPublicKey(config.CryptoKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	return config
}
