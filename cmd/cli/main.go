package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/MOONL0323/go-standards-mcp-server/internal/analyzer"
	"github.com/MOONL0323/go-standards-mcp-server/internal/config"
	"github.com/MOONL0323/go-standards-mcp-server/pkg/models"
	"go.uber.org/zap"
)

var (
	filePath   = flag.String("file", "", "Path to Go file to analyze")
	projectDir = flag.String("project", "", "Path to Go project directory")
	code       = flag.String("code", "", "Go code snippet to analyze")
	standard   = flag.String("standard", "standard", "Analysis standard: strict, standard, or relaxed")
	format     = flag.String("format", "json", "Output format: json or markdown")
	configPath = flag.String("config", "", "Path to custom config file")
	version    = flag.Bool("version", false, "Print version and exit")
	help       = flag.Bool("help", false, "Show detailed help message")
)

const (
	appName    = "go-standards"
	appVersion = "1.0.0"
)

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *help {
		printDetailedHelp()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Validate input
	if *filePath == "" && *projectDir == "" && *code == "" {
		fmt.Fprintln(os.Stderr, "Error: Must specify one of -file, -project, or -code\n")
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger (quiet mode for CLI)
	logger := zap.NewNop()

	// Initialize analyzer
	a, err := analyzer.NewAnalyzer(cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize analyzer: %v\n", err)
		os.Exit(1)
	}

	// Create analysis request
	req := &models.AnalysisRequest{
		Code:       *code,
		FilePath:   *filePath,
		ProjectDir: *projectDir,
		Standard:   *standard,
		Format:     *format,
	}

	// Perform analysis
	result, err := a.Analyze(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if *format == "json" {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format output: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(formatMarkdown(result))
	}

	// Exit with error code if issues found
	if result.Summary.ErrorCount > 0 {
		os.Exit(1)
	}
}

func formatMarkdown(result *models.AnalysisResult) string {
	md := fmt.Sprintf("# Code Analysis Report\n\n")
	md += fmt.Sprintf("Status: %s\n", result.Status)
	md += fmt.Sprintf("Score: %.1f/100\n\n", result.Summary.Score)
	
	md += fmt.Sprintf("## Summary\n\n")
	md += fmt.Sprintf("- Total Issues: %d\n", result.Summary.TotalIssues)
	md += fmt.Sprintf("- Errors: %d\n", result.Summary.ErrorCount)
	md += fmt.Sprintf("- Warnings: %d\n", result.Summary.WarningCount)
	md += fmt.Sprintf("- Files Analyzed: %d\n", result.Summary.FilesAnalyzed)
	md += fmt.Sprintf("- Duration: %s\n\n", result.Summary.Duration)

	if len(result.Issues) > 0 {
		md += "## Issues\n\n"
		for i, issue := range result.Issues {
			md += fmt.Sprintf("%d. [%s] %s\n", i+1, issue.Severity, issue.Message)
			md += fmt.Sprintf("   File: %s:%d:%d\n", issue.File, issue.Line, issue.Column)
			md += fmt.Sprintf("   Rule: %s (%s)\n\n", issue.Rule, issue.Source)
		}
	}

	if len(result.Suggestions) > 0 {
		md += "## Suggestions\n\n"
		for i, sug := range result.Suggestions {
			md += fmt.Sprintf("%d. [%s] %s\n", i+1, sug.Priority, sug.Title)
			md += fmt.Sprintf("   %s\n\n", sug.Description)
		}
	}

	return md
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [options]

Go code quality analysis tool with multiple standards.

OPTIONS:
  -file string
        Analyze a single Go file
        Example: -file main.go

  -project string
        Analyze entire Go project directory
        Example: -project .
        Example: -project /path/to/project

  -code string
        Analyze Go code snippet (use quotes)
        Example: -code 'package main; func main() { }'

  -standard string
        Analysis standard level (default: standard)
        Options:
          strict   - Highest standards (complexity ≤ 5, coverage ≥ 85%%)
          standard - Balanced standards (complexity ≤ 10, coverage ≥ 70%%)
          relaxed  - Basic standards (complexity ≤ 15, coverage ≥ 60%%)

  -format string
        Output format (default: json)
        Options:
          json     - Structured JSON output
          markdown - Human-readable Markdown report

  -config string
        Path to custom golangci-lint config file
        Example: -config .golangci.yml

  -version
        Print version information and exit

  -help
        Show this detailed help message

EXAMPLES:
  # Analyze a single file
  %s -file main.go

  # Analyze with strict mode
  %s -file main.go -standard strict

  # Analyze entire project with markdown output
  %s -project . -format markdown

  # Analyze code snippet
  %s -code 'package main
  func main() {
      x := 42
      println("hello")
  }'

  # Use custom config
  %s -project . -config .golangci.yml

EXIT CODES:
  0  Analysis successful, no errors found
  1  Analysis failed or errors detected

For more information, visit: https://github.com/MOONL0323/go-standards-mcp-server
`, appName, appName, appName, appName, appName, appName)
}

func printDetailedHelp() {
	printUsage()
	fmt.Fprintf(os.Stderr, `
ANALYSIS STANDARDS:

  strict:
    - Cyclomatic complexity ≤ 5
    - Test coverage ≥ 85%%
    - All linters enabled
    - Best for: Production systems, critical code

  standard: (recommended)
    - Cyclomatic complexity ≤ 10
    - Test coverage ≥ 70%%
    - Most linters enabled
    - Best for: Daily development, general projects

  relaxed:
    - Cyclomatic complexity ≤ 15
    - Test coverage ≥ 60%%
    - Core linters enabled
    - Best for: Prototypes, rapid development

WHAT IT CHECKS:

  ✓ Unused variables and functions
  ✓ Unchecked errors
  ✓ Code complexity
  ✓ Code formatting issues
  ✓ Potential bugs
  ✓ Performance issues
  ✓ Security vulnerabilities

OUTPUT FORMATS:

  json:
    - Complete structured data
    - Issue details with line numbers
    - Statistics and metrics
    - Machine-readable for automation

  markdown:
    - Human-readable report
    - Summary and issue list
    - Improvement suggestions
    - Great for documentation

USAGE IN SCRIPTS:

  # Git pre-commit hook
  %s -file "$file" || exit 1

  # CI/CD pipeline
  %s -project . -standard strict -format markdown > report.md

  # Batch analysis
  find . -name "*.go" -exec %s -file {} \;

TIPS:

  • Start with 'standard' mode for daily development
  • Use 'strict' mode before production deployment
  • Use 'relaxed' mode for quick prototyping
  • Combine with git hooks for automatic checking
  • Integrate into CI/CD for continuous quality control

`, appName, appName, appName)
}


