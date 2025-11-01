package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAnalysisResult_JSON(t *testing.T) {
	result := &AnalysisResult{
		ID:     "test-123",
		Status: "success",
		Issues: []Issue{
			{
				File:     "main.go",
				Line:     10,
				Column:   5,
				Severity: "warning",
				Category: "format",
				Rule:     "gofmt",
				Message:  "File is not formatted",
				Source:   "golangci-lint",
			},
		},
		Summary: Summary{
			TotalIssues:  1,
			WarningCount: 1,
			Score:        98.0,
			Duration:     1 * time.Second,
		},
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	var decoded AnalysisResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if decoded.ID != result.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, result.ID)
	}
	if decoded.Summary.TotalIssues != result.Summary.TotalIssues {
		t.Errorf("TotalIssues mismatch: got %d, want %d",
			decoded.Summary.TotalIssues, result.Summary.TotalIssues)
	}
}

func TestHealthStatus_JSON(t *testing.T) {
	status := &HealthStatus{
		Status:    "healthy",
		Version:   "1.0.0",
		Uptime:    1 * time.Hour,
		Checks: map[string]string{
			"analyzer": "ok",
			"config":   "ok",
		},
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal status: %v", err)
	}

	var decoded HealthStatus
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal status: %v", err)
	}

	if decoded.Status != status.Status {
		t.Errorf("Status mismatch: got %s, want %s", decoded.Status, status.Status)
	}
}
