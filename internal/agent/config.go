package agent

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

// Config is a struct that represents configuration
type Config struct {
	PollInterval   int // in seconds
	ReportInterval int // in seconds
	ServerAddress  string
}

// newConfig returns a new Config struct with default values
func newConfig() Config {
	return Config{
		PollInterval:   2,
		ReportInterval: 10,
		ServerAddress:  "http://localhost:8080",
	}
}

// LoadConfig loads the configuration from envs and command-line flags
func LoadConfig() (Config, error) {
	config := newConfig()

	if envPollInterval, exists := os.LookupEnv("POLL_INTERVAL"); exists {
		parsed, err := strconv.Atoi(envPollInterval)
		if err == nil {
			config.PollInterval = parsed
		}
	}
	if envReportInterval, exists := os.LookupEnv("REPORT_INTERVAL"); exists {
		parsed, err := strconv.Atoi(envReportInterval)
		if err == nil {
			config.ReportInterval = parsed
		}
	}
	if envAddress, exists := os.LookupEnv("ADDRESS"); exists {
		config.ServerAddress = envAddress
	}

	serverAddr := flag.String("a", config.ServerAddress, "HTTP server endpoint address")
	reportInterval := flag.Int("r", config.ReportInterval, "Frequency of sending metrics to the server (in seconds)")
	pollInterval := flag.Int("p", config.PollInterval, "Frequency of polling metrics from the runtime package (in seconds)")

	flag.Parse()

	if len(flag.Args()) > 0 {
		return config, errors.New("unexpected arguments provided")
	}

	config.ServerAddress = *serverAddr
	config.ReportInterval = *reportInterval
	config.PollInterval = *pollInterval

	return config, nil
}
