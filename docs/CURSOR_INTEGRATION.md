# Cursor IDE 集成指南

本指南将帮助您在 Cursor IDE 中集成 Go Standards MCP Server。

## 前置条件

1. 已安装 Cursor IDE
2. 已安装 Go Standards MCP Server
3. 确认服务器可以正常运行

## 配置步骤

### 1. 构建服务器

首先构建 MCP 服务器：

```bash
cd /path/to/go-standards-mcp-server
go build -o bin/mcp-server ./cmd/server
```

### 2. 配置 Cursor

打开 Cursor 设置（Settings），找到 MCP 服务器配置部分。

#### Windows 配置

在 Cursor 配置文件中添加（通常在 `%APPDATA%\Cursor\User\settings.json`）：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "D:\\path\\to\\go-standards-mcp-server\\bin\\mcp-server.exe",
      "args": [],
      "env": {
        "MCP_SERVER_MODE": "stdio",
        "MCP_LOG_LEVEL": "info"
      }
    }
  }
}
```

#### macOS/Linux 配置

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/go-standards-mcp-server/bin/mcp-server",
      "args": [],
      "env": {
        "MCP_SERVER_MODE": "stdio",
        "MCP_LOG_LEVEL": "info"
      }
    }
  }
}
```

### 3. 重启 Cursor

配置完成后，重启 Cursor IDE 以加载新的 MCP 服务器。

## 使用方法

### 分析当前文件

1. 打开一个 Go 文件
2. 打开 Cursor 的命令面板（Ctrl/Cmd + Shift + P）
3. 输入 "MCP: Analyze Code"
4. 选择 "go-standards: analyze_code"
5. 选择分析标准（strict/standard/relaxed）
6. 查看分析结果

### 在聊天中使用

在 Cursor 的 AI 聊天窗口中，您可以直接请求代码分析：

```
@go-standards 分析当前文件，使用 standard 模式
```

或者：

```
使用 go-standards 的 strict 模式分析这段代码的质量问题
```

### 示例对话

**用户**: 
```
@go-standards 请分析以下代码：

package main

import "fmt"

func main() {
    x := 1
    fmt.Println(x)
}
```

**AI（使用 MCP 工具）**:
```json
{
  "id": "abc-123",
  "status": "success",
  "summary": {
    "total_issues": 0,
    "score": 100,
    "files_analyzed": 1
  },
  "issues": []
}
```

代码质量良好，没有发现问题！

## 高级配置

### 自定义配置文件

如果您想使用自定义配置文件，可以在启动参数中指定：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/mcp-server",
      "args": ["--config", "/path/to/custom-config.yaml"],
      "env": {
        "MCP_SERVER_MODE": "stdio"
      }
    }
  }
}
```

### 日志配置

启用详细日志以便调试：

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/mcp-server",
      "args": ["--log-level", "debug"],
      "env": {
        "MCP_SERVER_MODE": "stdio",
        "MCP_LOG_OUTPUT": "/tmp/mcp-server.log"
      }
    }
  }
}
```

## 常见问题

### Q: 服务器无法启动

**A**: 检查以下几点：
1. 确认 mcp-server 可执行文件路径正确
2. 确认文件有执行权限（Unix 系统）
3. 查看 Cursor 的开发者工具日志
4. 检查服务器日志文件

### Q: 分析很慢

**A**: 
1. 大型项目分析需要时间，请耐心等待
2. 可以先分析单个文件或小模块
3. 考虑使用宽松模式（relaxed）

### Q: 找不到 golangci-lint

**A**: 
1. 确保 golangci-lint 已安装
2. 将 golangci-lint 添加到 PATH 环境变量
3. 或在配置中禁用 golangci-lint，仅使用 go vet

### Q: 如何查看详细的问题信息

**A**: 
1. 使用 format: "markdown" 获取更易读的输出
2. 在 options 中设置 detailed: true
3. 查看生成的报告文件

## 工作流示例

### 1. 开发前检查

在开始编写代码前，检查项目结构：

```
@go-standards 使用 strict 模式分析整个项目
```

### 2. 开发中检查

编写代码时，定期检查当前文件：

```
@go-standards 快速检查当前文件有没有明显问题
```

### 3. 提交前检查

提交代码前，进行全面检查：

```
@go-standards 使用 standard 模式分析所有修改的文件，生成详细报告
```

### 4. 代码审查

审查他人代码时：

```
@go-standards 分析这个 PR 的代码质量，重点关注安全和性能问题
```

## 最佳实践

1. **逐步提升标准**: 从 relaxed → standard → strict
2. **持续集成**: 在 CI/CD 中也运行相同的检查
3. **团队一致**: 团队使用相同的配置标准
4. **及时修复**: 发现问题及时修复，不要累积
5. **学习改进**: 查看建议，学习最佳实践

## 性能优化

1. **增量分析**: 只分析修改的文件
2. **调整超时**: 根据项目大小调整超时设置
3. **并发控制**: 合理设置并发限制
4. **缓存利用**: 启用结果缓存（如果配置了 Redis）

## 故障排除

### 启用调试模式

```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/mcp-server",
      "args": ["--log-level", "debug"],
      "env": {
        "MCP_SERVER_MODE": "stdio",
        "MCP_LOG_OUTPUT": "/tmp/mcp-debug.log"
      }
    }
  }
}
```

### 查看日志

```bash
# 实时查看日志
tail -f /tmp/mcp-debug.log

# 在 Windows 上
Get-Content /tmp/mcp-debug.log -Wait
```

### 手动测试

```bash
# 手动运行服务器测试
./bin/mcp-server --mode stdio

# 然后可以发送 JSON-RPC 请求进行测试
```

## 更多资源

- [MCP 协议文档](https://modelcontextprotocol.io/)
- [项目 GitHub](https://github.com/MOONL0323/go-standards-mcp-server)
- [API 文档](../docs/API.md)
- [贡献指南](../CONTRIBUTING.md)

## 支持

如有问题，请：
1. 查看项目文档
2. 搜索 GitHub Issues
3. 提交新的 Issue
4. 加入社区讨论

---

**提示**: 配置完成后，建议先用简单的代码测试，确认服务正常工作后再用于大型项目。
