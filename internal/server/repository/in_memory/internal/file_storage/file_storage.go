package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileStorage is a struct that represents a file storage.
type FileStorage struct {
	filePath string
}

// NewFileStorage returns a new FileStorage struct.
func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

// SaveMetrics saves metrics to file.
func (fs *FileStorage) SaveMetrics(metrics Metrics) error {
	bytes, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	err = os.WriteFile(fs.filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadMetrics loads metrics from file.
func (fs *FileStorage) LoadMetrics() (Metrics, error) {
	var metrics Metrics
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		return metrics, nil // If the file doesn't exist, just skip the restore step
	}

	bytes, err := os.ReadFile(fs.filePath)
	if err != nil {
		return metrics, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		return metrics, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	return metrics, nil
}
