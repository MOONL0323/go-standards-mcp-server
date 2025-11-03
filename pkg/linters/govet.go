package linters

import (
	"context"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"go-standards-mcp-server/pkg/models"
	"go.uber.org/zap"
)

// GoVet implements the go vet linter
type GoVet struct {
	logger *zap.Logger
}

// NewGoVet creates a new GoVet instance
func NewGoVet(logger *zap.Logger) *GoVet {
	return &GoVet{
		logger: logger,
	}
}

// Name returns the name of the linter
func (g *GoVet) Name() string {
	return "govet"
}

// IsAvailable checks if go vet is available
func (g *GoVet) IsAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

// Run executes go vet
func (g *GoVet) Run(ctx context.Context, workDir, configPath string) ([]models.Issue, error) {
	cmd := exec.CommandContext(ctx, "go", "vet", "./...")
	cmd.Dir = workDir

	g.logger.Debug("Running go vet", zap.String("workDir", workDir))

	output, err := cmd.CombinedOutput()
	
	// go vet returns non-zero exit code when issues are found
	// Parse the output regardless of exit code
	issues := g.parseOutput(workDir, string(output))

	g.logger.Debug("go vet completed",
		zap.Int("issues", len(issues)),
		zap.Error(err))

	return issues, nil
}

// parseOutput parses go vet output
func (g *GoVet) parseOutput(workDir, output string) []models.Issue {
	if output == "" {
		return []models.Issue{}
	}

	// go vet output format: path/file.go:line:column: message
	// or: path/file.go:line: message
	re := regexp.MustCompile(`^(.+?):(\d+):(?:(\d+):)?\s*(.+)$`)

	var issues []models.Issue
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		file := matches[1]
		lineNum, _ := strconv.Atoi(matches[2])
		column := 0
		message := matches[4]

		if matches[3] != "" {
			column, _ = strconv.Atoi(matches[3])
		}

		// Make path relative
		relPath, err := filepath.Rel(workDir, file)
		if err != nil {
			relPath = file
		}

		issue := models.Issue{
			File:     relPath,
			Line:     lineNum,
			Column:   column,
			Severity: "warning",
			Category: g.categorize(message),
			Rule:     "govet",
			Message:  message,
			Source:   "govet",
		}

		issues = append(issues, issue)
	}

	return issues
}

// categorize attempts to categorize the issue based on the message
func (g *GoVet) categorize(message string) string {
	message = strings.ToLower(message)

	if strings.Contains(message, "shadow") {
		return "logic"
	}
	if strings.Contains(message, "printf") || strings.Contains(message, "format") {
		return "format"
	}
	if strings.Contains(message, "composite literal") {
		return "style"
	}
	if strings.Contains(message, "unreachable") {
		return "dead-code"
	}
	if strings.Contains(message, "nil") {
		return "error-handling"
	}

	return "logic"
}
