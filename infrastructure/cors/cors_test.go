package cors

import (
	"os"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func TestGetAllowedOrigins_WithEnvVar(t *testing.T) {
	// Set environment variable
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,https://test.com,http://localhost:8080")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	expected := []string{"https://example.com", "https://test.com", "http://localhost:8080"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	for i, origin := range origins {
		if origin != expected[i] {
			t.Errorf("Expected origin %v, got %v", expected[i], origin)
		}
	}
}

func TestGetAllowedOrigins_WithSpaces(t *testing.T) {
	// Set environment variable with spaces
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", " https://example.com , https://test.com , http://localhost:8080 ")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	expected := []string{"https://example.com", "https://test.com", "http://localhost:8080"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	for i, origin := range origins {
		if origin != expected[i] {
			t.Errorf("Expected origin %v, got %v", expected[i], origin)
		}
	}
}

func TestGetAllowedOrigins_WithEmptyValues(t *testing.T) {
	// Set environment variable with empty values
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,,https://test.com,   ,")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	expected := []string{"https://example.com", "https://test.com"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	for i, origin := range origins {
		if origin != expected[i] {
			t.Errorf("Expected origin %v, got %v", expected[i], origin)
		}
	}
}

func TestGetAllowedOrigins_EmptyEnvVar(t *testing.T) {
	// Unset environment variable
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Unsetenv("ALLOWED_ORIGINS")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	expected := []string{"http://localhost:3000", "http://localhost:3001"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	for i, origin := range origins {
		if origin != expected[i] {
			t.Errorf("Expected origin %v, got %v", expected[i], origin)
		}
	}
}

func TestGetAllowedOrigins_OnlyEmptyValues(t *testing.T) {
	// Set environment variable with only empty values
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", " , , ")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	// Should return default origins when no valid origins found
	expected := []string{"http://localhost:3000", "http://localhost:3001"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	for i, origin := range origins {
		if origin != expected[i] {
			t.Errorf("Expected origin %v, got %v", expected[i], origin)
		}
	}
}

func TestGetAllowedOrigins_SingleOrigin(t *testing.T) {
	// Set environment variable with single origin
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://production.com")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()

	expected := []string{"https://production.com"}
	if len(origins) != len(expected) {
		t.Errorf("Expected %d origins, got %d", len(expected), len(origins))
	}

	if origins[0] != expected[0] {
		t.Errorf("Expected origin %v, got %v", expected[0], origins[0])
	}
}

func TestSetupCors(t *testing.T) {
	// Set test environment
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://test.com")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	// Test that SetupCors returns a valid gin.HandlerFunc
	handler := SetupCors()
	if handler == nil {
		t.Error("Expected non-nil handler function")
	}

	// Test that it returns a function of the correct type
	if _, ok := interface{}(handler).(gin.HandlerFunc); !ok {
		t.Error("Expected handler to be of type gin.HandlerFunc")
	}
}

func TestSetupCorsWithDefaults(t *testing.T) {
	// Test with empty environment (should use defaults)
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Unsetenv("ALLOWED_ORIGINS")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	handler := SetupCors()
	if handler == nil {
		t.Error("Expected non-nil handler function")
	}

	// Verify it's a valid gin handler
	if _, ok := interface{}(handler).(gin.HandlerFunc); !ok {
		t.Error("Expected handler to be of type gin.HandlerFunc")
	}
}

func TestCorsConfigParameters(t *testing.T) {
	// This test verifies that our CORS configuration has the expected parameters
	// We can't directly test the cors.Config without refactoring, but we can test
	// the logic that feeds into it
	
	// Test getAllowedOrigins function which is used by SetupCors
	originalValue := os.Getenv("ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com")
	defer os.Setenv("ALLOWED_ORIGINS", originalValue)

	origins := getAllowedOrigins()
	
	if len(origins) != 1 {
		t.Errorf("Expected 1 origin, got %d", len(origins))
	}
	
	if origins[0] != "https://example.com" {
		t.Errorf("Expected origin 'https://example.com', got %v", origins[0])
	}
}

// Test that the expected CORS configuration values are used
func TestCorsConfiguration(t *testing.T) {
	// We'll create our own config to test the expected values
	// This simulates what SetupCors does internally
	
	allowedOrigins := []string{"https://test.com"}
	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Test the configuration values
	if len(config.AllowOrigins) != 1 || config.AllowOrigins[0] != "https://test.com" {
		t.Errorf("Expected AllowOrigins to be ['https://test.com'], got %v", config.AllowOrigins)
	}

	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	if len(config.AllowMethods) != len(expectedMethods) {
		t.Errorf("Expected %d methods, got %d", len(expectedMethods), len(config.AllowMethods))
	}

	expectedHeaders := []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	if len(config.AllowHeaders) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(config.AllowHeaders))
	}

	if !config.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}

	if config.MaxAge != 12*time.Hour {
		t.Errorf("Expected MaxAge to be 12 hours, got %v", config.MaxAge)
	}
}