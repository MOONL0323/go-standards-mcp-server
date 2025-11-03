package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GitDetector handles git-based incremental detection
type GitDetector struct {
	repoPath string
}

// NewGitDetector creates a new git detector
func NewGitDetector(repoPath string) *GitDetector {
	return &GitDetector{
		repoPath: repoPath,
	}
}

// DiffMode represents different git diff modes
type DiffMode string

const (
	DiffModeStaged   DiffMode = "staged"   // git diff --cached
	DiffModeModified DiffMode = "modified" // git diff
	DiffModeBranch   DiffMode = "branch"   // git diff branch1..branch2
	DiffModeCommit   DiffMode = "commit"   // git diff commit1..commit2
)

// GetChangedFiles returns list of changed Go files
func (g *GitDetector) GetChangedFiles(mode DiffMode, args ...string) ([]string, error) {
	var cmd *exec.Cmd
	
	switch mode {
	case DiffModeStaged:
		// Get staged files (for pre-commit hook)
		cmd = exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	case DiffModeModified:
		// Get modified files in working directory
		cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=ACM")
	case DiffModeBranch:
		// Get files changed between branches
		if len(args) < 1 {
			return nil, fmt.Errorf("branch diff requires branch name argument")
		}
		cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=ACM", args[0])
	case DiffModeCommit:
		// Get files changed between commits
		if len(args) < 1 {
			return nil, fmt.Errorf("commit diff requires commit range argument")
		}
		cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=ACM", args[0])
	default:
		return nil, fmt.Errorf("unsupported diff mode: %s", mode)
	}
	
	cmd.Dir = g.repoPath
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git command failed: %w", err)
	}
	
	// Parse output and filter Go files
	var goFiles []string
	lines := strings.Split(out.String(), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Only include .go files
		if strings.HasSuffix(line, ".go") {
			// Convert to absolute path
			absPath := filepath.Join(g.repoPath, line)
			goFiles = append(goFiles, absPath)
		}
	}
	
	return goFiles, nil
}

// IsGitRepository checks if the path is a git repository
func (g *GitDetector) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.repoPath
	return cmd.Run() == nil
}

// GetCurrentBranch returns the current git branch name
func (g *GitDetector) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.repoPath
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	
	return strings.TrimSpace(out.String()), nil
}

// GetUnstagedFiles returns all unstaged Go files
func (g *GitDetector) GetUnstagedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--modified", "--exclude-standard")
	cmd.Dir = g.repoPath
	
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get unstaged files: %w", err)
	}
	
	var goFiles []string
	lines := strings.Split(out.String(), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasSuffix(line, ".go") {
			absPath := filepath.Join(g.repoPath, line)
			goFiles = append(goFiles, absPath)
		}
	}
	
	return goFiles, nil
}

// InstallGitHook installs a git hook
func (g *GitDetector) InstallGitHook(hookType string, content string) error {
	gitDir := filepath.Join(g.repoPath, ".git")
	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, hookType)
	
	// Create hook file
	file, err := os.Create(hookPath)
	if err != nil {
		return fmt.Errorf("failed to create hook file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write hook content: %w", err)
	}
	
	// Make executable (Unix systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(hookPath, 0755); err != nil {
			return fmt.Errorf("failed to make hook executable: %w", err)
		}
	}
	
	return nil
}

// GeneratePreCommitHook generates pre-commit hook script
func (g *GitDetector) GeneratePreCommitHook(serverPath string) string {
	return fmt.Sprintf(`#!/bin/sh
# Auto-generated pre-commit hook for go-standards-mcp-server

echo "Running code quality checks on staged files..."

# Run go-standards with staged files check
%s --git-mode staged --auto

if [ $? -ne 0 ]; then
    echo "Code quality check failed. Please fix the issues before committing."
    exit 1
fi

echo "Code quality check passed."
exit 0
`, serverPath)
}

// GeneratePrePushHook generates pre-push hook script
func (g *GitDetector) GeneratePrePushHook(serverPath string, baseBranch string) string {
	if baseBranch == "" {
		baseBranch = "origin/main"
	}
	
	return fmt.Sprintf(`#!/bin/sh
# Auto-generated pre-push hook for go-standards-mcp-server

echo "Running code quality checks on branch changes..."

# Run go-standards with branch diff check
%s --git-mode branch --git-ref %s --auto

if [ $? -ne 0 ]; then
    echo "Code quality check failed. Please fix the issues before pushing."
    exit 1
fi

echo "Code quality check passed."
exit 0
`, serverPath, baseBranch)
}
