# 使用演示

本文档提供了 Go Standards MCP Server 的实际使用示例。

## 场景 1: 分析简单代码片段

### 输入

通过 MCP 工具 `analyze_code` 分析以下代码：

```go
package main

func main() {
    x := 42
    println("hello")
}
```

### MCP 请求

```json
{
  "code": "package main\n\nfunc main() {\n    x := 42\n    println(\"hello\")\n}",
  "standard": "standard",
  "format": "json"
}
```

### 预期输出

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "success",
  "issues": [
    {
      "file": "main.go",
      "line": 4,
      "column": 5,
      "severity": "warning",
      "category": "dead-code",
      "rule": "unused",
      "message": "declared and not used: x",
      "source": "govet"
    },
    {
      "file": "main.go",
      "line": 5,
      "column": 5,
      "severity": "warning",
      "category": "style",
      "rule": "println",
      "message": "use of println not recommended",
      "source": "govet"
    }
  ],
  "summary": {
    "total_issues": 2,
    "error_count": 0,
    "warning_count": 2,
    "info_count": 0,
    "files_analyzed": 1,
    "lines_analyzed": 6,
    "duration": "1.2s",
    "score": 96.0,
    "category_counts": {
      "dead-code": 1,
      "style": 1
    }
  },
  "suggestions": [],
  "metadata": {
    "standard": "standard",
    "tools_used": ["golangci-lint", "govet"],
    "go_version": "1.21",
    "server_version": "1.0.0"
  },
  "created_at": "2025-11-01T12:00:00Z"
}
```

---

## 场景 2: 分析项目目录

### 项目结构

```
myproject/
├── main.go
├── handler.go
└── utils.go
```

### MCP 请求

```json
{
  "project_dir": "/path/to/myproject",
  "standard": "strict",
  "format": "markdown"
}
```

### 预期输出（Markdown）

```markdown
# Code Analysis Report

**Status**: success
**Score**: 78.5/100

## Summary

- Total Issues: 15
- Errors: 2
- Warnings: 10
- Files Analyzed: 3
- Duration: 3.5s

## Issues

### 1. Cyclomatic complexity too high
- **File**: handler.go:25:1
- **Severity**: error
- **Category**: complexity
- **Rule**: gocyclo

### 2. Error return value not checked
- **File**: main.go:18:5
- **Severity**: warning
- **Category**: error-handling
- **Rule**: errcheck

### 3. Function should have comment
- **File**: utils.go:10:1
- **Severity**: warning
- **Category**: style
- **Rule**: revive

... and 12 more issues

## Suggestions

### 1. Reduce Function Complexity
**Priority**: high
**Category**: complexity

Several functions exceed the complexity threshold. Consider breaking them into smaller functions.

### 2. Improve Error Handling
**Priority**: high
**Category**: error-handling

Multiple instances of unchecked error returns found. Always check errors.
```

---

## 场景 3: 使用自定义配置

### 自定义配置文件

```yaml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck

linters-settings:
  gocyclo:
    min-complexity: 8
  
  errcheck:
    check-type-assertions: true
    check-blank: true

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

### MCP 请求

```json
{
  "project_dir": "/path/to/myproject",
  "standard": "custom",
  "config": "linters:\n  enable:\n    - gofmt\n    - govet\n    - errcheck\nlinters-settings:\n  gocyclo:\n    min-complexity: 8",
  "format": "json"
}
```

---

## 场景 4: 在 Cursor 中使用

### 步骤 1: 配置 Cursor

在 Cursor 设置中添加：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "D:\\project\\go-standards-mcp-server\\bin\\mcp-server.exe"
    }
  }
}
```

### 步骤 2: 在聊天中使用

**用户提问**:
```
@go-standards 请帮我分析这个 Go 文件的代码质量，使用 standard 模式

package main

import "fmt"

func main() {
    x := 42
    y := 10
    
    result := x + y
    fmt.Println(result)
}
```

**AI 响应**:
```
我已经使用 Go Standards MCP Server 分析了您的代码。以下是分析结果：

**质量评分**: 95.0/100

**发现的问题**:
1. 变量 x 和 y 命名过于简单，建议使用更有意义的名称
2. 没有错误处理

**建议**:
- 使用更具描述性的变量名
- 添加注释说明代码功能

总体来说，代码质量良好，只需要进行一些小的改进。
```

---

## 场景 5: 查看可用模板

### MCP 请求

```json
{
  "action": "list"
}
```

### 输出

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

---

## 场景 6: 健康检查

### MCP 请求

```json
{}
```

### 输出

```json
{
  "status": "healthy",
  "timestamp": "2025-11-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": "24h30m15s",
  "checks": {
    "analyzer": "ok",
    "config": "ok"
  }
}
```

---

## 场景 7: 比较不同模式

### Relaxed Mode

```json
{
  "code": "package main\n\nfunc complexFunc() {\n    // 15 行复杂逻辑\n}",
  "standard": "relaxed",
  "format": "json"
}
```

**结果**: 0 个问题，分数 100

### Standard Mode

```json
{
  "code": "package main\n\nfunc complexFunc() {\n    // 15 行复杂逻辑\n}",
  "standard": "standard",
  "format": "json"
}
```

**结果**: 1 个警告（复杂度稍高），分数 98

### Strict Mode

```json
{
  "code": "package main\n\nfunc complexFunc() {\n    // 15 行复杂逻辑\n}",
  "standard": "strict",
  "format": "json"
}
```

**结果**: 2 个错误（复杂度过高，缺少注释），分数 90

---

## 实用技巧

### 技巧 1: 逐步提升标准

```bash
# 第一周：使用 relaxed 模式熟悉工具
analyze --standard relaxed

# 第二周：升级到 standard 模式
analyze --standard standard

# 第三周：目标 strict 模式
analyze --standard strict
```

### 技巧 2: 聚焦特定类型问题

在自定义配置中只启用关注的 linter：

```yaml
linters:
  enable:
    - gosec  # 只检查安全问题
```

### 技巧 3: 排除测试文件

```yaml
issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
        - gosec
```

### 技巧 4: 增量分析

只分析修改的文件（需要 Git）：

```bash
# 获取修改的文件列表
git diff --name-only | grep '\.go$'

# 逐个分析
for file in $(git diff --name-only | grep '\.go$'); do
  analyze --file $file
done
```

---

## CI/CD 集成示例

### GitHub Actions

```yaml
name: Go Code Quality

on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run Go Standards MCP Server
        run: |
          # 安装服务器
          go install github.com/MOONL0323/go-standards-mcp-server/cmd/server@latest
          
          # 运行分析（使用 CLI 模式，待实现）
          mcp-server analyze --path . --standard standard --format json > report.json
          
      - name: Upload Report
        uses: actions/upload-artifact@v2
        with:
          name: code-quality-report
          path: report.json
```

---

## 常见问题解答

### Q: 如何只检查安全问题？

**A**: 使用自定义配置，只启用 gosec：

```json
{
  "standard": "custom",
  "config": "linters:\n  enable:\n    - gosec"
}
```

### Q: 如何调整复杂度阈值？

**A**: 在自定义配置中设置：

```yaml
linters-settings:
  gocyclo:
    min-complexity: 15  # 调整为 15
```

### Q: 如何忽略特定文件？

**A**: 在配置中添加排除规则：

```yaml
run:
  skip-files:
    - legacy/.*\.go$
    - vendor/.*
```

---

## 更多示例

查看 `examples/` 目录获取更多实际代码示例。

## 反馈和建议

如果您有其他使用场景或建议，欢迎：
- 提交 Issue
- 贡献文档
- 分享使用经验

---

**提示**: 所有示例都可以在实际环境中直接使用。
