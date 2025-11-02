package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MOONL0323/go-standards-mcp-server/internal/analyzer"
	"github.com/MOONL0323/go-standards-mcp-server/internal/config"
	"github.com/MOONL0323/go-standards-mcp-server/internal/storage"
	"github.com/MOONL0323/go-standards-mcp-server/pkg/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

const (
	ServerName    = "go-standards-mcp-server"
	ServerVersion = "1.0.0"
)

// Server represents the MCP server
type Server struct {
	config        *config.Config
	logger        *zap.Logger
	analyzer      *analyzer.Analyzer
	srv           *server.MCPServer
	configStorage *storage.ConfigStorage
	docService    interface{} // DocumentService interface (避免循环依赖)
}

// NewServer creates a new MCP server instance
func NewServer(cfg *config.Config, logger *zap.Logger, analyzer *analyzer.Analyzer) (*Server, error) {
	// Initialize config storage
	configStorage, err := storage.NewConfigStorage("./configs/custom")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config storage: %w", err)
	}

	s := &Server{
		config:        cfg,
		logger:        logger,
		analyzer:      analyzer,
		configStorage: configStorage,
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithToolCapabilities(true),
	)

	// Register tools
	if err := s.registerTools(mcpServer); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	s.srv = mcpServer
	return s, nil
}

// registerTools registers all available MCP tools
func (s *Server) registerTools(mcpServer *server.MCPServer) error {
	tools := []struct {
		name        string
		description string
		schema      mcp.ToolInputSchema
		handler     server.ToolHandlerFunc
	}{
		{
			name:        "analyze_code",
			description: "Analyze Go code and return detailed inspection results with issues, suggestions, and quality metrics",
			schema:      s.getAnalyzeCodeSchema(),
			handler:     s.handleAnalyzeCode,
		},
		{
			name:        "manage_config",
			description: "Manage custom configuration files - upload, update, delete, or list configurations",
			schema:      s.getManageConfigSchema(),
			handler:     s.handleManageConfig,
		},
		{
			name:        "manage_templates",
			description: "Manage predefined configuration templates - list available templates and their details",
			schema:      s.getManageTemplatesSchema(),
			handler:     s.handleManageTemplates,
		},
		{
			name:        "generate_report",
			description: "Generate analysis reports in various formats (JSON, Markdown, HTML, PDF)",
			schema:      s.getGenerateReportSchema(),
			handler:     s.handleGenerateReport,
		},
		{
			name:        "batch_analyze",
			description: "Batch analyze multiple Go projects in parallel",
			schema:      s.getBatchAnalyzeSchema(),
			handler:     s.handleBatchAnalyze,
		},
		{
			name:        "health_check",
			description: "Check the health status of the service and its dependencies",
			schema:      s.getHealthCheckSchema(),
			handler:     s.handleHealthCheck,
		},
		{
			name:        "upload_document",
			description: "Upload team code standard document (PDF, TXT, Markdown) and auto-convert to golangci-lint config",
			schema:      s.getUploadDocumentSchema(),
			handler:     s.handleUploadDocument,
		},
		{
			name:        "list_documents",
			description: "List all uploaded team standard documents with metadata",
			schema:      s.getListDocumentsSchema(),
			handler:     s.handleListDocuments,
		},
		{
			name:        "get_document",
			description: "Get uploaded document details and generated configuration",
			schema:      s.getGetDocumentSchema(),
			handler:     s.handleGetDocument,
		},
		{
			name:        "delete_document",
			description: "Delete an uploaded document and its configuration",
			schema:      s.getDeleteDocumentSchema(),
			handler:     s.handleDeleteDocument,
		},
	}

	for _, tool := range tools {
		mcpServer.AddTool(mcp.Tool{
			Name:        tool.name,
			Description: tool.description,
			InputSchema: tool.schema,
		}, tool.handler)
		s.logger.Info("Registered tool", zap.String("name", tool.name))
	}

	return nil
}

// getAnalyzeCodeSchema returns the JSON schema for analyze_code tool
func (s *Server) getAnalyzeCodeSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Go code snippet to analyze",
			},
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "Path to a single Go file to analyze",
			},
			"project_dir": map[string]interface{}{
				"type":        "string",
				"description": "Path to a Go project directory to analyze",
			},
			"standard": map[string]interface{}{
				"type":        "string",
				"description": "Configuration standard to use: strict, standard, relaxed, or custom",
				"enum":        []string{"strict", "standard", "relaxed", "custom"},
				"default":     "standard",
			},
			"config": map[string]interface{}{
				"type":        "string",
				"description": "Custom configuration content (YAML format, required if standard is 'custom')",
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Output format for the analysis result",
				"enum":        []string{"json", "markdown", "html", "pdf"},
				"default":     "json",
			},
			"options": map[string]interface{}{
				"type":        "object",
				"description": "Additional analysis options",
				"properties": map[string]interface{}{
					"include_suggestions": map[string]interface{}{
						"type":    "boolean",
						"default": true,
					},
					"detailed": map[string]interface{}{
						"type":    "boolean",
						"default": false,
					},
				},
			},
		},
	}
}

// getManageConfigSchema returns the JSON schema for manage_config tool
func (s *Server) getManageConfigSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform",
				"enum":        []string{"upload", "update", "delete", "list", "get"},
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Configuration name (required for upload, update, delete, get)",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Configuration content in YAML format (required for upload, update)",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Configuration description (optional)",
			},
		},
	}
}

// getManageTemplatesSchema returns the JSON schema for manage_templates tool
func (s *Server) getManageTemplatesSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform",
				"enum":        []string{"list", "get"},
				"default":     "list",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Template name (required for get action)",
			},
		},
	}
}

// getGenerateReportSchema returns the JSON schema for generate_report tool
func (s *Server) getGenerateReportSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"analysis_id": map[string]interface{}{
				"type":        "string",
				"description": "ID of the analysis to generate report for",
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Report format",
				"enum":        []string{"json", "markdown", "html", "pdf"},
				"default":     "markdown",
			},
			"options": map[string]interface{}{
				"type":        "object",
				"description": "Report generation options",
			},
		},
	}
}

// getBatchAnalyzeSchema returns the JSON schema for batch_analyze tool
func (s *Server) getBatchAnalyzeSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"projects": map[string]interface{}{
				"type":        "array",
				"description": "List of projects to analyze",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Project name",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Project directory path",
						},
					},
					"required": []string{"name", "path"},
				},
			},
			"standard": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"strict", "standard", "relaxed", "custom"},
				"default": "standard",
			},
			"config": map[string]interface{}{
				"type":        "string",
				"description": "Custom configuration content",
			},
			"format": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"json", "markdown", "html"},
				"default": "json",
			},
		},
	}
}

// getHealthCheckSchema returns the JSON schema for health_check tool
func (s *Server) getHealthCheckSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:       "object",
		Properties: map[string]interface{}{},
	}
}

// handleAnalyzeCode handles the analyze_code tool invocation
func (s *Server) handleAnalyzeCode(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	s.logger.Info("Handling analyze_code request")

	// Parse arguments
	var req models.AnalysisRequest
	if err := parseArguments(arguments, &req); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Set defaults
	if req.Standard == "" {
		req.Standard = "standard"
	}
	if req.Format == "" {
		req.Format = "json"
	}

	// Perform analysis
	result, err := s.analyzer.Analyze(context.Background(), &req)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	// Format result
	content, err := formatResult(result, req.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to format result: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

// handleManageConfig handles the manage_config tool invocation
func (s *Server) handleManageConfig(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	s.logger.Info("Handling manage_config request")

	var args struct {
		Action      string `json:"action"`
		Name        string `json:"name"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}

	if err := parseArguments(arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	var response string
	var err error

	switch args.Action {
	case "list":
		configs, err := s.configStorage.List()
		if err != nil {
			return nil, fmt.Errorf("failed to list configs: %w", err)
		}
		data, err := json.MarshalIndent(configs, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal configs: %w", err)
		}
		response = string(data)

	case "upload", "update":
		if args.Name == "" || args.Content == "" {
			return nil, fmt.Errorf("name and content are required")
		}
		if err := s.configStorage.Save(args.Name, args.Content, args.Description); err != nil {
			return nil, fmt.Errorf("failed to save config: %w", err)
		}
		response = fmt.Sprintf(`{"message": "Config '%s' saved successfully"}`, args.Name)

	case "get":
		if args.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		config, err := s.configStorage.Get(args.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get config: %w", err)
		}
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		response = string(data)

	case "delete":
		if args.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		if err := s.configStorage.Delete(args.Name); err != nil {
			return nil, fmt.Errorf("failed to delete config: %w", err)
		}
		response = fmt.Sprintf(`{"message": "Config '%s' deleted successfully"}`, args.Name)

	default:
		return nil, fmt.Errorf("unknown action: %s (valid actions: list, upload, update, get, delete)", args.Action)
	}

	_ = err // avoid unused variable warning

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

// handleManageTemplates handles the manage_templates tool invocation
func (s *Server) handleManageTemplates(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	s.logger.Info("Handling manage_templates request")

	templates := []models.ConfigTemplate{
		{
			Name:        "strict",
			DisplayName: "Strict Mode",
			Description: "Highest standards for critical systems (complexity ≤ 5, coverage ≥ 85%)",
			Level:       "strict",
		},
		{
			Name:        "standard",
			DisplayName: "Standard Mode",
			Description: "Balanced standards for general projects (complexity ≤ 10, coverage ≥ 70%)",
			Level:       "standard",
		},
		{
			Name:        "relaxed",
			DisplayName: "Relaxed Mode",
			Description: "Basic standards for prototypes (complexity ≤ 15, coverage ≥ 60%)",
			Level:       "relaxed",
		},
	}

	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal templates: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: string(data),
			},
		},
	}, nil
}

// handleGenerateReport handles the generate_report tool invocation
func (s *Server) handleGenerateReport(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	s.logger.Info("Handling generate_report request")

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: `{"message": "Report generation coming soon"}`,
			},
		},
	}, nil
}

// handleBatchAnalyze handles the batch_analyze tool invocation
func (s *Server) handleBatchAnalyze(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	s.logger.Info("Handling batch_analyze request")

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: `{"message": "Batch analysis coming soon"}`,
			},
		},
	}, nil
}

// handleHealthCheck handles the health_check tool invocation
func (s *Server) handleHealthCheck(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	health := models.HealthStatus{
		Status:  "healthy",
		Version: ServerVersion,
		Checks: map[string]string{
			"analyzer": "ok",
			"config":   "ok",
		},
	}

	data, err := json.MarshalIndent(health, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal health status: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: string(data),
			},
		},
	}, nil
}

// Document management handlers

// handleUploadDocument handles document upload
func (s *Server) handleUploadDocument(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	var args struct {
		Content     string `json:"content"`      // Base64 encoded file content or text content
		FileName    string `json:"file_name"`    // File name with extension
		Name        string `json:"name"`         // Configuration name
		Description string `json:"description"`  // Description
	}

	if err := parseArguments(arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	response := fmt.Sprintf(`{
		"message": "Document upload feature requires document service initialization",
		"file": "%s",
		"name": "%s",
		"note": "This feature is available in the next version. For now, use manage_config to upload YAML configurations directly."
	}`, args.FileName, args.Name)

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

// handleListDocuments lists all uploaded documents
func (s *Server) handleListDocuments(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	response := `{
		"documents": [],
		"message": "Document management feature coming soon. Use manage_config to list configurations."
	}`

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

// handleGetDocument gets document details
func (s *Server) handleGetDocument(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := parseArguments(arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	response := fmt.Sprintf(`{
		"message": "Document details feature coming soon",
		"document_id": "%s"
	}`, args.ID)

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

// handleDeleteDocument deletes a document
func (s *Server) handleDeleteDocument(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := parseArguments(arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	response := fmt.Sprintf(`{
		"message": "Document deletion feature coming soon",
		"document_id": "%s"
	}`, args.ID)

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

// Serve starts the MCP server
func (s *Server) Serve() error {
	s.logger.Info("Starting MCP server", zap.String("mode", s.config.Server.Mode))

	if err := server.ServeStdio(s.srv); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// parseArguments parses MCP tool arguments into a struct
func parseArguments(args interface{}, target interface{}) error {
	data, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal arguments: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	return nil
}

// formatResult formats the analysis result in the specified format
func formatResult(result *models.AnalysisResult, format string) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "markdown":
		return formatMarkdown(result), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// formatMarkdown formats the result as Markdown
func formatMarkdown(result *models.AnalysisResult) string {
	md := fmt.Sprintf("# Code Analysis Report\n\n")
	md += fmt.Sprintf("**Status**: %s\n", result.Status)
	md += fmt.Sprintf("**Score**: %.1f/100\n\n", result.Summary.Score)
	md += fmt.Sprintf("## Summary\n\n")
	md += fmt.Sprintf("- Total Issues: %d\n", result.Summary.TotalIssues)
	md += fmt.Sprintf("- Errors: %d\n", result.Summary.ErrorCount)
	md += fmt.Sprintf("- Warnings: %d\n", result.Summary.WarningCount)
	md += fmt.Sprintf("- Files Analyzed: %d\n", result.Summary.FilesAnalyzed)
	md += fmt.Sprintf("- Duration: %s\n\n", result.Summary.Duration)

	if len(result.Issues) > 0 {
		md += "## Issues\n\n"
		for i, issue := range result.Issues {
			if i >= 10 { // Limit to first 10 issues in markdown
				md += fmt.Sprintf("... and %d more issues\n", len(result.Issues)-10)
				break
			}
			md += fmt.Sprintf("### %d. %s\n", i+1, issue.Message)
			md += fmt.Sprintf("- **File**: %s:%d:%d\n", issue.File, issue.Line, issue.Column)
			md += fmt.Sprintf("- **Severity**: %s\n", issue.Severity)
			md += fmt.Sprintf("- **Category**: %s\n", issue.Category)
			md += fmt.Sprintf("- **Rule**: %s\n\n", issue.Rule)
		}
	}

	return md
}

// Document management schemas

func (s *Server) getUploadDocumentSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Document content (plain text for TXT/MD, base64 for PDF)",
			},
			"file_name": map[string]interface{}{
				"type":        "string",
				"description": "File name with extension (e.g., team-standard.pdf, coding-rules.md)",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Configuration name for this standard (e.g., team-standard-v1)",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Description of the coding standard",
			},
		},
	}
}

func (s *Server) getListDocumentsSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:       "object",
		Properties: map[string]interface{}{},
	}
}

func (s *Server) getGetDocumentSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Document ID",
			},
		},
	}
}

func (s *Server) getDeleteDocumentSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Document ID to delete",
			},
		},
	}
}
