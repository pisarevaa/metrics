package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pisarevaa/metrics/internal/server"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

const readTimeout = 5
const writeTimout = 10

func main() {
	config := server.GetConfig()
	log.Printf("Server is running on %v", config.Host)
	srv := &http.Server{
		Addr:         config.Host,
		Handler:      server.MetricsRouter(config),
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimout * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
