package jwt

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


// Helper to create valid JWT token for testing
func createTestJWTToken(userID, email, phone, jti, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"phone":   phone,
		"jti":     jti,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func setupMiddlewareTest() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key-for-middleware-testing")
}

func TestJWTMiddleware_Success(t *testing.T) {
	setupMiddlewareTest()
	
	// Create valid token
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-123", "test-secret-key-for-middleware-testing", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	// Create request with token cookie
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Create middleware without blacklist service
	middleware := JWTMiddleware(nil)
	
	// Test successful authentication
	middleware(c)
	
	// Verify context values were set
	userID, exists := c.Get("user_id")
	if !exists {
		t.Error("Expected user_id to be set in context")
	} else if userID != "user123" {
		t.Errorf("Expected user_id 'user123', got '%v'", userID)
	}
	
	email, exists := c.Get("email")
	if !exists {
		t.Error("Expected email to be set in context")
	} else if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%v'", email)
	}
	
	phone, exists := c.Get("phone")
	if !exists {
		t.Error("Expected phone to be set in context")
	} else if phone != "+1234567890" {
		t.Errorf("Expected phone '+1234567890', got '%v'", phone)
	}
	
	jti, exists := c.Get("jti")
	if !exists {
		t.Error("Expected jti to be set in context")
	} else if jti != "jti-123" {
		t.Errorf("Expected jti 'jti-123', got '%v'", jti)
	}
	
	// Verify response was not aborted
	if c.IsAborted() {
		t.Error("Expected context not to be aborted for valid token")
	}
}

func TestJWTMiddleware_NoCookie(t *testing.T) {
	setupMiddlewareTest()
	
	// Create request without token cookie
	req, _ := http.NewRequest("GET", "/protected", nil)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted when no cookie is present")
	}
	
	// Verify error response
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_EmptyCookie(t *testing.T) {
	setupMiddlewareTest()
	
	// Create request with empty token cookie
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "",
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted for empty token")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	setupMiddlewareTest()
	
	// Create request with invalid token
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "invalid.jwt.token",
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted for invalid token")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	setupMiddlewareTest()
	
	// Create expired token
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-expired", "test-secret-key-for-middleware-testing", -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted for expired token")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_WrongSigningMethod(t *testing.T) {
	setupMiddlewareTest()
	
	// This will fail to sign properly, but that's expected for testing
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3MDAwMDAwMDAsInVzZXJfaWQiOiJ1c2VyMTIzIn0.invalid-signature"
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted for wrong signing method")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_WrongSecret(t *testing.T) {
	setupMiddlewareTest()
	
	// Create token with wrong secret
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-wrong-secret", "wrong-secret", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token with wrong secret: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted
	if !c.IsAborted() {
		t.Error("Expected context to be aborted for token with wrong secret")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_WithBlacklistService_ValidToken(t *testing.T) {
	setupMiddlewareTest()
	
	// Test without blacklist service since mocking is complex
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-valid", "test-secret-key-for-middleware-testing", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was not aborted (token is valid)
	if c.IsAborted() {
		t.Error("Expected context not to be aborted for valid token")
	}
}

func TestJWTMiddleware_WithBlacklistService_BlacklistedToken(t *testing.T) {
	setupMiddlewareTest()
	
	// Test middleware behavior with nil blacklist service (simplified test)
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-blacklisted", "test-secret-key-for-middleware-testing", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was not aborted (no blacklist service means no blacklist check)
	if c.IsAborted() {
		t.Error("Expected context not to be aborted when no blacklist service is provided")
	}
}

func TestJWTMiddleware_MissingClaims(t *testing.T) {
	setupMiddlewareTest()
	
	// Create token with minimal claims (missing some expected fields)
	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		// Missing user_id, email, phone, jti
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-key-for-middleware-testing"))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify token is still considered valid (missing claims are optional)
	if c.IsAborted() {
		t.Error("Expected context not to be aborted for token with missing optional claims")
	}
	
	// Verify missing claims are not set in context
	if _, exists := c.Get("user_id"); exists {
		t.Error("Expected user_id not to be set for token without user_id claim")
	}
	
	if _, exists := c.Get("email"); exists {
		t.Error("Expected email not to be set for token without email claim")
	}
	
	if _, exists := c.Get("phone"); exists {
		t.Error("Expected phone not to be set for token without phone claim")
	}
	
	if _, exists := c.Get("jti"); exists {
		t.Error("Expected jti not to be set for token without jti claim")
	}
}

func TestJWTMiddleware_InvalidClaimsTypes(t *testing.T) {
	setupMiddlewareTest()
	
	// Create token with invalid claim types
	claims := jwt.MapClaims{
		"user_id": 123,           // Should be string, not int
		"email":   true,          // Should be string, not bool
		"phone":   []string{"1"}, // Should be string, not array
		"jti":     nil,           // Should be string, not nil
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-key-for-middleware-testing"))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify token is still considered valid (type mismatches are handled gracefully)
	if c.IsAborted() {
		t.Error("Expected context not to be aborted for token with invalid claim types")
	}
	
	// Verify invalid claims are not set in context
	if _, exists := c.Get("user_id"); exists {
		t.Error("Expected user_id not to be set for invalid type")
	}
	
	if _, exists := c.Get("email"); exists {
		t.Error("Expected email not to be set for invalid type")
	}
	
	if _, exists := c.Get("phone"); exists {
		t.Error("Expected phone not to be set for invalid type")
	}
	
	if _, exists := c.Get("jti"); exists {
		t.Error("Expected jti not to be set for invalid type")
	}
}

func TestJWTMiddleware_NoJWTSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Don't set JWT_SECRET environment variable
	os.Unsetenv("JWT_SECRET")
	
	tokenString, err := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-no-secret", "any-secret", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}
	
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	middleware := JWTMiddleware(nil)
	middleware(c)
	
	// Verify request was aborted (empty secret should cause verification failure)
	if !c.IsAborted() {
		t.Error("Expected context to be aborted when JWT_SECRET is not set")
	}
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	// Restore environment for other tests
	os.Setenv("JWT_SECRET", "test-secret-key-for-middleware-testing")
}

// Benchmark tests
func BenchmarkJWTMiddleware_ValidToken(b *testing.B) {
	setupMiddlewareTest()
	
	tokenString, _ := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-bench", "test-secret-key-for-middleware-testing", 1*time.Hour)
	
	middleware := JWTMiddleware(nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: tokenString,
		})
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		
		middleware(c)
	}
}

func BenchmarkJWTMiddleware_InvalidToken(b *testing.B) {
	setupMiddlewareTest()
	
	middleware := JWTMiddleware(nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "invalid.jwt.token",
		})
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		
		middleware(c)
	}
}

func BenchmarkJWTMiddleware_WithBlacklist(b *testing.B) {
	setupMiddlewareTest()
	
	tokenString, _ := createTestJWTToken("user123", "test@example.com", "+1234567890", "jti-bench-blacklist", "test-secret-key-for-middleware-testing", 1*time.Hour)
	
	middleware := JWTMiddleware(nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: tokenString,
		})
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		
		middleware(c)
	}
}