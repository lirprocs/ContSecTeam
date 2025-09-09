package config

type DefaultConfig struct {
	Workers   int
	QueueSize int
	Port      string
}

var Defaults = DefaultConfig{
	Workers:   4,
	QueueSize: 64,
	Port:      "8080",
}
