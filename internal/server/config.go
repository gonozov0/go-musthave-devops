package server

import (
	"errors"
	"flag"
	"os"
)

// Config is a struct that represents configuration
type Config struct {
	ServerAddress string
}

// newConfig returns a new Config struct with default values
func newConfig() Config {
	return Config{
		ServerAddress: "localhost:8080",
	}
}

// LoadConfig loads the configuration from envs and command-line flags
func LoadConfig() (Config, error) {
	config := newConfig()

	if envAddress, exists := os.LookupEnv("ADDRESS"); exists {
		config.ServerAddress = envAddress
	}

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "HTTP server endpoint address")

	flag.Parse()

	if len(flag.Args()) > 0 {
		return config, errors.New("unexpected arguments provided")
	}

	return config, nil
}
