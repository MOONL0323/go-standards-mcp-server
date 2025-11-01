# API Documentation

## MCP Tools

This document describes all available MCP tools provided by the Go Standards MCP Server.

---

## 1. analyze_code

Analyze Go code and return detailed inspection results with issues, suggestions, and quality metrics.

### Input Schema

```json
{
  "code": "string (optional)",
  "file_path": "string (optional)",
  "project_dir": "string (optional)",
  "standard": "string (default: standard)",
  "config": "string (optional)",
  "format": "string (default: json)",
  "options": {
    "include_suggestions": "boolean (default: true)",
    "detailed": "boolean (default: false)"
  }
}
```

**Note**: Must provide one of: `code`, `file_path`, or `project_dir`

### Parameters

- **code**: Go code snippet to analyze
- **file_path**: Path to a single Go file
- **project_dir**: Path to a Go project directory
- **standard**: Configuration standard to use
  - `strict`: Highest standards for critical systems
  - `standard`: Balanced standards for general projects
  - `relaxed`: Basic standards for prototypes
  - `custom`: Use custom configuration (requires `config` parameter)
- **config**: Custom configuration in YAML format (required if standard is "custom")
- **format**: Output format
  - `json`: Structured JSON output
  - `markdown`: Human-readable Markdown
  - `html`: Rich HTML report (coming soon)
  - `pdf`: Printable PDF report (coming soon)
- **options**: Additional analysis options

### Output

```json
{
  "id": "uuid",
  "status": "success|error|partial",
  "issues": [
    {
      "file": "main.go",
      "line": 10,
      "column": 5,
      "severity": "error|warning|info",
      "category": "format|logic|security|performance|...",
      "rule": "gofmt",
      "message": "File is not formatted",
      "source": "golangci-lint",
      "code": "    func main() {",
      "suggestion": "Run gofmt on this file"
    }
  ],
  "summary": {
    "total_issues": 10,
    "error_count": 2,
    "warning_count": 5,
    "info_count": 3,
    "files_analyzed": 15,
    "lines_analyzed": 1500,
    "duration": "2.5s",
    "score": 85.5,
    "category_counts": {
      "format": 5,
      "logic": 3,
      "security": 2
    }
  },
  "suggestions": [
    {
      "title": "Improve Code Formatting",
      "description": "Multiple formatting issues detected...",
      "priority": "medium",
      "category": "format",
      "examples": "go fmt ./..."
    }
  ],
  "metadata": {
    "standard": "standard",
    "tools_used": ["golangci-lint", "govet"],
    "config_hash": "abc123",
    "go_version": "1.21",
    "server_version": "1.0.0"
  },
  "created_at": "2025-11-01T12:00:00Z"
}
```

### Example Usage

#### Analyze Code Snippet

```json
{
  "code": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}",
  "standard": "standard",
  "format": "json"
}
```

#### Analyze File

```json
{
  "file_path": "/path/to/main.go",
  "standard": "strict",
  "format": "markdown"
}
```

#### Analyze Project

```json
{
  "project_dir": "/path/to/project",
  "standard": "standard",
  "format": "json",
  "options": {
    "include_suggestions": true,
    "detailed": true
  }
}
```

#### Use Custom Configuration

```json
{
  "project_dir": "/path/to/project",
  "standard": "custom",
  "config": "linters:\n  enable:\n    - gofmt\n    - govet",
  "format": "json"
}
```

---

## 2. manage_config

Manage custom configuration files - upload, update, delete, or list configurations.

### Input Schema

```json
{
  "action": "upload|update|delete|list|get",
  "name": "string (optional)",
  "content": "string (optional)",
  "description": "string (optional)"
}
```

### Actions

- **list**: List all custom configurations
- **get**: Get a specific configuration
- **upload**: Upload a new configuration
- **update**: Update an existing configuration
- **delete**: Delete a configuration

### Example Usage

#### List Configurations

```json
{
  "action": "list"
}
```

#### Upload Configuration

```json
{
  "action": "upload",
  "name": "my-team-config",
  "content": "linters:\n  enable:\n    - gofmt\n    - govet\n...",
  "description": "Our team's coding standards"
}
```

---

## 3. manage_templates

Manage predefined configuration templates - list available templates and their details.

### Input Schema

```json
{
  "action": "list|get",
  "name": "string (optional)"
}
```

### Output

```json
[
  {
    "name": "strict",
    "display_name": "Strict Mode",
    "description": "Highest standards for critical systems (complexity ≤ 5, coverage ≥ 85%)",
    "level": "strict"
  },
  {
    "name": "standard",
    "display_name": "Standard Mode",
    "description": "Balanced standards for general projects (complexity ≤ 10, coverage ≥ 70%)",
    "level": "standard"
  },
  {
    "name": "relaxed",
    "display_name": "Relaxed Mode",
    "description": "Basic standards for prototypes (complexity ≤ 15, coverage ≥ 60%)",
    "level": "relaxed"
  }
]
```

### Example Usage

```json
{
  "action": "list"
}
```

---

## 4. generate_report

Generate analysis reports in various formats.

### Input Schema

```json
{
  "analysis_id": "string",
  "format": "json|markdown|html|pdf",
  "options": {}
}
```

### Example Usage

```json
{
  "analysis_id": "550e8400-e29b-41d4-a716-446655440000",
  "format": "html"
}
```

---

## 5. batch_analyze

Batch analyze multiple Go projects in parallel.

### Input Schema

```json
{
  "projects": [
    {
      "name": "string",
      "path": "string"
    }
  ],
  "standard": "string",
  "config": "string (optional)",
  "format": "string"
}
```

### Example Usage

```json
{
  "projects": [
    {
      "name": "service-a",
      "path": "/path/to/service-a"
    },
    {
      "name": "service-b",
      "path": "/path/to/service-b"
    }
  ],
  "standard": "standard",
  "format": "json"
}
```

---

## 6. health_check

Check the health status of the service and its dependencies.

### Input Schema

```json
{}
```

### Output

```json
{
  "status": "healthy|degraded|unhealthy",
  "timestamp": "2025-11-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": "24h30m15s",
  "checks": {
    "analyzer": "ok",
    "config": "ok",
    "linters": "ok"
  }
}
```

### Example Usage

```json
{}
```

---

## Error Handling

All tools return errors in the following format:

```json
{
  "error": {
    "code": "error_code",
    "message": "Error description",
    "details": {}
  }
}
```

Common error codes:
- `invalid_arguments`: Invalid input parameters
- `analysis_failed`: Code analysis failed
- `config_not_found`: Configuration not found
- `internal_error`: Internal server error

---

## Best Practices

1. **Start with standard mode**: Use the `standard` template for most projects
2. **Use custom configs for teams**: Create and share custom configurations
3. **Check health regularly**: Monitor service health in production
4. **Handle errors gracefully**: Implement proper error handling in your client
5. **Review suggestions**: Pay attention to improvement suggestions
6. **Iterate incrementally**: Fix high-priority issues first

---

## Rate Limits

- Analysis operations: No limits in self-hosted mode
- Concurrent analyses: Configurable via `analyzer.concurrent_limit`
- Report generation: No limits

---

## Support

For issues or questions:
- GitHub Issues: https://github.com/MOONL0323/go-standards-mcp-server/issues
- Documentation: See README.md
- Email: support@example.com
