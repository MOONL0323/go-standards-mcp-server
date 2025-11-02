# Contributing to Go Standards MCP Server

感谢您对本项目的关注！我们欢迎各种形式的贡献。

## 贡献方式

### 报告问题

如果您发现 bug 或有功能建议，请在 GitHub Issues 中创建问题：

1. 搜索现有 Issues，避免重复
2. 使用清晰的标题和描述
3. 提供重现步骤（针对 bug）
4. 附上相关日志和配置信息

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/MOONL0323/go-standards-mcp-server.git
   cd go-standards-mcp-server
   ```

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **编写代码**
   - 遵循 Go 代码规范
   - 添加必要的测试
   - 更新相关文档

4. **运行测试**
   ```bash
   make test
   make lint
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

6. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建 Pull Request**
   - 提供清晰的 PR 描述
   - 关联相关 Issues
   - 等待代码审查

## 代码规范

### Go 代码风格

- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 检查代码质量
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型（type）：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(analyzer): add support for custom linters

Add ability to integrate custom linting tools through plugin system.

Closes #123
```

### 测试要求

- 新功能必须包含单元测试
- 测试覆盖率应保持在 70% 以上
- 关键路径需要集成测试
- 所有测试必须通过才能合并

### 文档要求

- 公开函数和类型必须有注释
- 复杂逻辑需要说明注释
- 新功能需要更新 README
- API 变更需要更新文档

## 开发环境设置

### 环境要求

- Go 1.21+
- golangci-lint
- Docker (可选)
- Make

### 安装依赖

```bash
make install
make setup-dev
```

### 运行开发服务器

```bash
make run
```

### 构建项目

```bash
make build
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定测试
go test -v ./internal/analyzer/...

# 生成覆盖率报告
make test
```

### 代码检查

```bash
make lint
make fmt
```

## 发布流程

项目维护者会定期发布新版本：

1. 更新版本号
2. 更新 CHANGELOG
3. 创建 Git tag
4. 构建 Docker 镜像
5. 发布 GitHub Release

## 社区准则

### 行为规范

- 尊重所有贡献者
- 保持友善和专业
- 接受建设性批评
- 关注最佳实践

### 沟通渠道

- GitHub Issues：问题报告和功能讨论
- GitHub Discussions：一般性讨论
- Pull Requests：代码审查

## 获取帮助

如果您有任何问题：

1. 查看项目文档
2. 搜索现有 Issues
3. 在 Discussions 中提问
4. 发送邮件至维护者

## 许可证

贡献的代码将遵循项目的 MIT 许可证。

---

感谢您的贡献！
