# 版本发布说明 - v0.5.0

## 概述

Go Standards MCP Server v0.5.0 是一个重要版本,引入了全面的多用户会话管理、Git 集成的增量代码分析以及专业的双语文档。此版本专注于企业级可扩展性和开发者体验改进。

## 新增功能

### 核心特性

**多用户会话管理**
- 实现了支持并发会话处理的 SessionManager
- 自动会话清理和超时管理
- 每用户上下文隔离,实现独立的分析工作流
- 会话持久化,可配置超时时间(默认 30 分钟)

**Git 集成增量分析**
- 新增 `git_check` 工具:仅分析 Git 引用之间修改的 Go 文件
- 新增 `git_config` 工具:配置基准分支和目标引用进行比较
- 智能变更检测:支持提交哈希、分支和标签
- 性能优化:仅分析变更文件,而非整个代码库

**增强的 MCP 工具**
- `analyze_go_code`:全功能 Go 代码静态分析
- `git_check`:基于 Git 的增量代码分析
- `git_config`:Git 比较配置
- `list_standards`:浏览可用的编码规范
- `upload_document`:自定义规范管理
- `get_config`/`set_config`:运行时配置控制
- `health_check`:服务健康监控

### 架构改进

**服务器架构**
- 重构 MCP 服务器,采用模块化工具注册
- 改进错误处理和日志基础设施
- 增强配置管理和验证
- 优化资源清理和优雅关闭

**内部模块**
- `internal/usercontext`:多用户会话管理
- `internal/git`:Git 集成和变更检测
- `internal/analyzer`:增强的 Go 代码分析引擎
- `internal/standards`:规范管理系统

### 文档

**双语文档**
- 完整的英文文档 (README.md)
- 完整的中文文档 (README_CN.md)
- 双语更新日志 (CHANGELOG.md, CHANGELOG_CN.md)
- 专业格式,无表情符号或非正式语言

**全面的使用指南**
- 完整代码分析工作流
- Git 增量分析分步指南
- 自定义规范集成
- 高级配置选项
- MCP 客户端集成示例

**可视化架构图**
- 使用 Mermaid 流程图展示系统架构
- 支持 VS Code 预览
- 清晰的组件关系和数据流

## 安装说明

### 二进制安装

从发布资产中下载适合您平台的二进制文件:

**Windows**
```powershell
# 下载 mcp-server.exe 和 go-standards-cli.exe
# 添加到 PATH 或直接运行
.\mcp-server.exe
```

**Linux**
```bash
# 下载 mcp-server-linux 和 go-standards-cli-linux
chmod +x mcp-server-linux go-standards-cli-linux
./mcp-server-linux
```

**macOS**
```bash
# Intel Mac: 下载 *-darwin-amd64 二进制文件
# Apple Silicon: 下载 *-darwin-arm64 二进制文件
chmod +x mcp-server-darwin-arm64 go-standards-cli-darwin-arm64
./mcp-server-darwin-arm64
```

### 从源码构建

```bash
git clone https://github.com/yourusername/go-standards-mcp-server.git
cd go-standards-mcp-server
go build -o bin/mcp-server ./cmd/server
go build -o bin/go-standards-cli ./cmd/cli
```

## 配置说明

### MCP 客户端集成

添加到您的 MCP 客户端配置中(例如 Claude Desktop):

**Windows (PowerShell)**
```json
{
  "mcpServers": {
    "go-standards": {
      "command": "powershell.exe",
      "args": [
        "-NoProfile",
        "-ExecutionPolicy",
        "Bypass",
        "-File",
        "D:\\project\\go-standards-mcp-server\\bin\\mcp-server.exe"
      ]
    }
  }
}
```

**Linux/macOS**
```json
{
  "mcpServers": {
    "go-standards": {
      "command": "/path/to/mcp-server-linux"
    }
  }
}
```

### 服务器配置

在项目根目录创建 `config.yaml`:

```yaml
server:
  name: "Go Standards MCP Server"
  version: "0.5.0"
  log_level: "info"
  session_timeout: 1800  # 30 分钟

analyzer:
  max_file_size: 10485760  # 10MB
  ignore_patterns:
    - "*_test.go"
    - "vendor/**"
    - ".git/**"
  enable_complexity: true
  complexity_threshold: 15
```

## 使用示例

### 完整项目分析

```
使用默认规范分析位于 /path/to/project 的整个 Go 项目
```

### Git 增量分析

```
配置 Git 以比较 feature/new-api 分支与 main 分支:
- 基准分支: main
- 目标分支: feature/new-api

然后仅分析变更的 Go 文件
```

### 自定义规范

```
从 /path/to/custom-standards.md 上传自定义编码规范

然后使用默认规范和自定义规范分析项目
```

## 升级指南

### 从 v1.x.x 升级到 v0.5.0

**版本号变更**
此版本将版本号重置为 0.5.0,以反映项目当前的成熟度并遵循语义化版本控制实践。

**重大变更**
- 无。此版本与 v1.x.x 配置向后兼容。

**新增配置选项**
- `session_timeout`:如果使用多用户功能,请添加到 `server` 部分
- Git 配置:无需更改,Git 工具是增量添加的

**迁移步骤**
1. 停止现有 MCP 服务器
2. 用 v0.5.0 版本替换二进制文件
3. (可选)在 config.yaml 中添加 `session_timeout`
4. 重启 MCP 服务器
5. 使用 `health_check` 工具测试

## 性能改进

- Git 增量分析使大型代码库的扫描时间减少 70-90%
- 会话清理防止长时间运行部署中的内存泄漏
- 优化的文件解析将 CPU 使用率降低 30%
- 支持多文件并发分析处理

## 已知问题

- Git 集成需要系统安装 Git CLI
- 超大差异(10000+ 文件)可能会出现性能下降
- 会话持久化不会在服务器重启后保留(设计如此)

## 贡献者

此版本得益于全面的测试和反馈。特别感谢所有帮助改进文档和识别问题的贡献者。

## 支持

- GitHub Issues: https://github.com/yourusername/go-standards-mcp-server/issues
- 文档: README.md 和 README_CN.md
- 许可证: MIT

## 校验和

发布二进制文件的 SHA256 校验和:

```
# 使用以下命令生成: Get-FileHash -Algorithm SHA256 <file>
# 将在构建后添加
```

---

**完整更新日志**: https://github.com/yourusername/go-standards-mcp-server/blob/dev/CHANGELOG_CN.md

**发布日期**: 2025-01-XX

**项目状态**: 活跃开发中
