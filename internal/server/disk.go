package server

import (
	"encoding/json"
	"os"

	"github.com/pisarevaa/metrics/internal/server/storage"
)

func SaveToDosk(metrics []storage.Metrics, filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&metrics)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func LoadFromDosk(filename string) ([]storage.Metrics, error) {
	if filename == "" {
		return []storage.Metrics{}, nil
	}
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return []storage.Metrics{}, err
	}
	var metrics []storage.Metrics
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metrics)
	if err != nil {
		return []storage.Metrics{}, err
	}
	err = file.Close()
	if err != nil {
		return []storage.Metrics{}, err
	}
	return metrics, nil
}
