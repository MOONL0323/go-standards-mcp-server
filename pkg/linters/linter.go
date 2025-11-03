package linters

import (
	"context"

	"go-standards-mcp-server/pkg/models"
)

// Linter is the interface that all linters must implement
type Linter interface {
	// Name returns the name of the linter
	Name() string

	// Run executes the linter on the specified directory with the given config
	Run(ctx context.Context, workDir, configPath string) ([]models.Issue, error)

	// IsAvailable checks if the linter is available on the system
	IsAvailable() bool
}
