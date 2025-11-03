package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypePDF      DocumentType = "pdf"
	DocumentTypeDOC      DocumentType = "doc"
	DocumentTypeDOCX     DocumentType = "docx"
	DocumentTypeTXT      DocumentType = "txt"
	DocumentTypeMarkdown DocumentType = "markdown"
)

// ParsedDocument represents a parsed document with metadata
type ParsedDocument struct {
	Type     DocumentType // 文档类型
	Content  string       // 提取的文本内容
	FileName string       // 文件名
	FileSize int64        // 文件大小（字节）
	PageCount int         // 页数（如果适用）
}

// DocumentParser handles parsing of various document formats
type DocumentParser struct{}

// NewDocumentParser creates a new document parser
func NewDocumentParser() *DocumentParser {
	return &DocumentParser{}
}

// ParseFile parses a document file and extracts text content
func (p *DocumentParser) ParseFile(filePath string) (*ParsedDocument, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// 确定文档类型
	docType := p.detectType(filePath)
	if docType == "" {
		return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(filePath))
	}

	// 根据类型解析
	var content string
	var pageCount int

	switch docType {
	case DocumentTypePDF:
		content, pageCount, err = p.parsePDF(filePath)
	case DocumentTypeTXT, DocumentTypeMarkdown:
		content, err = p.parseText(filePath)
	case DocumentTypeDOCX:
		content, err = p.parseDOCX(filePath)
	case DocumentTypeDOC:
		return nil, fmt.Errorf("DOC format requires conversion to DOCX first (use Office or online converter)")
	default:
		return nil, fmt.Errorf("unsupported document type: %s", docType)
	}

	if err != nil {
		return nil, err
	}

	return &ParsedDocument{
		Type:      docType,
		Content:   content,
		FileName:  filepath.Base(filePath),
		FileSize:  fileInfo.Size(),
		PageCount: pageCount,
	}, nil
}

// ParseBytes parses document from bytes with specified file extension
func (p *DocumentParser) ParseBytes(data []byte, fileName string) (*ParsedDocument, error) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "doc-parse-*"+filepath.Ext(fileName))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 写入数据
	if _, err := tmpFile.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// 解析临时文件
	return p.ParseFile(tmpFile.Name())
}

// detectType detects document type from file extension
func (p *DocumentParser) detectType(filePath string) DocumentType {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".pdf":
		return DocumentTypePDF
	case ".doc":
		return DocumentTypeDOC
	case ".docx":
		return DocumentTypeDOCX
	case ".txt":
		return DocumentTypeTXT
	case ".md", ".markdown":
		return DocumentTypeMarkdown
	default:
		return ""
	}
}

// parsePDF parses PDF file and extracts text
func (p *DocumentParser) parsePDF(filePath string) (string, int, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	totalPages := r.NumPage()
	var content strings.Builder

	// 提取每一页的文本
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// 尝试继续处理其他页
			continue
		}

		content.WriteString(text)
		content.WriteString("\n\n")
	}

	return strings.TrimSpace(content.String()), totalPages, nil
}

// parseText parses plain text file
func (p *DocumentParser) parseText(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}
	return string(data), nil
}

// parseDOCX parses DOCX file and extracts text
func (p *DocumentParser) parseDOCX(filePath string) (string, error) {
	// 使用 archive/zip 读取 DOCX (它实际上是 ZIP 文件)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat DOCX: %w", err)
	}

	// 读取所有内容到内存
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX: %w", err)
	}

	// 使用简单的 XML 解析提取文本
	content, err := p.extractDOCXText(data, fileInfo.Size())
	if err != nil {
		return "", fmt.Errorf("failed to extract DOCX text: %w", err)
	}

	return content, nil
}

// extractDOCXText extracts text from DOCX zip data
func (p *DocumentParser) extractDOCXText(data []byte, size int64) (string, error) {
	// 这是一个简化版本，实际应该使用专门的 DOCX 库
	// 如 github.com/nguyenthenguyen/docx 或 github.com/unidoc/unioffice
	
	// 对于 MVP，我们返回一个提示信息
	return "", fmt.Errorf("DOCX parsing requires additional library. Please convert to PDF or TXT format, or use: go get github.com/nguyenthenguyen/docx")
}

// SupportedFormats returns list of supported document formats
func (p *DocumentParser) SupportedFormats() []string {
	return []string{
		"pdf",
		"txt", 
		"md",
		"markdown",
		// "docx", // Requires additional library
		// "doc",  // Requires conversion
	}
}

// ValidateFormat checks if file format is supported
func (p *DocumentParser) ValidateFormat(fileName string) error {
	ext := strings.ToLower(filepath.Ext(fileName))
	ext = strings.TrimPrefix(ext, ".")
	
	supported := p.SupportedFormats()
	for _, format := range supported {
		if ext == format {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported format: %s (supported: %s)", ext, strings.Join(supported, ", "))
}

