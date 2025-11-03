package linters

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"go-standards-mcp-server/pkg/models"
	"go.uber.org/zap"
)

// GolangciLint implements the golangci-lint linter
type GolangciLint struct {
	logger *zap.Logger
}

// NewGolangciLint creates a new GolangciLint instance
func NewGolangciLint(logger *zap.Logger) (*GolangciLint, error) {
	g := &GolangciLint{
		logger: logger,
	}

	if !g.IsAvailable() {
		return nil, fmt.Errorf("golangci-lint not found in PATH")
	}

	return g, nil
}

// Name returns the name of the linter
func (g *GolangciLint) Name() string {
	return "golangci-lint"
}

// IsAvailable checks if golangci-lint is available
func (g *GolangciLint) IsAvailable() bool {
	_, err := exec.LookPath("golangci-lint")
	return err == nil
}

// Run executes golangci-lint
func (g *GolangciLint) Run(ctx context.Context, workDir, configPath string) ([]models.Issue, error) {
	args := []string{
		"run",
		"--out-format=json",
		"--print-issued-lines=false",
	}

	if configPath != "" {
		args = append(args, "--config", configPath)
	}

	args = append(args, "./...")

	cmd := exec.CommandContext(ctx, "golangci-lint", args...)
	cmd.Dir = workDir

	g.logger.Debug("Running golangci-lint",
		zap.String("workDir", workDir),
		zap.String("config", configPath),
		zap.Strings("args", args))

	output, err := cmd.CombinedOutput()
	
	// golangci-lint returns non-zero exit code when issues are found
	// We only treat it as error if output parsing fails
	if err != nil && len(output) == 0 {
		return nil, fmt.Errorf("golangci-lint failed: %w", err)
	}

	// Parse JSON output
	var result GolangciLintResult
	if len(output) > 0 {
		if err := json.Unmarshal(output, &result); err != nil {
			g.logger.Warn("Failed to parse golangci-lint output",
				zap.Error(err),
				zap.String("output", string(output)))
			return []models.Issue{}, nil
		}
	}

	// Convert to our Issue format
	issues := make([]models.Issue, 0, len(result.Issues))
	for _, issue := range result.Issues {
		issues = append(issues, models.Issue{
			File:     g.relativePath(workDir, issue.Pos.Filename),
			Line:     issue.Pos.Line,
			Column:   issue.Pos.Column,
			Severity: g.mapSeverity(issue.Severity),
			Category: g.mapCategory(issue.FromLinter),
			Rule:     issue.FromLinter,
			Message:  issue.Text,
			Source:   "golangci-lint",
			Code:     issue.SourceLines,
		})
	}

	g.logger.Debug("golangci-lint completed", zap.Int("issues", len(issues)))
	return issues, nil
}

// relativePath returns a relative path if possible
func (g *GolangciLint) relativePath(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}

// mapSeverity maps golangci-lint severity to our severity levels
func (g *GolangciLint) mapSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "error":
		return "error"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}

// mapCategory maps linter names to categories
func (g *GolangciLint) mapCategory(linter string) string {
	categoryMap := map[string]string{
		"gofmt":       "format",
		"goimports":   "format",
		"gosec":       "security",
		"govet":       "logic",
		"staticcheck": "logic",
		"errcheck":    "error-handling",
		"ineffassign": "performance",
		"unused":      "dead-code",
		"deadcode":    "dead-code",
		"varcheck":    "dead-code",
		"structcheck": "dead-code",
		"gocyclo":     "complexity",
		"gocognit":    "complexity",
		"nestif":      "complexity",
		"dupl":        "duplication",
		"goconst":     "maintainability",
	}

	if category, ok := categoryMap[linter]; ok {
		return category
	}
	return "other"
}

// GolangciLintResult represents golangci-lint JSON output
type GolangciLintResult struct {
	Issues []GolangciLintIssue `json:"Issues"`
}

// GolangciLintIssue represents a single issue from golangci-lint
type GolangciLintIssue struct {
	FromLinter  string              `json:"FromLinter"`
	Text        string              `json:"Text"`
	Severity    string              `json:"Severity"`
	SourceLines string              `json:"SourceLines"`
	Pos         GolangciLintPosition `json:"Pos"`
}

// GolangciLintPosition represents a position in the source code
type GolangciLintPosition struct {
	Filename string `json:"Filename"`
	Line     int    `json:"Line"`
	Column   int    `json:"Column"`
}
