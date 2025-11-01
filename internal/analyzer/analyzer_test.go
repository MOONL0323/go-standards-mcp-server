package analyzer

import (
	"context"
	"testing"
	"time"

	"github.com/MOONL0323/go-standards-mcp-server/internal/config"
	"github.com/MOONL0323/go-standards-mcp-server/pkg/models"
	"go.uber.org/zap"
)

func TestAnalyzer_Analyze(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		Analyzer: config.AnalyzerConfig{
			Timeout:         5 * time.Minute,
			ConcurrentLimit: 5,
			TempDir:         "../../tmp",
		},
		Linters: config.LintersConfig{
			GolangciLint: config.GolangciLintConfig{
				Enabled: false, // Disable for unit test
			},
			Govet: config.LinterConfig{
				Enabled: true,
			},
		},
	}

	analyzer, err := NewAnalyzer(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	tests := []struct {
		name    string
		request *models.AnalysisRequest
		wantErr bool
	}{
		{
			name: "analyze simple code",
			request: &models.AnalysisRequest{
				Code:     "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
				Standard: "standard",
				Format:   "json",
			},
			wantErr: false,
		},
		{
			name: "analyze with strict standard",
			request: &models.AnalysisRequest{
				Code:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n",
				Standard: "strict",
				Format:   "json",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := analyzer.Analyze(ctx, tt.request)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if result.ID == "" {
					t.Error("Result ID is empty")
				}
				if result.Status == "" {
					t.Error("Result status is empty")
				}
			}
		})
	}
}

func TestAnalyzer_calculateSummary(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		Analyzer: config.AnalyzerConfig{
			TempDir: "../../tmp",
		},
	}

	analyzer, _ := NewAnalyzer(cfg, logger)

	issues := []models.Issue{
		{Severity: "error", Category: "logic"},
		{Severity: "warning", Category: "format"},
		{Severity: "warning", Category: "format"},
		{Severity: "info", Category: "style"},
	}

	summary := analyzer.calculateSummary(issues, ".", 1*time.Second)

	if summary.TotalIssues != 4 {
		t.Errorf("Expected 4 total issues, got %d", summary.TotalIssues)
	}
	if summary.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", summary.ErrorCount)
	}
	if summary.WarningCount != 2 {
		t.Errorf("Expected 2 warnings, got %d", summary.WarningCount)
	}
	if summary.InfoCount != 1 {
		t.Errorf("Expected 1 info, got %d", summary.InfoCount)
	}
	if summary.Score >= 100 || summary.Score < 0 {
		t.Errorf("Score out of range: %f", summary.Score)
	}
}
