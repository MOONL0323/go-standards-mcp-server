# 项目结构说明

本文档详细说明了 Go Standards MCP Server 项目的目录结构和文件组织。

## 目录结构

```
go-standards-mcp-server/
├── cmd/                        # 命令行应用入口
│   └── server/                 # MCP 服务器主程序
│       └── main.go            # 主入口文件
│
├── internal/                   # 内部私有包（不对外暴露）
│   ├── mcp/                   # MCP 协议实现
│   │   └── server.go          # MCP 服务器核心逻辑
│   ├── analyzer/              # 代码分析引擎
│   │   ├── analyzer.go        # 分析器主逻辑
│   │   └── analyzer_test.go   # 分析器测试
│   ├── config/                # 配置管理
│   │   └── config.go          # 配置加载和验证
│   ├── report/                # 报告生成器（待实现）
│   │   └── generator.go       
│   └── storage/               # 数据存储（待实现）
│       └── storage.go         
│
├── pkg/                        # 公共可复用包
│   ├── models/                # 数据模型
│   │   ├── types.go           # 核心类型定义
│   │   └── types_test.go      # 模型测试
│   └── linters/               # Linter 工具集成
│       ├── linter.go          # Linter 接口定义
│       ├── golangci.go        # golangci-lint 集成
│       └── govet.go           # go vet 集成
│
├── configs/                    # 配置文件
│   ├── default.yaml           # 默认配置
│   └── templates/             # 预设配置模板
│       ├── strict.yaml        # 严格模式
│       ├── standard.yaml      # 标准模式
│       └── relaxed.yaml       # 宽松模式
│
├── docs/                       # 项目文档
│   ├── API.md                 # API 文档
│   ├── CURSOR_INTEGRATION.md  # Cursor IDE 集成指南
│   └── DEPLOYMENT.md          # 部署指南
│
├── examples/                   # 示例代码
│   ├── sample.go              # 示例 Go 文件
│   └── README.md              # 示例说明
│
├── scripts/                    # 工具脚本
│   ├── quickstart.sh          # 快速启动脚本 (Unix)
│   └── quickstart.ps1         # 快速启动脚本 (Windows)
│
├── tests/                      # 测试文件（集成测试）
│   ├── integration/           # 集成测试
│   └── fixtures/              # 测试数据
│
├── .gitignore                 # Git 忽略文件
├── CHANGELOG.md               # 变更日志
├── CONTRIBUTING.md            # 贡献指南
├── Dockerfile                 # Docker 镜像构建文件
├── docker-compose.yml         # Docker Compose 配置
├── go.mod                     # Go 模块定义
├── go.sum                     # Go 模块依赖锁定
├── LICENSE                    # 开源协议
├── Makefile                   # 构建脚本
├── README.md                  # 项目说明
└── 需求.md                    # 项目需求文档
```

## 核心组件说明

### 1. cmd/server/main.go

应用程序入口，负责：
- 解析命令行参数
- 加载配置文件
- 初始化日志系统
- 启动 MCP 服务器

### 2. internal/mcp/server.go

MCP 协议服务器实现：
- 注册 MCP 工具（analyze_code, manage_config 等）
- 处理工具调用请求
- 参数验证和错误处理
- 结果格式化

### 3. internal/analyzer/analyzer.go

代码分析引擎：
- 管理多个 linter 工具
- 协调分析流程
- 聚合和统计分析结果
- 生成改进建议

### 4. internal/config/config.go

配置管理系统：
- 从文件和环境变量加载配置
- 配置验证
- 默认值处理

### 5. pkg/models/types.go

核心数据模型：
- `AnalysisRequest`: 分析请求
- `AnalysisResult`: 分析结果
- `Issue`: 代码问题
- `Summary`: 统计摘要
- 其他辅助类型

### 6. pkg/linters/

Linter 工具集成：
- `linter.go`: Linter 接口定义
- `golangci.go`: golangci-lint 集成
- `govet.go`: go vet 集成
- 可扩展支持更多工具

## 配置文件说明

### configs/default.yaml

服务器默认配置，包含：
- 服务器模式和端口
- 日志配置
- 分析器配置
- Linter 开关
- 存储配置
- 缓存配置

### configs/templates/

预设的代码规范模板：
- **strict.yaml**: 最严格标准，适用于关键系统
- **standard.yaml**: 平衡标准，适用于一般项目
- **relaxed.yaml**: 宽松标准，适用于原型开发

## 数据流

```
1. 客户端（Cursor/CLI）
   ↓
2. MCP 协议层 (internal/mcp)
   ↓
3. 分析引擎 (internal/analyzer)
   ↓
4. Linter 工具 (pkg/linters)
   ↓
5. 结果聚合和格式化
   ↓
6. 返回给客户端
```

## 扩展点

### 添加新的 Linter

1. 在 `pkg/linters/` 创建新文件（如 `staticcheck.go`）
2. 实现 `Linter` 接口
3. 在 `analyzer.go` 中注册新 linter

### 添加新的 MCP 工具

1. 在 `internal/mcp/server.go` 定义工具 schema
2. 实现工具处理函数
3. 在 `registerTools()` 中注册

### 添加新的报告格式

1. 在 `internal/report/` 实现格式生成器
2. 在 `formatResult()` 中添加格式支持

### 添加新的配置模板

1. 在 `configs/templates/` 创建新的 YAML 文件
2. 遵循 golangci-lint 配置格式
3. 在模板管理中注册

## 测试结构

```
tests/
├── integration/              # 集成测试
│   ├── analyzer_test.go     # 分析器集成测试
│   └── mcp_test.go          # MCP 协议集成测试
└── fixtures/                # 测试数据
    ├── sample1/
    │   └── main.go
    └── sample2/
        └── main.go
```

## 构建产物

```
bin/                         # 编译后的二进制文件
├── mcp-server              # Unix 可执行文件
└── mcp-server.exe          # Windows 可执行文件

reports/                     # 生成的分析报告
tmp/                        # 临时文件
data/                       # 数据库文件（SQLite）
```

## 开发工作流

1. **本地开发**
   ```bash
   make run              # 运行开发服务器
   make test             # 运行测试
   make lint             # 代码检查
   ```

2. **构建**
   ```bash
   make build            # 构建二进制文件
   make docker-build     # 构建 Docker 镜像
   ```

3. **部署**
   ```bash
   make docker-compose-up   # Docker Compose 部署
   ```

## 依赖管理

项目使用 Go Modules 管理依赖：

```bash
go mod download        # 下载依赖
go mod tidy           # 清理依赖
go mod verify         # 验证依赖
```

主要依赖：
- `github.com/mark3labs/mcp-go`: MCP 协议支持
- `github.com/spf13/viper`: 配置管理
- `go.uber.org/zap`: 日志系统
- `github.com/google/uuid`: UUID 生成

## 代码规范

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化
- 使用 `golangci-lint` 检查
- 公开 API 必须有文档注释
- 测试覆盖率 > 70%

## 版本控制

- 主分支: `main`
- 功能分支: `feature/xxx`
- 修复分支: `fix/xxx`
- 发布分支: `release/vx.x.x`

## 相关文档

- [API 文档](docs/API.md)
- [贡献指南](CONTRIBUTING.md)
- [部署指南](docs/DEPLOYMENT.md)
- [变更日志](CHANGELOG.md)

---

**注意**: 本项目持续开发中，结构可能会随版本更新而变化。
