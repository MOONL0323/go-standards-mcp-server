package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/MOONL0323/go-standards-mcp-server/internal/converter"
	"github.com/MOONL0323/go-standards-mcp-server/internal/parser"
	"github.com/MOONL0323/go-standards-mcp-server/internal/storage"
	"go.uber.org/zap"
)

// DocumentService handles document upload and conversion workflows
type DocumentService struct {
	parser    *parser.DocumentParser
	converter *converter.AIConverter
	docStore  *storage.DocumentStorage
	cfgStore  *storage.ConfigStorage
	logger    *zap.Logger
}

// NewDocumentService creates a new document service
func NewDocumentService(logger *zap.Logger) (*DocumentService, error) {
	// 初始化存储
	docStore, err := storage.NewDocumentStorage("./data/documents")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize document storage: %w", err)
	}

	cfgStore, err := storage.NewConfigStorage("./configs/custom")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config storage: %w", err)
	}

	return &DocumentService{
		parser:    parser.NewDocumentParser(),
		converter: converter.NewAIConverter(),
		docStore:  docStore,
		cfgStore:  cfgStore,
		logger:    logger,
	}, nil
}

// UploadDocumentRequest represents a document upload request
type UploadDocumentRequest struct {
	File        io.Reader // 文件内容
	FileName    string    // 文件名
	Name        string    // 配置名称（用户指定）
	Description string    // 描述
}

// UploadDocumentResponse represents the response of document upload
type UploadDocumentResponse struct {
	DocumentID        string   `json:"document_id"`
	ConfigName        string   `json:"config_name"`
	Summary           string   `json:"summary"`
	ExtractedRules    []string `json:"extracted_rules"`
	Confidence        float64  `json:"confidence"`
	Success           bool     `json:"success"`
	Message           string   `json:"message"`
}

// UploadDocument handles the complete document upload and conversion workflow
func (s *DocumentService) UploadDocument(ctx context.Context, req *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	s.logger.Info("Starting document upload", 
		zap.String("file", req.FileName),
		zap.String("name", req.Name))

	// 1. 验证文件格式
	if err := s.parser.ValidateFormat(req.FileName); err != nil {
		return nil, fmt.Errorf("invalid file format: %w", err)
	}

	// 2. 读取文件内容到内存（需要多次使用）
	fileData, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 3. 解析文档
	parsed, err := s.parser.ParseBytes(fileData, req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	s.logger.Info("Document parsed successfully",
		zap.String("type", string(parsed.Type)),
		zap.Int64("size", parsed.FileSize))

	// 4. 使用 AI 转换为配置
	convResult, err := s.converter.Convert(ctx, parsed.Content, parsed.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	// 5. 保存文档
	metadata := &storage.DocumentMetadata{
		Name:              req.Name,
		OriginalFile:      req.FileName,
		FileType:          string(parsed.Type),
		FileSize:          parsed.FileSize,
		ConfigName:        req.Name,
		Description:       req.Description,
		ConversionSummary: convResult.Summary,
		ExtractedRules:    convResult.Rules,
		Confidence:        convResult.Confidence,
	}

	if err := s.docStore.SaveDocument(bytes.NewReader(fileData), metadata); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// 6. 保存生成的配置
	if err := s.docStore.SaveConfig(metadata.ID, convResult.Config); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	// 7. 将配置注册到配置管理系统
	if err := s.cfgStore.Save(req.Name, convResult.Config, req.Description); err != nil {
		s.logger.Warn("Failed to register config", zap.Error(err))
		// 不中断流程，继续返回结果
	}

	s.logger.Info("Document uploaded and converted successfully",
		zap.String("document_id", metadata.ID),
		zap.String("config_name", req.Name),
		zap.Float64("confidence", convResult.Confidence))

	return &UploadDocumentResponse{
		DocumentID:     metadata.ID,
		ConfigName:     req.Name,
		Summary:        convResult.Summary,
		ExtractedRules: convResult.Rules,
		Confidence:     convResult.Confidence,
		Success:        true,
		Message:        fmt.Sprintf("Document uploaded and converted successfully. Confidence: %.1f%%", convResult.Confidence*100),
	}, nil
}

// ListDocuments lists all uploaded documents
func (s *DocumentService) ListDocuments() ([]*storage.DocumentMetadata, error) {
	return s.docStore.List()
}

// GetDocument retrieves document details
func (s *DocumentService) GetDocument(id string) (*storage.DocumentMetadata, error) {
	return s.docStore.Get(id)
}

// GetDocumentConfig retrieves the generated configuration for a document
func (s *DocumentService) GetDocumentConfig(id string) (string, error) {
	return s.docStore.GetConfig(id)
}

// DeleteDocument deletes a document and its configuration
func (s *DocumentService) DeleteDocument(id string) error {
	// 获取文档信息以删除关联的配置
	doc, err := s.docStore.Get(id)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// 删除配置
	if doc.ConfigName != "" {
		if err := s.cfgStore.Delete(doc.ConfigName); err != nil {
			s.logger.Warn("Failed to delete config", 
				zap.String("config", doc.ConfigName), 
				zap.Error(err))
		}
	}

	// 删除文档
	return s.docStore.Delete(id)
}

// GetStats returns document storage statistics
func (s *DocumentService) GetStats() (map[string]interface{}, error) {
	return s.docStore.GetStats()
}

