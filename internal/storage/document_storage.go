package storage

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// DocumentMetadata stores metadata about uploaded documents
type DocumentMetadata struct {
	ID           string    `json:"id"`            // 文档唯一ID（基于文件哈希）
	Name         string    `json:"name"`          // 文档名称
	OriginalFile string    `json:"original_file"` // 原始文件名
	FileType     string    `json:"file_type"`     // 文件类型
	FileSize     int64     `json:"file_size"`     // 文件大小
	UploadedAt   time.Time `json:"uploaded_at"`   // 上传时间
	ConfigName   string    `json:"config_name"`   // 关联的配置名称
	Description  string    `json:"description"`   // 描述
	Version      int       `json:"version"`       // 版本号
	Hash         string    `json:"hash"`          // 文件哈希值
	ExtractedText string   `json:"-"`             // 提取的文本（不序列化到 JSON）
	
	// AI 转换结果
	ConversionSummary string   `json:"conversion_summary"` // 转换摘要
	ExtractedRules    []string `json:"extracted_rules"`    // 提取的规则
	Confidence        float64  `json:"confidence"`         // 转换置信度
}

// DocumentStorage manages uploaded documents
type DocumentStorage struct {
	baseDir string // 存储根目录
}

// NewDocumentStorage creates a new document storage
func NewDocumentStorage(baseDir string) (*DocumentStorage, error) {
	// 创建必要的目录结构
	dirs := []string{
		baseDir,
		filepath.Join(baseDir, "documents"), // 原始文档
		filepath.Join(baseDir, "metadata"),  // 元数据
		filepath.Join(baseDir, "configs"),   // 生成的配置
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return &DocumentStorage{
		baseDir: baseDir,
	}, nil
}

// SaveDocument saves a document with its metadata
func (s *DocumentStorage) SaveDocument(file io.Reader, metadata *DocumentMetadata) error {
	// 计算文件哈希并保存文件
	hash, filePath, err := s.saveFile(file, metadata.OriginalFile)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	
	// 更新元数据
	metadata.Hash = hash
	metadata.ID = hash[:16] // 使用哈希前16位作为ID
	metadata.UploadedAt = time.Now()
	
	// 检查是否已存在相同文件
	if existing, err := s.GetByHash(hash); err == nil {
		// 文件已存在，返回现有记录的信息
		return fmt.Errorf("document already exists: %s (uploaded at %s)", 
			existing.OriginalFile, existing.UploadedAt.Format(time.RFC3339))
	}
	
	// 保存元数据
	if err := s.saveMetadata(metadata); err != nil {
		// 如果元数据保存失败，删除已保存的文件
		os.Remove(filePath)
		return fmt.Errorf("failed to save metadata: %w", err)
	}
	
	return nil
}

// SaveConfig saves the generated configuration for a document
func (s *DocumentStorage) SaveConfig(documentID, configContent string) error {
	configPath := filepath.Join(s.baseDir, "configs", documentID+".yaml")
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	return nil
}

// GetConfig retrieves the generated configuration for a document
func (s *DocumentStorage) GetConfig(documentID string) (string, error) {
	configPath := filepath.Join(s.baseDir, "configs", documentID+".yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("config not found for document: %s", documentID)
		}
		return "", fmt.Errorf("failed to read config: %w", err)
	}
	
	return string(data), nil
}

// Get retrieves document metadata by ID
func (s *DocumentStorage) Get(id string) (*DocumentMetadata, error) {
	metaPath := filepath.Join(s.baseDir, "metadata", id+".json")
	
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("document not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}
	
	var metadata DocumentMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	
	return &metadata, nil
}

// GetByHash retrieves document metadata by file hash
func (s *DocumentStorage) GetByHash(hash string) (*DocumentMetadata, error) {
	// 遍历所有元数据文件查找匹配的哈希
	metaDir := filepath.Join(s.baseDir, "metadata")
	files, err := os.ReadDir(metaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata directory: %w", err)
	}
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(metaDir, file.Name()))
		if err != nil {
			continue
		}
		
		var metadata DocumentMetadata
		if err := json.Unmarshal(data, &metadata); err != nil {
			continue
		}
		
		if metadata.Hash == hash {
			return &metadata, nil
		}
	}
	
	return nil, fmt.Errorf("document not found with hash: %s", hash)
}

// List lists all uploaded documents
func (s *DocumentStorage) List() ([]*DocumentMetadata, error) {
	metaDir := filepath.Join(s.baseDir, "metadata")
	files, err := os.ReadDir(metaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata directory: %w", err)
	}
	
	var documents []*DocumentMetadata
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(metaDir, file.Name()))
		if err != nil {
			continue
		}
		
		var metadata DocumentMetadata
		if err := json.Unmarshal(data, &metadata); err != nil {
			continue
		}
		
		documents = append(documents, &metadata)
	}
	
	// 按上传时间排序（最新的在前）
	sort.Slice(documents, func(i, j int) bool {
		return documents[i].UploadedAt.After(documents[j].UploadedAt)
	})
	
	return documents, nil
}

// Delete deletes a document and its associated files
func (s *DocumentStorage) Delete(id string) error {
	// 删除元数据
	metaPath := filepath.Join(s.baseDir, "metadata", id+".json")
	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}
	
	// 删除配置文件
	configPath := filepath.Join(s.baseDir, "configs", id+".yaml")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete config: %w", err)
	}
	
	// 删除原始文档（通过查找匹配的文件）
	docDir := filepath.Join(s.baseDir, "documents")
	files, err := os.ReadDir(docDir)
	if err == nil {
		for _, file := range files {
			if filepath.Base(file.Name())[:len(id)] == id {
				os.Remove(filepath.Join(docDir, file.Name()))
			}
		}
	}
	
	return nil
}

// GetStats returns storage statistics
func (s *DocumentStorage) GetStats() (map[string]interface{}, error) {
	docs, err := s.List()
	if err != nil {
		return nil, err
	}
	
	var totalSize int64
	typeCount := make(map[string]int)
	
	for _, doc := range docs {
		totalSize += doc.FileSize
		typeCount[doc.FileType]++
	}
	
	return map[string]interface{}{
		"total_documents": len(docs),
		"total_size":      totalSize,
		"type_count":      typeCount,
		"storage_path":    s.baseDir,
	}, nil
}

// saveFile saves uploaded file and returns its hash
func (s *DocumentStorage) saveFile(file io.Reader, originalName string) (string, string, error) {
	// 读取文件内容并计算哈希
	hasher := sha256.New()
	var buf []byte
	var err error
	
	if seeker, ok := file.(io.ReadSeeker); ok {
		// 如果支持 Seek，先读取计算哈希，然后 Seek 回去
		buf, err = io.ReadAll(io.TeeReader(file, hasher))
		if err != nil {
			return "", "", err
		}
		seeker.Seek(0, 0)
	} else {
		// 否则直接读取
		buf, err = io.ReadAll(file)
		if err != nil {
			return "", "", err
		}
		hasher.Write(buf)
	}
	
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	
	// 保存文件
	ext := filepath.Ext(originalName)
	fileName := hash[:16] + ext
	filePath := filepath.Join(s.baseDir, "documents", fileName)
	
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		return "", "", err
	}
	
	return hash, filePath, nil
}

// saveMetadata saves document metadata
func (s *DocumentStorage) saveMetadata(metadata *DocumentMetadata) error {
	metaPath := filepath.Join(s.baseDir, "metadata", metadata.ID+".json")
	
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(metaPath, data, 0644)
}

