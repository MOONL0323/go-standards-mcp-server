package converter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AIConverter converts document content to golangci-lint configuration using AI
type AIConverter struct {
	apiKey   string       // OpenAI API Key 或其他 AI 服务的 Key
	apiURL   string       // API 地址
	model    string       // 模型名称
	client   *http.Client // HTTP 客户端
}

// NewAIConverter creates a new AI converter
func NewAIConverter() *AIConverter {
	// 从环境变量读取配置
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("AI_API_KEY") // 支持通用的 AI API KEY
	}

	apiURL := os.Getenv("AI_API_URL")
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/chat/completions" // 默认使用 OpenAI
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4" // 默认模型
	}

	return &AIConverter{
		apiKey: apiKey,
		apiURL: apiURL,
		model:  model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ConversionResult represents the result of document conversion
type ConversionResult struct {
	Config      string   // 生成的 golangci-lint 配置
	Summary     string   // 配置摘要说明
	Rules       []string // 提取的主要规则列表
	Confidence  float64  // 转换置信度 (0-1)
	Suggestions []string // 改进建议
}

// Convert converts document content to golangci-lint configuration
func (c *AIConverter) Convert(ctx context.Context, documentContent, fileName string) (*ConversionResult, error) {
	if c.apiKey == "" {
		// 如果没有配置 API Key，使用模板生成
		return c.convertWithTemplate(documentContent)
	}

	// 使用 AI 进行转换
	return c.convertWithAI(ctx, documentContent, fileName)
}

// convertWithAI uses AI to convert document to configuration
func (c *AIConverter) convertWithAI(ctx context.Context, content, fileName string) (*ConversionResult, error) {
	prompt := c.buildPrompt(content, fileName)

	// 构建请求
	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert in Go code quality standards and golangci-lint configuration. Convert code standard documents into valid golangci-lint YAML configurations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3, // 降低随机性，提高准确性
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// 解析 AI 返回的结果
	return c.parseAIResponse(response.Choices[0].Message.Content)
}

// buildPrompt builds the conversion prompt
func (c *AIConverter) buildPrompt(content, fileName string) string {
	return fmt.Sprintf(`Convert the following code standard document to a golangci-lint YAML configuration.

Document: %s
Content:
---
%s
---

Requirements:
1. Extract all code quality rules and convert to appropriate golangci-lint linters and settings
2. Map complexity limits, naming conventions, error handling rules, etc.
3. Generate a complete, valid YAML configuration
4. Include explanatory comments in the YAML
5. Provide a summary of the main rules extracted
6. Rate your confidence in the conversion (0-1)

Return the result in this JSON format:
{
  "config": "... complete YAML configuration ...",
  "summary": "Brief summary of the configuration",
  "rules": ["rule1", "rule2", ...],
  "confidence": 0.95,
  "suggestions": ["suggestion1", "suggestion2", ...]
}
`, fileName, truncateContent(content, 4000))
}

// parseAIResponse parses AI response into ConversionResult
func (c *AIConverter) parseAIResponse(content string) (*ConversionResult, error) {
	// 尝试从 markdown 代码块中提取 JSON
	content = extractJSONFromMarkdown(content)

	var result ConversionResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果解析失败，尝试直接将内容作为配置
		return &ConversionResult{
			Config:     content,
			Summary:    "AI generated configuration",
			Confidence: 0.7,
		}, nil
	}

	return &result, nil
}

// convertWithTemplate uses template-based conversion (fallback when no AI)
func (c *AIConverter) convertWithTemplate(content string) (*ConversionResult, error) {
	// 基于关键字的简单规则提取
	rules := c.extractRulesFromContent(content)
	
	config := c.generateTemplateConfig(rules)
	
	return &ConversionResult{
		Config:      config,
		Summary:     fmt.Sprintf("Template-based configuration with %d rules detected", len(rules)),
		Rules:       rules,
		Confidence:  0.6, // 模板方法置信度较低
		Suggestions: []string{
			"Set OPENAI_API_KEY environment variable for AI-powered conversion",
			"Review and adjust the generated configuration manually",
			"Test the configuration on your codebase",
		},
	}, nil
}

// extractRulesFromContent extracts rules using keyword matching
func (c *AIConverter) extractRulesFromContent(content string) []string {
	var rules []string
	
	keywords := map[string]string{
		"complexity":        "gocyclo",
		"error handling":    "errcheck",
		"unused":            "unused",
		"format":            "gofmt",
		"import":            "goimports",
		"naming":            "revive",
		"security":          "gosec",
		"performance":       "prealloc",
		"shadow":            "shadow",
		"type assertion":    "errcheck",
		"function length":   "funlen",
		"cognitive complexity": "gocognit",
	}
	
	contentLower := content
	for keyword, linter := range keywords {
		if contains(contentLower, keyword) {
			rules = append(rules, linter)
		}
	}
	
	if len(rules) == 0 {
		// 返回默认规则
		rules = []string{"govet", "errcheck", "staticcheck", "unused", "gofmt"}
	}
	
	return deduplicate(rules)
}

// generateTemplateConfig generates a template configuration based on detected rules
func (c *AIConverter) generateTemplateConfig(rules []string) string {
	config := `# Auto-generated configuration from team document
# Generated at: ` + time.Now().Format(time.RFC3339) + `

linters:
  enable:
`
	
	for _, rule := range rules {
		config += fmt.Sprintf("    - %s\n", rule)
	}
	
	config += `
linters-settings:
  gocyclo:
    min-complexity: 10
  funlen:
    lines: 80
    statements: 50
  errcheck:
    check-type-assertions: true
    check-blank: true

issues:
  max-issues-per-linter: 0
  max-same-issues: 0

run:
  timeout: 5m
  tests: true
`
	
	return config
}

// Helper functions

func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "\n... (truncated)"
}

func extractJSONFromMarkdown(content string) string {
	// 提取 ```json ... ``` 代码块
	start := -1
	end := -1
	
	lines := bytes.Split([]byte(content), []byte("\n"))
	for i, line := range lines {
		lineStr := string(bytes.TrimSpace(line))
		if start == -1 && (lineStr == "```json" || lineStr == "```") {
			start = i + 1
		} else if start != -1 && lineStr == "```" {
			end = i
			break
		}
	}
	
	if start != -1 && end != -1 {
		return string(bytes.Join(lines[start:end], []byte("\n")))
	}
	
	return content
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func deduplicate(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

