package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host           string `env:"ADDRESS"`
	PollInterval   int    `env:"REPORT_INTERVAL"`
	ReportInterval int    `env:"POLL_INTERVAL"`
}

func GetConfigs() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "server host")
	flag.IntVar(&config.PollInterval, "p", 2, "frequency of sending metrics to the server")
	flag.IntVar(&config.ReportInterval, "r", 10, "frequency of polling metrics from the runtime package")
	flag.Parse()
	if len(flag.Args()) > 0 {
		panic("used not declared arguments")
	}

	var envConfig Config
	err := env.Parse(&envConfig)
	if err != nil {
		panic(err)
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

	return config
}
