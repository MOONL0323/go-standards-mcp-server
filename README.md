# Go Standards MCP Server

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

基于权威标准的 Go 代码规范检测 MCP 服务器，集成 golangci-lint、staticcheck、gosec 等专业工具链，支持团队自定义配置。

## 🌟 核心特性

- **权威标准**: 基于 Effective Go、Google Go Style Guide、Uber Go Style Guide
- **灵活配置**: 支持预设模板（严格/标准/宽松）和自定义配置
- **专业工具链**: 集成 golangci-lint、staticcheck、gosec、go vet
- **多格式报告**: 支持 JSON、Markdown、HTML、PDF 格式输出
- **MCP 协议**: 无缝集成 Cursor IDE、Claude Code、VS Code
- **团队协作**: 配置共享、版本管理、权限控制

## 📋 快速开始

### 前置要求

- Go 1.21 或更高版本
- Docker 和 Docker Compose (可选)
- golangci-lint (自动安装)

### 安装

```bash
# 克隆项目
git clone https://github.com/MOONL0323/go-standards-mcp-server.git
cd go-standards-mcp-server

# 安装依赖
go mod download

# 构建
go build -o bin/mcp-server ./cmd/server
```

### 运行

#### 1. Stdio 模式 (本地集成)

```bash
# 直接运行
./bin/mcp-server

# 或使用 go run
go run ./cmd/server
```

#### 2. HTTP/SSE 模式 (远程访问)

```bash
# 启动 HTTP 服务器
./bin/mcp-server --mode http --port 8080
```

#### 3. Docker 部署

```bash
# 使用 Docker Compose
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 🔧 配置说明

### 预设模板

服务器内置三种配置模板：

1. **严格模式** (`strict`): 最高标准，适用于关键系统
   - 圈复杂度 ≤ 5
   - 测试覆盖率 ≥ 85%
   - 启用所有检查规则

2. **标准模式** (`standard`): 平衡标准，适用于一般项目
   - 圈复杂度 ≤ 10
   - 测试覆盖率 ≥ 70%
   - 启用大部分检查规则

3. **宽松模式** (`relaxed`): 基础标准，适用于原型项目
   - 圈复杂度 ≤ 15
   - 测试覆盖率 ≥ 60%
   - 启用核心检查规则

### 自定义配置

创建 `.golangci.yml` 文件：

```yaml
linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    - gosec
    - errcheck
    
linters-settings:
  gocyclo:
    min-complexity: 10
  govet:
    check-shadowing: true
```

## 🎯 MCP 工具列表

### 1. analyze_code

分析 Go 代码并返回详细的检查结果。

```json
{
  "code": "package main\n\nfunc main() {\n    println(\"hello\")\n}",
  "standard": "standard",
  "format": "json"
}
```

### 2. manage_config

管理自定义配置文件。

```json
{
  "action": "upload",
  "name": "my-team-config",
  "content": "..."
}
```

### 3. manage_templates

管理预设配置模板。

```json
{
  "action": "list"
}
```

### 4. generate_report

生成分析报告。

```json
{
  "analysis_id": "uuid",
  "format": "markdown"
}
```

### 5. batch_analyze

批量分析多个项目。

```json
{
  "projects": [
    {"path": "/path/to/project1"},
    {"path": "/path/to/project2"}
  ]
}
```

### 6. health_check

检查服务健康状态。

```json
{}
```

## 🏗️ 项目结构

```
go-standards-mcp-server/
├── cmd/
│   └── server/           # 主程序入口
├── internal/
│   ├── mcp/             # MCP 协议实现
│   ├── analyzer/        # 代码分析引擎
│   ├── config/          # 配置管理
│   ├── report/          # 报告生成
│   └── storage/         # 数据存储
├── pkg/
│   ├── linters/         # Linter 工具集成
│   └── models/          # 数据模型
├── configs/
│   ├── templates/       # 预设模板
│   └── default.yaml     # 默认配置
├── scripts/             # 部署和工具脚本
├── docs/                # 项目文档
├── tests/               # 测试文件
├── docker-compose.yml   # Docker 编排
├── Dockerfile           # Docker 镜像
└── README.md
```

## 📊 使用示例

### Cursor IDE 集成

1. 打开 Cursor 设置
2. 添加 MCP 服务器配置：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/mcp-server",
      "args": []
    }
  }
}
```

3. 重启 Cursor，即可使用代码检查工具

### CLI 使用

```bash
# 分析单个文件
./bin/mcp-server analyze --file main.go --standard strict

# 分析整个项目
./bin/mcp-server analyze --path ./myproject --standard standard

# 使用自定义配置
./bin/mcp-server analyze --path ./myproject --config .golangci.yml

# 生成 HTML 报告
./bin/mcp-server analyze --path ./myproject --format html --output report.html
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行集成测试
go test -tags=integration ./tests/integration/...
```

## 📈 性能指标

- 单文件分析: < 5 秒
- 小型项目（< 100 文件）: < 30 秒
- 中型项目（< 1000 文件）: < 2 分钟
- 并发支持: 100+ 请求
- 内存使用: < 1GB

## 🤝 贡献指南

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 开源协议

本项目采用 MIT 协议 - 详见 [LICENSE](LICENSE) 文件

## 👥 作者

**MOONL0323**

- GitHub: [@MOONL0323](https://github.com/MOONL0323)

## 🙏 致谢

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP Go SDK
- [golangci-lint](https://github.com/golangci/golangci-lint) - Go linters aggregator
- [staticcheck](https://staticcheck.io/) - Advanced Go linter
- [gosec](https://github.com/securego/gosec) - Go security checker

## 📞 支持

- 提交 Issue: [GitHub Issues](https://github.com/MOONL0323/go-standards-mcp-server/issues)
- 讨论区: [GitHub Discussions](https://github.com/MOONL0323/go-standards-mcp-server/discussions)
- 邮箱: support@example.com

---

**项目版本**: v1.0.0  
**最后更新**: 2025-11-01
