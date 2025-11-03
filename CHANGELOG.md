# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

**[中文版本](CHANGELOG_CN.md)**

---

## [1.1.0] - 2025-11-03

### Added
- **Git Incremental Checking**: Analyze only changed files
  - Three modes: staged, modified, branch comparison
  - Pre-commit and pre-push hook support
  - Automatic hook installation
  - Configurable base branch
- **Multi-User Architecture**: Support concurrent users
  - User workspace isolation
  - Shared team standards
  - Session management with timeout
  - Per-user resource limits
- **Document Upload**: Team standards management
  - PDF/TXT/Markdown support
  - Automatic rule extraction
  - Auto-generate golangci-lint configs
- **CLI Tool**: Standalone command-line interface
  - All MCP features in CLI
  - Git integration
  - CI/CD ready
- **Platform Independence**: Generic module path
  - Module: `go-standards-mcp-server`
  - Easy migration to any Git platform

### New Tools
- `upload_document`: Upload coding standards
- `git_config`: Configure Git checking
- `git_check`: Check if path is Git repo

### Improved
- Complete README rewrite
- ASCII architecture diagrams
- Professional project structure
- Bilingual documentation (EN/CN)

### Configuration
- `.go-standards-git.yaml`: Git integration config

---

## [1.0.0] - 2025-11-01

### Added
- Initial release
- MCP protocol support (stdio)
- golangci-lint and go vet integration
- Three config templates (strict/standard/relaxed)
- Detailed issue reporting
- JSON and Markdown reports
- Docker support

### Features
- `analyze_go_code`: Analyze with configurable standards
- `list_standards`: List coding standards
- `get_config`/`set_config`: Manage configs

### Templates
- **Strict**: complexity ≤ 5, coverage ≥ 85%
- **Standard**: complexity ≤ 10, coverage ≥ 70%
- **Relaxed**: complexity ≤ 15, coverage ≥ 60%

---

## Roadmap

### v1.2 (Planned)
- HTML/PDF report generation
- Batch project analysis
- Performance metrics tracking
- Web dashboard

### v2.0 (Future)
- IDE extensions (VS Code, GoLand)
- Team collaboration features
- Custom linter plugins
- Multi-language support

---

[1.1.0]: https://github.com/MOONL0323/go-standards-mcp-server/releases/tag/v1.1.0
[1.0.0]: https://github.com/MOONL0323/go-standards-mcp-server/releases/tag/v1.0.0
