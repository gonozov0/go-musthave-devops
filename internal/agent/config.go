package agent

import (
	"fmt"
	"os"
	"strconv"
)

// Config is a struct that represents configuration
type Config struct {
	CollectInterval int // in seconds
	SendInterval    int // in seconds
	ServerAddress   string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (Config, error) {
	config := Config{}

	collectIntervalStr := os.Getenv("COLLECT_INTERVAL")
	if collectIntervalStr == "" {
		config.CollectInterval = 2 // default value
	} else {
		var err error
		config.CollectInterval, err = strconv.Atoi(collectIntervalStr)
		if err != nil {
			return config, fmt.Errorf("error parsing COLLECT_INTERVAL: %v", err)
		}
	}

	sendIntervalStr := os.Getenv("SEND_INTERVAL")
	if sendIntervalStr == "" {
		config.SendInterval = 10 // default value
	} else {
		var err error
		config.SendInterval, err = strconv.Atoi(sendIntervalStr)
		if err != nil {
			return config, fmt.Errorf("error parsing REPORT_INTERVAL: %v", err)
		}
	}

	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	if config.ServerAddress == "" {
		config.ServerAddress = "http://127.0.0.1:8080" // default value
	}

	return config, nil
}
