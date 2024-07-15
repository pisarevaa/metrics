package agent

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host           string `env:"ADDRESS"`
	PollInterval   int    `env:"REPORT_INTERVAL"`
	ReportInterval int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// Получение конфигурации агента.
func GetConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "server host")
	flag.IntVar(&config.PollInterval, "p", 2, "frequency of sending metrics to the server")
	flag.IntVar(&config.ReportInterval, "r", 10, "frequency of polling metrics from the runtime package")
	flag.StringVar(&config.Key, "k", "", "Key for hashing")
	flag.IntVar(&config.RateLimit, "l", 20, "Rate limit to send HTTP requests")

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

	return config
}
