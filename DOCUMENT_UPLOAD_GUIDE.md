# 团队代码规范文档上传指南

## 功能概述

支持上传团队的代码规范文档（PDF、TXT、Markdown 等格式），自动转换为 golangci-lint 配置，实现团队自定义代码检查标准。

## 核心特性

- **多格式支持**: PDF、TXT、Markdown
- **AI 智能转换**: 自动解析文档并生成配置
- **去重检测**: 防止重复上传相同文档
- **版本管理**: 记录所有上传历史和元数据
- **可信度评分**: AI 转换结果包含置信度评分

## 系统架构

```
团队文档 → 文档解析器 → AI转换器 → golangci-lint配置 → 代码检查
    ↓           ↓           ↓              ↓
  存储系统    元数据管理   版本控制      使用统计
```

## 使用方式

### 方式一：MCP 工具（推荐）

通过 Cursor IDE 集成使用：

```javascript
// 1. 上传文档（文本格式）
@go-standards upload_document {
  content: "团队代码规范\n\n1. 函数复杂度不超过10\n2. 必须检查所有错误...",
  file_name: "team-standard.md",
  name: "team-standard-v1",
  description: "团队代码规范 v1.0 - 2025"
}

// 2. 上传 PDF（Base64 编码）
@go-standards upload_document {
  content: "JVBERi0xLjcKCjEgMCBvYmo...",  // Base64 encoded PDF
  file_name: "coding-rules.pdf",
  name: "coding-rules-2025",
  description: "2025年度编码规范"
}

// 3. 查看所有上传的文档
@go-standards list_documents

// 4. 获取文档详情和生成的配置
@go-standards get_document { id: "a1b2c3d4" }

// 5. 使用上传的规范进行代码检查
@go-standards analyze_code {
  project_dir: "./myproject",
  standard: "custom",
  config: "team-standard-v1"
}

// 6. 删除文档
@go-standards delete_document { id: "a1b2c3d4" }
```

### 方式二：直接 API 集成

通过文档服务 API 使用：

```go
import (
    "github.com/MOONL0323/go-standards-mcp-server/internal/service"
)

// 初始化服务
docService, _ := service.NewDocumentService(logger)

// 上传文档
response, err := docService.UploadDocument(ctx, &service.UploadDocumentRequest{
    File:        file,
    FileName:    "team-standard.pdf",
    Name:        "team-standard-v1",
    Description: "团队代码规范",
})

// 列出所有文档
documents, _ := docService.ListDocuments()

// 获取生成的配置
config, _ := docService.GetDocumentConfig(documentID)
```

## 文档格式要求

### 支持的格式

| 格式 | 扩展名 | 支持状态 | 说明 |
|------|--------|---------|------|
| PDF | `.pdf` | ✅ 完全支持 | 自动提取文本 |
| 纯文本 | `.txt` | ✅ 完全支持 | 直接读取 |
| Markdown | `.md`, `.markdown` | ✅ 完全支持 | 直接读取 |
| DOCX | `.docx` | 🔜 计划支持 | 需要额外库 |
| DOC | `.doc` | ❌ 需转换 | 请先转为 PDF/DOCX |

### 文档内容建议

为获得最佳转换效果，文档应包含：

1. **明确的规则定义**
   ```
   ✅ 好的示例：
   "函数的圈复杂度不得超过 10"
   "所有错误必须进行检查和处理"
   
   ❌ 避免模糊：
   "函数不能太复杂"
   "注意错误处理"
   ```

2. **具体的数值指标**
   ```
   ✅ 包含数值：
   - 最大行数：80
   - 测试覆盖率：≥ 70%
   - 复杂度阈值：10
   
   ❌ 缺少具体指标：
   - 代码要简洁
   - 有足够的测试
   ```

3. **清晰的分类**
   ```
   建议按以下分类组织：
   - 代码质量规范
   - 命名规范
   - 注释规范
   - 错误处理规范
   - 性能规范
   - 安全规范
   ```

## AI 转换机制

### 转换流程

```
1. 文档解析
   ↓
2. 文本提取
   ↓
3. AI 分析（或模板匹配）
   ↓
4. 生成 golangci-lint 配置
   ↓
5. 置信度评估
   ↓
6. 保存配置
```

### 配置 AI 服务（可选）

使用 OpenAI 或兼容的 AI 服务提升转换质量：

```bash
# OpenAI API
export OPENAI_API_KEY="sk-..."
export AI_MODEL="gpt-4"

# 自定义 AI 服务
export AI_API_URL="https://your-ai-service.com/v1/chat/completions"
export AI_API_KEY="your-api-key"
export AI_MODEL="your-model"
```

未配置 AI 时，系统将使用基于关键字的模板转换（置信度较低）。

### 置信度评分

| 置信度 | 说明 | 建议 |
|--------|------|------|
| 0.9-1.0 | 优秀 | 可直接使用 |
| 0.7-0.9 | 良好 | 建议人工审核 |
| 0.5-0.7 | 中等 | 需要人工调整 |
| < 0.5 | 较低 | 建议重新编写配置 |

## 存储管理

### 目录结构

```
data/documents/
├── documents/           # 原始文档
│   ├── a1b2c3d4.pdf
│   └── e5f6g7h8.md
├── metadata/           # 文档元数据
│   ├── a1b2c3d4.json
│   └── e5f6g7h8.json
└── configs/            # 生成的配置
    ├── a1b2c3d4.yaml
    └── e5f6g7h8.yaml
```

### 元数据结构

```json
{
  "id": "a1b2c3d4",
  "name": "team-standard-v1",
  "original_file": "team-standard.pdf",
  "file_type": "pdf",
  "file_size": 153600,
  "uploaded_at": "2025-11-03T00:00:00Z",
  "config_name": "team-standard-v1",
  "description": "团队代码规范 v1.0",
  "version": 1,
  "hash": "a1b2c3d4e5f6...",
  "conversion_summary": "Successfully extracted 15 rules",
  "extracted_rules": [
    "gocyclo",
    "errcheck",
    "govet",
    ...
  ],
  "confidence": 0.92
}
```

## 重复检测

系统使用 SHA256 哈希算法检测重复文档：

```
相同文件 → 相同哈希 → 拒绝上传 → 提示已存在
```

如需更新规范，请：
1. 删除旧文档
2. 上传新版本

或直接使用 `manage_config` 工具更新配置。

## 使用示例

### 示例 1：上传 Markdown 规范

```bash
# 创建规范文档
cat > team-rules.md << 'EOF'
# 团队 Go 代码规范

## 代码复杂度
- 圈复杂度不超过 10
- 认知复杂度不超过 15
- 函数长度不超过 80 行

## 错误处理
- 必须检查所有错误返回值
- 禁止忽略错误
- 类型断言必须检查

## 代码风格
- 使用 gofmt 格式化
- 导入路径使用 goimports 排序
- 遵循 Go 命名规范

## 安全要求
- 使用 gosec 检查安全问题
- 中等及以上严重级别必须修复
EOF

# 通过 Cursor 上传
@go-standards upload_document {
  content: "$(cat team-rules.md)",
  file_name: "team-rules.md",
  name: "team-2025",
  description: "2025年团队规范"
}
```

### 示例 2：批量管理

```javascript
// 列出所有文档
const docs = await callTool('list_documents');
console.log(`Total documents: ${docs.length}`);

// 检查每个文档的置信度
for (const doc of docs) {
  if (doc.confidence < 0.7) {
    console.log(`低置信度文档: ${doc.name} (${doc.confidence})`);
    // 获取详情进行人工审核
    const detail = await callTool('get_document', { id: doc.id });
  }
}
```

### 示例 3：版本升级

```javascript
// 1. 列出现有版本
@go-standards list_documents

// 2. 删除旧版本
@go-standards delete_document { id: "old-version-id" }

// 3. 上传新版本
@go-standards upload_document {
  content: "...",
  file_name: "team-standard-v2.pdf",
  name: "team-standard-v2",
  description: "团队规范 v2.0 - 增加性能检查"
}
```

## 故障排除

### 常见问题

**Q: 上传失败提示"document already exists"**

A: 文档哈希相同，说明已上传过该文件。可以：
- 使用 `list_documents` 查找现有文档
- 使用 `delete_document` 删除旧文档后重新上传
- 或直接使用现有配置

**Q: 转换置信度低**

A: 可能原因：
- 未配置 AI API（使用的是模板匹配）
- 文档内容不够结构化
- 规则描述过于模糊

解决方法：
- 配置 OpenAI API Key
- 重写文档，使用更明确的规则表述
- 人工编辑生成的配置文件

**Q: PDF 解析失败**

A: 检查：
- PDF 是否为扫描版（需要 OCR）
- PDF 是否加密或有权限限制
- 尝试另存为新的 PDF
- 或转换为 TXT/Markdown 格式

**Q: 如何查看生成的配置**

A: 使用 `get_document` 工具：
```javascript
@go-standards get_document { id: "document-id" }
```

返回结果包含完整的 golangci-lint 配置。

## 最佳实践

1. **文档命名规范**
   - 使用版本号：`team-standard-v1.0.md`
   - 包含日期：`coding-rules-2025-Q1.pdf`
   - 清晰描述：`backend-api-standards.md`

2. **定期审核**
   - 每季度检查文档置信度
   - 人工审核低置信度配置
   - 根据实际使用反馈优化规则

3. **版本控制**
   - 保留历史版本记录
   - 使用描述字段记录更新内容
   - 重大更新时创建新配置名称

4. **团队协作**
   - 规范文档应由团队共同制定
   - 配置更新前进行代码库测试
   - 记录规则变更原因和影响范围

## 技术实现

### 核心模块

1. **document_parser**: 文档解析
   - 支持多种格式
   - 提取纯文本内容
   - 保留文档结构信息

2. **ai_converter**: AI 转换
   - 调用 OpenAI API
   - 规则提取和映射
   - 生成标准化配置

3. **document_storage**: 存储管理
   - 文件存储
   - 元数据管理
   - 重复检测
   - 版本控制

4. **document_service**: 业务编排
   - 完整工作流
   - 错误处理
   - 事务管理

## 未来计划

- [ ] 支持 DOCX 格式
- [ ] OCR 支持扫描版 PDF
- [ ] 配置对比和合并工具
- [ ] Web UI 管理界面
- [ ] 配置效果分析报告
- [ ] 多语言文档支持

---

**提示**: 这是一个专业级的代码规范管理系统，建议团队 Tech Lead 负责文档上传和配置管理。

