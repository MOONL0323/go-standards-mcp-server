# 快速开始指南

本指南将帮助您在 5 分钟内启动并运行 Go Standards MCP Server。

## 前置要求检查

在开始之前，请确保已安装：

- ✅ Go 1.21 或更高版本
- ✅ Git
- ⚠️ golangci-lint（推荐，但可选）
- ⚠️ Docker（可选，用于容器化部署）

## 快速启动步骤

### 方式 1: 使用快速启动脚本（推荐）

#### Windows (PowerShell)
```powershell
.\scripts\quickstart.ps1
```

#### Linux/macOS
```bash
chmod +x scripts/quickstart.sh
./scripts/quickstart.sh
```

### 方式 2: 手动安装

#### 1. 克隆项目

```bash
git clone https://github.com/MOONL0323/go-standards-mcp-server.git
cd go-standards-mcp-server
```

#### 2. 安装依赖

```bash
go mod download
```

#### 3. 构建服务器

```bash
# 使用 Make（推荐）
make build

# 或者直接使用 Go
go build -o bin/mcp-server ./cmd/server
```

#### 4. 运行服务器

```bash
# stdio 模式（用于 MCP 集成）
./bin/mcp-server

# HTTP 模式（用于远程访问）
./bin/mcp-server --mode http --port 8080
```

### 方式 3: 使用 Docker

```bash
# 使用 Docker Compose
docker-compose up -d

# 查看日志
docker-compose logs -f mcp-server
```

## 验证安装

### 1. 检查服务器运行状态

如果使用 HTTP 模式：

```bash
curl http://localhost:8080/health
```

预期输出：
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "checks": {
    "analyzer": "ok",
    "config": "ok"
  }
}
```

### 2. 测试代码分析

创建一个测试文件 `test.go`：

```go
package main

func main() {
    println("hello")
}
```

使用 MCP 工具分析（通过 Cursor 或其他 MCP 客户端）：

```json
{
  "code": "package main\n\nfunc main() {\n    println(\"hello\")\n}",
  "standard": "standard",
  "format": "json"
}
```

## 在 Cursor IDE 中使用

### 1. 配置 Cursor

编辑 Cursor 配置文件（通常在用户设置中）：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/bin/mcp-server",
      "args": []
    }
  }
}
```

### 2. 重启 Cursor

保存配置后，重启 Cursor IDE。

### 3. 使用工具

在 Cursor 的 AI 聊天中：

```
@go-standards 分析这段代码的质量
```

## 常见使用场景

### 场景 1: 分析代码片段

```json
{
  "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tx := 42\n\tfmt.Println(\"hello\")\n}",
  "standard": "standard",
  "format": "markdown"
}
```

### 场景 2: 分析整个项目

```json
{
  "project_dir": "/path/to/your/go/project",
  "standard": "strict",
  "format": "json"
}
```

### 场景 3: 使用自定义配置

```json
{
  "project_dir": "/path/to/project",
  "standard": "custom",
  "config": "linters:\n  enable:\n    - gofmt\n    - govet\n    - staticcheck",
  "format": "json"
}
```

### 场景 4: 查看可用模板

```json
{
  "action": "list"
}
```

## 配置说明

### 三种预设模板

1. **strict（严格模式）**
   - 圈复杂度 ≤ 5
   - 测试覆盖率 ≥ 85%
   - 启用所有检查规则
   - 适用于：关键系统、金融应用、安全敏感项目

2. **standard（标准模式）** ⭐ 推荐
   - 圈复杂度 ≤ 10
   - 测试覆盖率 ≥ 70%
   - 启用大部分检查规则
   - 适用于：大多数生产项目、团队协作

3. **relaxed（宽松模式）**
   - 圈复杂度 ≤ 15
   - 测试覆盖率 ≥ 60%
   - 启用核心检查规则
   - 适用于：原型开发、学习项目、快速实验

### 自定义配置

在项目根目录创建 `.golangci.yml` 文件：

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

issues:
  exclude-use-default: true
```

## 下一步

### 深入学习

- 📖 [完整 API 文档](docs/API.md)
- 🔧 [配置详解](configs/default.yaml)
- 🚀 [部署指南](docs/DEPLOYMENT.md)
- 💻 [Cursor 集成](docs/CURSOR_INTEGRATION.md)

### 高级用法

- 批量分析多个项目
- 生成 HTML/PDF 报告
- 集成到 CI/CD 流程
- 自定义 Linter 规则

### 最佳实践

1. **从 standard 开始**: 先使用标准模式熟悉工具
2. **逐步提升**: 代码质量稳定后升级到 strict
3. **团队一致**: 团队使用相同的配置标准
4. **持续改进**: 定期运行检查，及时修复问题
5. **学习建议**: 关注工具给出的改进建议

## 常见问题

### Q: golangci-lint 未找到

**A**: 安装 golangci-lint：
```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows
# 从 GitHub releases 下载预编译二进制文件
```

### Q: 分析速度慢

**A**: 
- 使用 `relaxed` 模式加快分析
- 分析单个文件而不是整个项目
- 增加并发限制配置

### Q: 如何查看日志

**A**:
```bash
# 启用调试日志
./bin/mcp-server --log-level debug

# 输出到文件
./bin/mcp-server --log-level info 2>&1 | tee mcp-server.log
```

### Q: Docker 容器无法启动

**A**:
- 检查端口 8080 是否被占用
- 查看容器日志: `docker-compose logs mcp-server`
- 确认配置文件正确挂载

## 获取帮助

遇到问题？

1. 📚 查看 [文档](docs/)
2. 🐛 提交 [Issue](https://github.com/MOONL0323/go-standards-mcp-server/issues)
3. 💬 加入 [讨论](https://github.com/MOONL0323/go-standards-mcp-server/discussions)
4. 📧 发送邮件至维护者

## 贡献

欢迎贡献！查看 [贡献指南](CONTRIBUTING.md) 了解如何参与。

---

**🎉 恭喜！您已经成功启动了 Go Standards MCP Server！**

现在开始使用它来提升您的 Go 代码质量吧！
