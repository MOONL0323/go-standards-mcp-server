# 更新日志

本文件记录项目的所有重要更改。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

**[English Version](CHANGELOG.md)**

---

## [1.1.0] - 2025-11-03

### 新增
- **Git 增量检测**：仅分析变更文件
  - 三种模式：暂存、已修改、分支对比
  - 支持 pre-commit 和 pre-push 钩子
  - 自动安装钩子
  - 可配置基准分支
- **多用户架构**：支持并发用户
  - 用户工作空间隔离
  - 共享团队规范
  - 会话管理与超时
  - 用户资源限制
- **文档上传**：团队规范管理
  - 支持 PDF/TXT/Markdown
  - 自动规则提取
  - 自动生成 golangci-lint 配置
- **CLI 工具**：独立命令行界面
  - CLI 支持所有 MCP 功能
  - Git 集成
  - 适配 CI/CD
- **平台独立**：通用模块路径
  - 模块：`go-standards-mcp-server`
  - 轻松迁移到任意 Git 平台

### 新工具
- `upload_document`：上传编码规范
- `git_config`：配置 Git 检测
- `git_check`：检查是否为 Git 仓库

### 改进
- 完全重写 README
- ASCII 架构图
- 专业化项目结构
- 中英双语文档

### 配置
- `.go-standards-git.yaml`：Git 集成配置

---

## [1.0.0] - 2025-11-01

### 新增
- 首次发布
- MCP 协议支持（stdio）
- 集成 golangci-lint 和 go vet
- 三套配置模板（严格/标准/宽松）
- 详细问题报告
- JSON 和 Markdown 报告
- Docker 支持

### 功能
- `analyze_go_code`：可配置标准分析
- `list_standards`：列出编码规范
- `get_config`/`set_config`：管理配置

### 模板
- **严格**：复杂度 ≤ 5，覆盖率 ≥ 85%
- **标准**：复杂度 ≤ 10，覆盖率 ≥ 70%
- **宽松**：复杂度 ≤ 15，覆盖率 ≥ 60%

---

## 路线图

### v1.2（计划中）
- HTML/PDF 报告生成
- 批量项目分析
- 性能指标追踪
- Web 仪表板

### v2.0（未来）
- IDE 扩展（VS Code、GoLand）
- 团队协作功能
- 自定义 linter 插件
- 多语言支持

---

[1.1.0]: https://github.com/MOONL0323/go-standards-mcp-server/releases/tag/v1.1.0
[1.0.0]: https://github.com/MOONL0323/go-standards-mcp-server/releases/tag/v1.0.0
