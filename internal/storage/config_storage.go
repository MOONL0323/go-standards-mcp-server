package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConfigMetadata stores metadata about a custom configuration
type ConfigMetadata struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConfigStorage manages custom configuration files
type ConfigStorage struct {
	baseDir string
}

// NewConfigStorage creates a new configuration storage
func NewConfigStorage(baseDir string) (*ConfigStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &ConfigStorage{
		baseDir: baseDir,
	}, nil
}

// Save saves a configuration
func (s *ConfigStorage) Save(name, content, description string) error {
	// Validate name
	if name == "" {
		return fmt.Errorf("config name cannot be empty")
	}

	metadata := ConfigMetadata{
		Name:        name,
		Description: description,
		Content:     content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Check if already exists for created time
	if existing, err := s.Get(name); err == nil {
		metadata.CreatedAt = existing.CreatedAt
	}

	// Save metadata
	metaPath := filepath.Join(s.baseDir, name+".json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Save config content
	configPath := filepath.Join(s.baseDir, name+".yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Get retrieves a configuration
func (s *ConfigStorage) Get(name string) (*ConfigMetadata, error) {
	metaPath := filepath.Join(s.baseDir, name+".json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found: %s", name)
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata ConfigMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
}

// GetConfigPath returns the path to the config file
func (s *ConfigStorage) GetConfigPath(name string) (string, error) {
	configPath := filepath.Join(s.baseDir, name+".yaml")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("config not found: %s", name)
		}
		return "", fmt.Errorf("failed to access config: %w", err)
	}
	return configPath, nil
}

// List lists all configurations
func (s *ConfigStorage) List() ([]ConfigMetadata, error) {
	files, err := filepath.Glob(filepath.Join(s.baseDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	var configs []ConfigMetadata
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var metadata ConfigMetadata
		if err := json.Unmarshal(data, &metadata); err != nil {
			continue
		}

		configs = append(configs, metadata)
	}

	return configs, nil
}

// Delete deletes a configuration
func (s *ConfigStorage) Delete(name string) error {
	metaPath := filepath.Join(s.baseDir, name+".json")
	configPath := filepath.Join(s.baseDir, name+".yaml")

	// Delete metadata
	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	// Delete config
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}

