package jwt

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

// Test service structure and cache operations only
func TestBlacklistService_Structure(t *testing.T) {
	logger := zap.NewNop()
	
	service := &BlacklistService{
		cache:  make(map[string]time.Time),
		mutex:  sync.RWMutex{},
		logger: logger,
	}
	
	if service == nil {
		t.Error("Expected non-nil blacklist service")
	}
	
	if service.cache == nil {
		t.Error("Expected cache to be initialized")
	}
	
	if service.logger == nil {
		t.Error("Expected logger to be set")
	}
}

// Test TokenBlacklist struct creation
func TestTokenBlacklist_Creation(t *testing.T) {
	now := time.Now()
	token := TokenBlacklist{
		JTI:       "test-jti-123",
		UserEmail: "test@example.com",
		ExpiresAt: now.Add(1 * time.Hour),
		CreatedAt: now,
	}
	
	if token.JTI != "test-jti-123" {
		t.Errorf("Expected JTI 'test-jti-123', got '%s'", token.JTI)
	}
	
	if token.UserEmail != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", token.UserEmail)
	}
	
	if token.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}
	
	if token.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

// Test cache operations directly
func TestBlacklistService_CacheOperations(t *testing.T) {
	service := &BlacklistService{
		cache:  make(map[string]time.Time),
		mutex:  sync.RWMutex{},
		logger: zap.NewNop(),
	}
	
	// Test adding to cache
	jti := "test-jti"
	expiresAt := time.Now().Add(1 * time.Hour)
	
	service.mutex.Lock()
	service.cache[jti] = expiresAt
	service.mutex.Unlock()
	
	// Test reading from cache
	service.mutex.RLock()
	cached, exists := service.cache[jti]
	service.mutex.RUnlock()
	
	if !exists {
		t.Error("Expected token to exist in cache")
	}
	
	if !cached.Equal(expiresAt) {
		t.Error("Expected cached expiration time to match")
	}
	
	// Test cache cleanup logic
	service.mutex.Lock()
	// Add expired token
	service.cache["expired"] = time.Now().Add(-1 * time.Hour)
	service.mutex.Unlock()
	
	// Run cleanup
	service.CleanupExpiredTokens()
	
	// Verify cleanup
	service.mutex.RLock()
	_, expiredExists := service.cache["expired"]
	_, validExists := service.cache[jti]
	service.mutex.RUnlock()
	
	if expiredExists {
		t.Error("Expected expired token to be removed")
	}
	
	if !validExists {
		t.Error("Expected valid token to remain")
	}
}

// Test concurrent cache access
func TestBlacklistService_ConcurrentCacheAccess(t *testing.T) {
	service := &BlacklistService{
		cache:  make(map[string]time.Time),
		mutex:  sync.RWMutex{},
		logger: zap.NewNop(),
	}
	
	var wg sync.WaitGroup
	numGoroutines := 5
	numOperations := 10
	
	// Start multiple goroutines performing cache operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < numOperations; j++ {
				jti := fmt.Sprintf("concurrent-jti-%d-%d", id, j)
				expiresAt := time.Now().Add(1 * time.Hour)
				
				// Add to cache
				service.mutex.Lock()
				service.cache[jti] = expiresAt
				service.mutex.Unlock()
				
				// Read from cache
				service.mutex.RLock()
				_, exists := service.cache[jti]
				service.mutex.RUnlock()
				
				if !exists {
					t.Errorf("Expected token %s to exist in cache", jti)
					return
				}
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify final cache state
	service.mutex.RLock()
	cacheSize := len(service.cache)
	service.mutex.RUnlock()
	
	expectedSize := numGoroutines * numOperations
	if cacheSize != expectedSize {
		t.Errorf("Expected %d tokens in cache, got %d", expectedSize, cacheSize)
	}
}

// Benchmark cache operations
func BenchmarkBlacklistService_CacheRead(b *testing.B) {
	service := &BlacklistService{
		cache:  make(map[string]time.Time),
		mutex:  sync.RWMutex{},
		logger: zap.NewNop(),
	}
	
	// Pre-populate cache
	service.cache["benchmark-jti"] = time.Now().Add(1 * time.Hour)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.mutex.RLock()
		_, exists := service.cache["benchmark-jti"]
		service.mutex.RUnlock()
		_ = exists
	}
}

func BenchmarkBlacklistService_CacheWrite(b *testing.B) {
	service := &BlacklistService{
		cache:  make(map[string]time.Time),
		mutex:  sync.RWMutex{},
		logger: zap.NewNop(),
	}
	
	expiresAt := time.Now().Add(1 * time.Hour)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jti := fmt.Sprintf("benchmark-jti-%d", i)
		service.mutex.Lock()
		service.cache[jti] = expiresAt
		service.mutex.Unlock()
	}
}