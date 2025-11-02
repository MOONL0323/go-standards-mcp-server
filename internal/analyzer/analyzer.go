package analyzer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MOONL0323/go-standards-mcp-server/internal/config"
	"github.com/MOONL0323/go-standards-mcp-server/pkg/linters"
	"github.com/MOONL0323/go-standards-mcp-server/pkg/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Analyzer handles code analysis operations
type Analyzer struct {
	config  *config.Config
	logger  *zap.Logger
	linters map[string]linters.Linter
}

// NewAnalyzer creates a new Analyzer instance
func NewAnalyzer(cfg *config.Config, logger *zap.Logger) (*Analyzer, error) {
	a := &Analyzer{
		config:  cfg,
		logger:  logger,
		linters: make(map[string]linters.Linter),
	}

	// Initialize linters
	if err := a.initLinters(); err != nil {
		return nil, fmt.Errorf("failed to initialize linters: %w", err)
	}

	return a, nil
}

// initLinters initializes all configured linters
func (a *Analyzer) initLinters() error {
	// Initialize golangci-lint
	if a.config.Linters.GolangciLint.Enabled {
		golangci, err := linters.NewGolangciLint(a.logger)
		if err != nil {
			a.logger.Warn("Failed to initialize golangci-lint", zap.Error(err))
		} else {
			a.linters["golangci-lint"] = golangci
			a.logger.Info("Initialized golangci-lint")
		}
	}

	// Initialize other linters
	if a.config.Linters.Govet.Enabled {
		govet := linters.NewGoVet(a.logger)
		a.linters["govet"] = govet
		a.logger.Info("Initialized govet")
	}

	if len(a.linters) == 0 {
		return fmt.Errorf("no linters available")
	}

	return nil
}

// Analyze performs code analysis based on the request
func (a *Analyzer) Analyze(ctx context.Context, req *models.AnalysisRequest) (*models.AnalysisResult, error) {
	startTime := time.Now()
	analysisID := uuid.New().String()

	a.logger.Info("Starting analysis",
		zap.String("id", analysisID),
		zap.String("standard", req.Standard))

	// Prepare working directory
	workDir, cleanup, err := a.prepareWorkDir(req)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare work directory: %w", err)
	}
	defer cleanup()

	// Load configuration
	configPath, err := a.loadConfig(req.Standard, req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Run analysis
	issues, err := a.runLinters(ctx, workDir, configPath)
	if err != nil {
		a.logger.Error("Analysis failed", zap.Error(err))
		return &models.AnalysisResult{
			ID:        analysisID,
			Status:    "error",
			Issues:    []models.Issue{},
			Summary:   models.Summary{},
			CreatedAt: time.Now(),
		}, err
	}

	// Calculate summary
	summary := a.calculateSummary(issues, workDir, time.Since(startTime))

	// Generate suggestions
	suggestions := a.generateSuggestions(issues)

	result := &models.AnalysisResult{
		ID:          analysisID,
		Status:      "success",
		Issues:      issues,
		Summary:     summary,
		Suggestions: suggestions,
		Metadata: models.Metadata{
			Standard:      req.Standard,
			ToolsUsed:     a.getToolNames(),
			ServerVersion: "1.0.0",
		},
		CreatedAt: time.Now(),
	}

	a.logger.Info("Analysis completed",
		zap.String("id", analysisID),
		zap.Int("issues", len(issues)),
		zap.Float64("score", summary.Score),
		zap.Duration("duration", summary.Duration))

	return result, nil
}

// prepareWorkDir prepares the working directory for analysis
func (a *Analyzer) prepareWorkDir(req *models.AnalysisRequest) (string, func(), error) {
	// If analyzing a project directory, use it directly
	if req.ProjectDir != "" {
		return req.ProjectDir, func() {}, nil
	}

	// If analyzing a file, use its directory
	if req.FilePath != "" {
		return filepath.Dir(req.FilePath), func() {}, nil
	}

	// If analyzing code snippet, create a temporary file
	if req.Code != "" {
		tempDir := filepath.Join(a.config.Analyzer.TempDir, uuid.New().String())
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
		}

		tempFile := filepath.Join(tempDir, "main.go")
		if err := os.WriteFile(tempFile, []byte(req.Code), 0644); err != nil {
			os.RemoveAll(tempDir)
			return "", nil, fmt.Errorf("failed to write temp file: %w", err)
		}

		cleanup := func() {
			os.RemoveAll(tempDir)
		}

		return tempDir, cleanup, nil
	}

	return "", nil, fmt.Errorf("no code, file, or directory specified")
}

// loadConfig loads the appropriate configuration
func (a *Analyzer) loadConfig(standard, customConfig string) (string, error) {
	if standard == "custom" && customConfig != "" {
		// Save custom config to temp file
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(customConfig)))
		configPath := filepath.Join(a.config.Analyzer.TempDir, fmt.Sprintf("config-%s.yaml", hash[:8]))
		
		if err := os.WriteFile(configPath, []byte(customConfig), 0644); err != nil {
			return "", fmt.Errorf("failed to write custom config: %w", err)
		}
		
		return configPath, nil
	}

	// Use predefined template
	// Try multiple paths to find the config file
	possiblePaths := []string{
		filepath.Join("configs", "templates", fmt.Sprintf("%s.yaml", standard)),
		filepath.Join("..", "configs", "templates", fmt.Sprintf("%s.yaml", standard)),
		filepath.Join(getExecutableDir(), "..", "configs", "templates", fmt.Sprintf("%s.yaml", standard)),
	}

	for _, templatePath := range possiblePaths {
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}
	}

	return "", fmt.Errorf("template not found: %s (tried: configs/templates/%s.yaml)", standard, standard)
}

// getExecutableDir returns the directory of the executable
func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(ex)
}

// runLinters runs all configured linters
func (a *Analyzer) runLinters(ctx context.Context, workDir, configPath string) ([]models.Issue, error) {
	var allIssues []models.Issue

	for name, linter := range a.linters {
		a.logger.Debug("Running linter", zap.String("linter", name))

		issues, err := linter.Run(ctx, workDir, configPath)
		if err != nil {
			a.logger.Warn("Linter failed", zap.String("linter", name), zap.Error(err))
			continue
		}

		allIssues = append(allIssues, issues...)
		a.logger.Debug("Linter completed", zap.String("linter", name), zap.Int("issues", len(issues)))
	}

	return allIssues, nil
}

// calculateSummary calculates analysis summary statistics
func (a *Analyzer) calculateSummary(issues []models.Issue, workDir string, duration time.Duration) models.Summary {
	summary := models.Summary{
		TotalIssues:    len(issues),
		Duration:       duration,
		CategoryCounts: make(map[string]int),
	}

	// Count issues by severity
	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			summary.ErrorCount++
		case "warning":
			summary.WarningCount++
		case "info":
			summary.InfoCount++
		}

		// Count by category
		summary.CategoryCounts[issue.Category]++
	}

	// Count files
	files, err := a.countGoFiles(workDir)
	if err == nil {
		summary.FilesAnalyzed = files
	}

	// Calculate quality score (0-100)
	// Simple scoring: start at 100, deduct points for issues
	score := 100.0
	score -= float64(summary.ErrorCount) * 5.0
	score -= float64(summary.WarningCount) * 2.0
	score -= float64(summary.InfoCount) * 0.5

	if score < 0 {
		score = 0
	}

	summary.Score = score

	return summary
}

// countGoFiles counts the number of Go files in a directory
func (a *Analyzer) countGoFiles(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Skip test files and vendor
			if filepath.Base(path) != "main.go" && !filepath.HasPrefix(path, filepath.Join(dir, "vendor")) {
				count++
			}
		}
		return nil
	})
	return count, err
}

// generateSuggestions generates improvement suggestions based on issues
func (a *Analyzer) generateSuggestions(issues []models.Issue) []models.Suggestion {
	suggestions := []models.Suggestion{}

	// Count issues by category
	categoryCounts := make(map[string]int)
	for _, issue := range issues {
		categoryCounts[issue.Category]++
	}

	// Generate suggestions for high-count categories
	if categoryCounts["format"] > 5 {
		suggestions = append(suggestions, models.Suggestion{
			Title:       "Improve Code Formatting",
			Description: "Multiple formatting issues detected. Run 'gofmt' or 'goimports' to auto-fix.",
			Priority:    "medium",
			Category:    "format",
			Examples:    "go fmt ./... or goimports -w .",
		})
	}

	if categoryCounts["security"] > 0 {
		suggestions = append(suggestions, models.Suggestion{
			Title:       "Address Security Issues",
			Description: "Security vulnerabilities detected. Review and fix these issues immediately.",
			Priority:    "high",
			Category:    "security",
		})
	}

	if categoryCounts["performance"] > 3 {
		suggestions = append(suggestions, models.Suggestion{
			Title:       "Optimize Performance",
			Description: "Several performance issues found. Consider profiling and optimization.",
			Priority:    "medium",
			Category:    "performance",
		})
	}

	return suggestions
}

// getToolNames returns the names of all active linters
func (a *Analyzer) getToolNames() []string {
	names := make([]string, 0, len(a.linters))
	for name := range a.linters {
		names = append(names, name)
	}
	return names
}
