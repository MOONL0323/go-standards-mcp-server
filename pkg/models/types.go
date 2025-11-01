package models

import "time"

// AnalysisRequest represents a code analysis request
type AnalysisRequest struct {
	Code       string                 `json:"code,omitempty"`       // Code snippet to analyze
	FilePath   string                 `json:"file_path,omitempty"`  // Path to file
	ProjectDir string                 `json:"project_dir,omitempty"` // Path to project directory
	Standard   string                 `json:"standard"`             // strict, standard, relaxed, or custom
	Config     string                 `json:"config,omitempty"`     // Custom config content
	Format     string                 `json:"format"`               // json, markdown, html, pdf
	Options    map[string]interface{} `json:"options,omitempty"`    // Additional options
}

// AnalysisResult represents the result of code analysis
type AnalysisResult struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"` // success, error, partial
	Issues      []Issue                `json:"issues"`
	Summary     Summary                `json:"summary"`
	Metadata    Metadata               `json:"metadata"`
	Suggestions []Suggestion           `json:"suggestions,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Issue represents a single code issue
type Issue struct {
	File       string   `json:"file"`
	Line       int      `json:"line"`
	Column     int      `json:"column"`
	Severity   string   `json:"severity"` // error, warning, info
	Category   string   `json:"category"` // format, logic, security, performance, etc.
	Rule       string   `json:"rule"`
	Message    string   `json:"message"`
	Source     string   `json:"source"` // golangci-lint, staticcheck, etc.
	Code       string   `json:"code,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
}

// Summary provides statistics about the analysis
type Summary struct {
	TotalIssues    int            `json:"total_issues"`
	ErrorCount     int            `json:"error_count"`
	WarningCount   int            `json:"warning_count"`
	InfoCount      int            `json:"info_count"`
	FilesAnalyzed  int            `json:"files_analyzed"`
	LinesAnalyzed  int            `json:"lines_analyzed"`
	Duration       time.Duration  `json:"duration"`
	Score          float64        `json:"score"` // 0-100
	CategoryCounts map[string]int `json:"category_counts"`
}

// Metadata contains analysis metadata
type Metadata struct {
	Standard      string            `json:"standard"`
	ToolsUsed     []string          `json:"tools_used"`
	ConfigHash    string            `json:"config_hash"`
	GoVersion     string            `json:"go_version"`
	ServerVersion string            `json:"server_version"`
	Options       map[string]string `json:"options,omitempty"`
}

// Suggestion represents an improvement suggestion
type Suggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // high, medium, low
	Category    string `json:"category"`
	Examples    string `json:"examples,omitempty"`
}

// ConfigTemplate represents a predefined configuration template
type ConfigTemplate struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Level       string    `json:"level"` // strict, standard, relaxed
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CustomConfig represents a user-uploaded custom configuration
type CustomConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content"`
	Hash        string    `json:"hash"`
	Valid       bool      `json:"valid"`
	ValidationErrors []string `json:"validation_errors,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status    string            `json:"status"` // healthy, degraded, unhealthy
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    time.Duration     `json:"uptime"`
	Checks    map[string]string `json:"checks"`
}

// BatchAnalysisRequest represents a batch analysis request
type BatchAnalysisRequest struct {
	Projects []ProjectInfo          `json:"projects"`
	Standard string                 `json:"standard"`
	Config   string                 `json:"config,omitempty"`
	Format   string                 `json:"format"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ProjectInfo contains information about a project to analyze
type ProjectInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// BatchAnalysisResult contains results for multiple projects
type BatchAnalysisResult struct {
	ID        string                    `json:"id"`
	Status    string                    `json:"status"`
	Results   map[string]AnalysisResult `json:"results"` // project name -> result
	Summary   BatchSummary              `json:"summary"`
	CreatedAt time.Time                 `json:"created_at"`
}

// BatchSummary summarizes batch analysis results
type BatchSummary struct {
	TotalProjects     int           `json:"total_projects"`
	SuccessfulProjects int           `json:"successful_projects"`
	FailedProjects    int           `json:"failed_projects"`
	TotalIssues       int           `json:"total_issues"`
	AverageScore      float64       `json:"average_score"`
	Duration          time.Duration `json:"duration"`
}
