package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	QueueSize int
	Port      string
}

func New() Config {
	cfg := Config{
		Workers:   Defaults.Workers,
		QueueSize: Defaults.QueueSize,
		Port:      Defaults.Port,
	}
	if w, err := strconv.Atoi(os.Getenv("WORKERS")); err == nil {
		cfg.Workers = w
	}
	if qs, err := strconv.Atoi(os.Getenv("QUEUE_SIZE")); err == nil {
		cfg.QueueSize = qs
	}
	if p := os.Getenv("PORT"); p != "" {
		cfg.Port = p
	}
	return cfg
}
