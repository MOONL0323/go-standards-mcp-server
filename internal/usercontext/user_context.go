package usercontext

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UserContext represents a user's execution context
type UserContext struct {
	UserID       string
	SessionID    string
	WorkspaceDir string // User private workspace
	SharedDir    string // Shared resources directory
	CreatedAt    time.Time
	LastAccessAt time.Time
}

// NewUserContext creates a new user context
func NewUserContext(userID, sessionID, baseDir string) *UserContext {
	now := time.Now()
	ctx := &UserContext{
		UserID:       userID,
		SessionID:    sessionID,
		WorkspaceDir: filepath.Join(baseDir, "users", userID),
		SharedDir:    filepath.Join(baseDir, "shared"),
		CreatedAt:    now,
		LastAccessAt: now,
	}
	
	// Initialize user directories
	ctx.initDirectories()
	
	return ctx
}

// initDirectories creates necessary user directories
func (uc *UserContext) initDirectories() error {
	dirs := []string{
		uc.GetTempDir(),
		uc.GetCacheDir(),
		uc.GetHistoryDir(),
		uc.GetReportsDir(),
		uc.GetGitConfigDir(),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// GetTempDir returns user's temporary directory
func (uc *UserContext) GetTempDir() string {
	return filepath.Join(uc.WorkspaceDir, "temp")
}

// GetCacheDir returns user's cache directory
func (uc *UserContext) GetCacheDir() string {
	return filepath.Join(uc.WorkspaceDir, "cache")
}

// GetHistoryDir returns user's analysis history directory
func (uc *UserContext) GetHistoryDir() string {
	return filepath.Join(uc.WorkspaceDir, "history")
}

// GetReportsDir returns user's reports directory
func (uc *UserContext) GetReportsDir() string {
	return filepath.Join(uc.WorkspaceDir, "reports")
}

// GetGitConfigDir returns user's git configuration directory
func (uc *UserContext) GetGitConfigDir() string {
	return filepath.Join(uc.WorkspaceDir, "git-config")
}

// GetSharedDocumentsDir returns shared documents directory (all users)
func (uc *UserContext) GetSharedDocumentsDir() string {
	return filepath.Join(uc.SharedDir, "documents")
}

// GetSharedTemplatesDir returns shared templates directory (all users)
func (uc *UserContext) GetSharedTemplatesDir() string {
	return filepath.Join(uc.SharedDir, "templates")
}

// GetSharedConfigsDir returns shared configs directory (all users)
func (uc *UserContext) GetSharedConfigsDir() string {
	return filepath.Join(uc.SharedDir, "configs")
}

// UpdateLastAccess updates the last access timestamp
func (uc *UserContext) UpdateLastAccess() {
	uc.LastAccessAt = time.Now()
}

// IsExpired checks if the context has expired
func (uc *UserContext) IsExpired(timeout time.Duration) bool {
	return time.Since(uc.LastAccessAt) > timeout
}

// CleanupTempFiles removes temporary files for this user
func (uc *UserContext) CleanupTempFiles() error {
	tempDir := uc.GetTempDir()
	if err := os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("failed to cleanup temp files: %w", err)
	}
	return os.MkdirAll(tempDir, 0755)
}

// SessionManager manages user sessions and contexts
type SessionManager struct {
	sessions       map[string]*UserContext
	mu             sync.RWMutex
	baseDir        string
	sessionTimeout time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(baseDir string, sessionTimeout time.Duration) *SessionManager {
	if sessionTimeout == 0 {
		sessionTimeout = 3600 * time.Second // Default: 1 hour
	}
	
	sm := &SessionManager{
		sessions:       make(map[string]*UserContext),
		baseDir:        baseDir,
		sessionTimeout: sessionTimeout,
	}
	
	// Initialize shared directories
	sm.initSharedDirectories()
	
	// Start cleanup goroutine
	go sm.cleanupLoop()
	
	return sm
}

// initSharedDirectories creates shared resource directories
func (sm *SessionManager) initSharedDirectories() error {
	sharedDir := filepath.Join(sm.baseDir, "shared")
	dirs := []string{
		filepath.Join(sharedDir, "documents"),
		filepath.Join(sharedDir, "templates"),
		filepath.Join(sharedDir, "configs"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create shared directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// GetOrCreateUserContext gets existing or creates new user context
func (sm *SessionManager) GetOrCreateUserContext(userID, sessionID string) *UserContext {
	key := fmt.Sprintf("%s:%s", userID, sessionID)
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if ctx, exists := sm.sessions[key]; exists {
		ctx.UpdateLastAccess()
		return ctx
	}
	
	ctx := NewUserContext(userID, sessionID, sm.baseDir)
	sm.sessions[key] = ctx
	
	return ctx
}

// GetUserContext gets existing user context
func (sm *SessionManager) GetUserContext(userID, sessionID string) (*UserContext, bool) {
	key := fmt.Sprintf("%s:%s", userID, sessionID)
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	ctx, exists := sm.sessions[key]
	if exists {
		ctx.UpdateLastAccess()
	}
	
	return ctx, exists
}

// RemoveSession removes a user session
func (sm *SessionManager) RemoveSession(userID, sessionID string) error {
	key := fmt.Sprintf("%s:%s", userID, sessionID)
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if ctx, exists := sm.sessions[key]; exists {
		// Cleanup temporary files
		if err := ctx.CleanupTempFiles(); err != nil {
			return err
		}
		delete(sm.sessions, key)
	}
	
	return nil
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	cleaned := 0
	for key, ctx := range sm.sessions {
		if ctx.IsExpired(sm.sessionTimeout) {
			// Cleanup temporary files
			ctx.CleanupTempFiles()
			delete(sm.sessions, key)
			cleaned++
		}
	}
	
	return cleaned
}

// cleanupLoop periodically cleans up expired sessions
func (sm *SessionManager) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		cleaned := sm.CleanupExpiredSessions()
		if cleaned > 0 {
			fmt.Printf("Cleaned up %d expired sessions\n", cleaned)
		}
	}
}

// GetActiveSessionCount returns the number of active sessions
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// GetUserSessionCount returns the number of sessions for a specific user
func (sm *SessionManager) GetUserSessionCount(userID string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	count := 0
	for _, ctx := range sm.sessions {
		if ctx.UserID == userID {
			count++
		}
	}
	return count
}

// ListActiveSessions returns all active session keys
func (sm *SessionManager) ListActiveSessions() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	keys := make([]string, 0, len(sm.sessions))
	for key := range sm.sessions {
		keys = append(keys, key)
	}
	return keys
}

// GetStats returns session manager statistics
func (sm *SessionManager) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	userCounts := make(map[string]int)
	for _, ctx := range sm.sessions {
		userCounts[ctx.UserID]++
	}
	
	return map[string]interface{}{
		"total_sessions":     len(sm.sessions),
		"unique_users":       len(userCounts),
		"sessions_per_user":  userCounts,
		"session_timeout":    sm.sessionTimeout.String(),
	}
}
