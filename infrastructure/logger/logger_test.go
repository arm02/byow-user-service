package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLoggerTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createTestLogger() (*zap.Logger, *bytes.Buffer) {
	// Create a buffer to capture log output
	buffer := &bytes.Buffer{}
	
	// Create a custom core that writes to our buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(buffer), zapcore.InfoLevel)
	logger := zap.New(core)
	
	return logger, buffer
}

func TestLogRequestBody_GET_Request(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// GET requests should not log request body
	logOutput := buffer.String()
	if strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected no request body logging for GET request")
	}
}

func TestLogRequestBody_POST_WithBody(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	requestBody := `{"test": "data"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// POST requests should log request body (unless in skip paths)
	logOutput := buffer.String()
	if !strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected request body logging for POST request")
	}

	if !strings.Contains(logOutput, "POST") {
		t.Error("Expected method 'POST' in log output")
	}

	if !strings.Contains(logOutput, `test`) && !strings.Contains(logOutput, `data`) {
		t.Error("Expected request body content in log output")
	}
}

func TestLogRequestBody_POST_SkipPath(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/auth/users/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	requestBody := `{"email": "test@example.com", "password": "secret"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/users/login", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Login path should be skipped for security
	logOutput := buffer.String()
	if strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected no request body logging for login path")
	}
}

func TestLogRequestBody_POST_RegisterSkipPath(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/auth/users/register", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	requestBody := `{"email": "test@example.com", "password": "secret"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/users/register", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Register path should be skipped for security
	logOutput := buffer.String()
	if strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected no request body logging for register path")
	}
}

func TestLogRequestBody_POST_ChangePasswordSkipPath(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/auth/users/change-password", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	requestBody := `{"old_password": "old", "new_password": "new"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/users/change-password", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Change password path should be skipped for security
	logOutput := buffer.String()
	if strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected no request body logging for change password path")
	}
}

func TestLogRequestBody_PUT_WithBody(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.PUT("/api/update", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "updated"})
	})

	requestBody := `{"name": "updated name"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/update", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// PUT requests should log request body
	logOutput := buffer.String()
	if !strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected request body logging for PUT request")
	}

	if !strings.Contains(logOutput, "PUT") {
		t.Error("Expected method 'PUT' in log output")
	}
}

func TestLogRequestBody_EmptyBody(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should still log even with empty body
	logOutput := buffer.String()
	if !strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected request body logging even for empty body")
	}
}

func TestLogRequestBody_NilBody(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Should not crash with nil body
	logOutput := buffer.String()
	if strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected no request body logging for nil body")
	}
}

func TestLogRequestBody_AllSkipPaths(t *testing.T) {
	// Test that all skip paths are properly handled
	skipPaths := []string{
		"/auth/users/login",
		"/auth/users/change-password",
		"/auth/users/register",
	}

	for _, path := range skipPaths {
		t.Run(path, func(t *testing.T) {
			logger, buffer := createTestLogger()
			router := setupLoggerTestRouter()
			
			router.Use(LogRequestBody(logger))
			router.POST(path, func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})

			requestBody := `{"sensitive": "data"}`
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", path, strings.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			logOutput := buffer.String()
			if strings.Contains(logOutput, "Request Payload") {
				t.Errorf("Expected no request body logging for skip path %s", path)
			}
		})
	}
}

func TestLogRequestBody_NonSkipPath(t *testing.T) {
	logger, buffer := createTestLogger()
	router := setupLoggerTestRouter()
	
	router.Use(LogRequestBody(logger))
	router.POST("/api/public/endpoint", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	requestBody := `{"public": "data"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/public/endpoint", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Non-skip paths should log request body
	logOutput := buffer.String()
	if !strings.Contains(logOutput, "Request Payload") {
		t.Error("Expected request body logging for non-skip path")
	}

	if !strings.Contains(logOutput, "/api/public/endpoint") {
		t.Error("Expected path in log output")
	}
}