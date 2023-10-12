package server

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config is a struct that represents configuration
type Config struct {
	ServerAddress   string
	StoreInterval   uint64 // in seconds
	FileStoragePath string
	RestoreFlag     bool
	DatabaseDSN     string
}

// newConfig returns a new Config struct with default values
func newConfig() Config {
	return Config{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		RestoreFlag:     true,
		DatabaseDSN:     "", // "postgres://postgres:postgres@localhost:5442/development?sslmode=disable",
	}
}

// LoadConfig loads the configuration from envs and command-line flags
func LoadConfig() (Config, error) {
	config := newConfig()

	if envAddress, exists := os.LookupEnv("ADDRESS"); exists {
		config.ServerAddress = envAddress
	}
	if envInterval, exists := os.LookupEnv("STORE_INTERVAL"); exists {
		uintEnvInterval, err := strconv.ParseUint(envInterval, 10, 64)
		if err != nil {
			return config, fmt.Errorf("failed to parse STORE_INTERVAL: %w", err)
		}
		config.StoreInterval = uintEnvInterval
	}
	if envFileStoragePath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		config.FileStoragePath = envFileStoragePath
	}
	if envRestoreFlag, exists := os.LookupEnv("RESTORE"); exists {
		boolEnvRestoreFlag, err := strconv.ParseBool(envRestoreFlag)
		if err != nil {
			return config, fmt.Errorf("failed to parse RESTORE_FLAG: %w", err)
		}
		config.RestoreFlag = boolEnvRestoreFlag
	}
	if envDatabaseDSN, exists := os.LookupEnv("DATABASE_DSN"); exists {
		config.DatabaseDSN = envDatabaseDSN
	}

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "HTTP server endpoint address")
	flag.Uint64Var(&config.StoreInterval, "i", config.StoreInterval, "Metrics store interval in seconds")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File storage path")
	flag.BoolVar(&config.RestoreFlag, "r", config.RestoreFlag, "Restore metrics from file storage")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database server address (postgres)")

	flag.Parse()

	if len(flag.Args()) > 0 {
		return config, errors.New("unexpected arguments provided")
	}

	return config, nil
}
