package server

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

func GetConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "address and port to run server")
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

	return config
}
