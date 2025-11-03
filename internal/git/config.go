package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IncrementalConfig represents git incremental check configuration
type IncrementalConfig struct {
	Enabled      bool   `json:"enabled"`       // Enable incremental checking
	AutoCommit   bool   `json:"auto_commit"`   // Auto check on commit
	AutoPush     bool   `json:"auto_push"`     // Auto check on push
	BaseBranch   string `json:"base_branch"`   // Base branch for comparison (default: origin/main)
	ConfigFile   string `json:"config_file"`   // Path to golangci-lint config
	FailOnError  bool   `json:"fail_on_error"` // Fail git operation on error
	HooksInstalled bool `json:"hooks_installed"` // Whether git hooks are installed
}

// DefaultIncrementalConfig returns default configuration
func DefaultIncrementalConfig() *IncrementalConfig {
	return &IncrementalConfig{
		Enabled:      false,
		AutoCommit:   false,
		AutoPush:     false,
		BaseBranch:   "origin/main",
		FailOnError:  true,
		HooksInstalled: false,
	}
}

// ConfigManager manages incremental check configuration
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new config manager
func NewConfigManager(repoPath string) *ConfigManager {
	configPath := filepath.Join(repoPath, ".go-standards.json")
	return &ConfigManager{
		configPath: configPath,
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() (*IncrementalConfig, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultIncrementalConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	var config IncrementalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return &config, nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save(config *IncrementalConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	
	return nil
}

// Enable enables incremental checking
func (cm *ConfigManager) Enable() error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.Enabled = true
	return cm.Save(config)
}

// Disable disables incremental checking
func (cm *ConfigManager) Disable() error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.Enabled = false
	return cm.Save(config)
}

// EnableAutoCommit enables auto-check on commit
func (cm *ConfigManager) EnableAutoCommit() error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.AutoCommit = true
	return cm.Save(config)
}

// EnableAutoPush enables auto-check on push
func (cm *ConfigManager) EnableAutoPush() error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.AutoPush = true
	return cm.Save(config)
}

// SetBaseBranch sets the base branch for comparison
func (cm *ConfigManager) SetBaseBranch(branch string) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.BaseBranch = branch
	return cm.Save(config)
}

// SetConfigFile sets the golangci-lint config file path
func (cm *ConfigManager) SetConfigFile(path string) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.ConfigFile = path
	return cm.Save(config)
}

// MarkHooksInstalled marks git hooks as installed
func (cm *ConfigManager) MarkHooksInstalled(installed bool) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.HooksInstalled = installed
	return cm.Save(config)
}
