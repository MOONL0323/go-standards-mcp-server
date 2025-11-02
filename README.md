# Go Standards MCP Server

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

基于 Model Context Protocol 的 Go 代码质量检测工具，集成 golangci-lint 和 go vet，遵循 Effective Go、Google Go Style Guide 等业界标准。

## 特性

- **多层次分析**：支持代码片段、文件、项目三个级别
- **标准化配置**：提供 strict、standard、relaxed 三种预设模板
- **两种使用方式**：MCP 服务器模式（IDE 集成）+ CLI 命令行模式
- **详细报告**：JSON 和 Markdown 格式输出
- **工具链集成**：golangci-lint (40+ linters) 和 go vet

## 快速开始

### 安装

```bash
# 克隆项目
git clone https://github.com/MOONL0323/go-standards-mcp-server.git
cd go-standards-mcp-server

# 安装依赖
go mod download

# 构建
make build-all    # 构建 MCP 服务器和 CLI 工具
```

构建完成后：
- MCP 服务器：`bin/mcp-server`
- CLI 工具：`bin/go-standards`

### 使用方式

#### 方式一：MCP 模式（Cursor IDE 集成）

**1. 配置 Cursor**

在 Cursor 设置中添加：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/your/path/to/bin/mcp-server"
    }
  }
}
```

**2. 重启 Cursor**

**3. 开始使用**

```
@go-standards 检查当前文件
@go-standards 用 strict 模式分析代码
@go-standards health_check
```

#### 方式二：CLI 模式（命令行）

**查看帮助**

```bash
./bin/go-standards --help
./bin/go-standards -help    # 查看详细帮助
```

**基本用法**

```bash
# 分析单个文件
./bin/go-standards -file main.go

# 分析整个项目
./bin/go-standards -project . -standard strict

# Markdown 格式输出
./bin/go-standards -file main.go -format markdown

# 分析代码片段
./bin/go-standards -code 'package main
func main() {
    x := 42
}'
```

## 分析模式

| 模式 | 复杂度阈值 | 覆盖率要求 | 适用场景 |
|------|-----------|-----------|---------|
| **strict** | ≤ 5 | ≥ 85% | 生产环境、关键系统 |
| **standard** | ≤ 10 | ≥ 70% | 日常开发（推荐） |
| **relaxed** | ≤ 15 | ≥ 60% | 原型开发、快速迭代 |

## 检测内容

- 未使用的变量和函数
- 未检查的错误
- 代码复杂度
- 代码格式问题
- 潜在 bug
- 性能问题
- 安全漏洞

## 实际应用

### Git Pre-commit Hook

```bash
#!/bin/bash
files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
for file in $files; do
    ./bin/go-standards -file "$file" || exit 1
done
```

### CI/CD 集成

```yaml
# GitHub Actions
- name: Code Quality Check
  run: ./bin/go-standards -project . -standard strict -format markdown > report.md
```

### 批量分析

```bash
find . -name "*.go" -not -path "./vendor/*" -exec ./bin/go-standards -file {} \;
```

## 自定义配置

### 方式一：本地配置文件（CLI 模式）

创建 `.golangci.yml` 文件：

```yaml
linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    - errcheck

linters-settings:
  gocyclo:
    min-complexity: 10
  govet:
    check-shadowing: true
```

使用：

```bash
./bin/go-standards -project . -config .golangci.yml
```

### 方式二：上传团队配置（MCP 模式）

#### 方法 A：上传 YAML 配置（推荐）

```javascript
// 直接上传 golangci-lint 配置
@go-standards manage_config {
  action: "upload",
  name: "team-standard",
  description: "团队代码规范 v1.0",
  content: "... golangci-lint YAML 内容 ..."
}
```

#### 方法 B：上传规范文档（AI 自动转换）🆕

```javascript
// 1. 上传 PDF/TXT/Markdown 格式的团队规范文档
@go-standards upload_document {
  content: "团队代码规范\n\n1. 函数复杂度≤10\n2. 必须检查所有错误...",
  file_name: "team-standard.md",
  name: "team-standard-v1",
  description: "团队代码规范 2025"
}
// 系统自动解析文档并转换为 golangci-lint 配置

// 2. 查看所有上传的文档
@go-standards list_documents

// 3. 获取文档和生成的配置
@go-standards get_document { id: "文档ID" }

// 4. 使用转换后的配置进行检查
@go-standards analyze_code {
  project_dir: "./myproject",
  standard: "custom",
  config: "team-standard-v1"
}
```

配置文件示例：`examples/team-config.yaml`

详细使用指南：[DOCUMENT_UPLOAD_GUIDE.md](DOCUMENT_UPLOAD_GUIDE.md)

## MCP 工具列表

| 工具 | 说明 | 状态 |
|-----|------|------|
| `analyze_code` | 代码质量分析 | ✅ 生产可用 |
| `manage_templates` | 管理配置模板 | ✅ 生产可用 |
| `manage_config` | 自定义配置管理（YAML上传） | ✅ 生产可用 |
| `upload_document` | 上传团队规范文档（PDF/TXT/MD）并自动转换 | ✅ 生产可用 |
| `list_documents` | 列出所有上传的规范文档 | ✅ 生产可用 |
| `get_document` | 获取文档详情和生成的配置 | ✅ 生产可用 |
| `delete_document` | 删除上传的文档 | ✅ 生产可用 |
| `health_check` | 服务健康检查 | ✅ 生产可用 |
| `generate_report` | 生成分析报告 | 🔜 开发中 |
| `batch_analyze` | 批量项目分析 | 🔜 开发中 |

## 项目结构

```
go-standards-mcp-server/
├── cmd/
│   ├── server/          # MCP 服务器
│   └── cli/             # CLI 工具
├── internal/
│   ├── mcp/            # MCP 协议实现
│   ├── analyzer/       # 分析引擎
│   ├── config/         # 配置管理
│   └── storage/        # 配置存储
├── pkg/
│   ├── models/         # 数据模型
│   └── linters/        # Linter 集成
├── configs/
│   ├── templates/      # 预设模板
│   └── custom/         # 自定义配置存储
└── examples/
    └── team-config.yaml # 团队配置示例
```

## 开发

```bash
# 运行测试
make test

# 代码检查
make lint

# 构建
make build-all

# 清理
make clean
```

## 常见问题

**Q: 如何在任何目录使用 CLI 工具？**

添加到 PATH：
```bash
export PATH=$PATH:/path/to/bin
# 或创建软链接
sudo ln -s /path/to/bin/go-standards /usr/local/bin/
```

**Q: MCP 模式和 CLI 模式有什么区别？**

- MCP 模式：集成到 IDE，实时分析，AI 辅助理解
- CLI 模式：独立运行，适合脚本和自动化

**Q: 支持哪些 linter？**

当前集成：golangci-lint、go vet  
计划集成：staticcheck、gosec

## 技术栈

- **语言**: Go 1.21+
- **MCP SDK**: [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- **配置管理**: [spf13/viper](https://github.com/spf13/viper)
- **日志系统**: [uber-go/zap](https://go.uber.org/zap)
- **分析工具**: golangci-lint, go vet

## 贡献

欢迎提交 Issue 和 Pull Request！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解开发规范。

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 致谢

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP Go SDK
- [golangci-lint](https://github.com/golangci/golangci-lint) - Go linters 聚合工具
- Go 社区的各类优秀工具和标准

---

**版本**: v1.0.0  
**作者**: MOONL0323  
**仓库**: [github.com/MOONL0323/go-standards-mcp-server](https://github.com/MOONL0323/go-standards-mcp-server)
