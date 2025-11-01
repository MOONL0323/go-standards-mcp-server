# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-01

### Added
- Initial release of Go Standards MCP Server
- MCP protocol support (stdio mode)
- Integration with golangci-lint and go vet
- Three predefined configuration templates (strict, standard, relaxed)
- Code analysis with detailed issue reporting
- JSON and Markdown report formats
- Docker support for easy deployment
- Comprehensive documentation and examples
- Unit tests for core components

### Features
- **analyze_code**: Analyze Go code with configurable standards
- **manage_config**: Manage custom configurations (stub implementation)
- **manage_templates**: List and view predefined templates
- **generate_report**: Generate analysis reports (stub implementation)
- **batch_analyze**: Batch analysis support (stub implementation)
- **health_check**: Service health monitoring

### Configuration Templates
- **Strict Mode**: Highest standards (complexity ≤ 5, coverage ≥ 85%)
- **Standard Mode**: Balanced standards (complexity ≤ 10, coverage ≥ 70%)
- **Relaxed Mode**: Basic standards (complexity ≤ 15, coverage ≥ 60%)

### Tools Integrated
- golangci-lint: Comprehensive linter aggregator
- go vet: Official Go static analysis tool

### Documentation
- Complete README with usage examples
- Contributing guidelines
- Docker deployment guide
- API documentation for MCP tools

## [Unreleased]

### Planned
- HTTP/SSE mode support
- Additional linter integrations (staticcheck, gosec)
- HTML and PDF report generation
- Custom configuration upload and management
- Batch analysis implementation
- PostgreSQL support
- Redis caching
- Incremental analysis
- CI/CD integration examples
- VS Code extension
- Web dashboard

---

[1.0.0]: https://github.com/MOONL0323/go-standards-mcp-server/releases/tag/v1.0.0
