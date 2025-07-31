package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with details",
			appError: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test message",
				Status:  400,
				Details: "Additional details",
			},
			expected: "TEST_ERROR: Test message (Additional details)",
		},
		{
			name: "error without details",
			appError: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test message",
				Status:  400,
			},
			expected: "TEST_ERROR: Test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			if result != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	message := "validation failed"
	err := NewValidationError(message)
	
	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected code 'VALIDATION_ERROR', got %v", err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%v', got %v", message, err.Message)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, err.Status)
	}
}

func TestNewNotFoundError(t *testing.T) {
	resource := "user"
	err := NewNotFoundError(resource)
	
	if err.Code != "NOT_FOUND" {
		t.Errorf("Expected code 'NOT_FOUND', got %v", err.Code)
	}
	if err.Message != "user not found" {
		t.Errorf("Expected message 'user not found', got %v", err.Message)
	}
	if err.Status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, err.Status)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "unauthorized access"
	err := NewUnauthorizedError(message)
	
	if err.Code != "UNAUTHORIZED" {
		t.Errorf("Expected code 'UNAUTHORIZED', got %v", err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%v', got %v", message, err.Message)
	}
	if err.Status != http.StatusUnauthorized {
		t.Errorf("Expected status %v, got %v", http.StatusUnauthorized, err.Status)
	}
}

func TestNewConflictError(t *testing.T) {
	message := "resource conflict"
	err := NewConflictError(message)
	
	if err.Code != "CONFLICT" {
		t.Errorf("Expected code 'CONFLICT', got %v", err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%v', got %v", message, err.Message)
	}
	if err.Status != http.StatusConflict {
		t.Errorf("Expected status %v, got %v", http.StatusConflict, err.Status)
	}
}

func TestNewInternalError(t *testing.T) {
	message := "internal server error"
	err := NewInternalError(message)
	
	if err.Code != "INTERNAL_ERROR" {
		t.Errorf("Expected code 'INTERNAL_ERROR', got %v", err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%v', got %v", message, err.Message)
	}
	if err.Status != http.StatusInternalServerError {
		t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, err.Status)
	}
}

func TestNewBadRequestError(t *testing.T) {
	message := "bad request"
	err := NewBadRequestError(message)
	
	if err.Code != "BAD_REQUEST" {
		t.Errorf("Expected code 'BAD_REQUEST', got %v", err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message '%v', got %v", message, err.Message)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, err.Status)
	}
}

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "app error",
			err:      &AppError{Code: "TEST", Message: "test", Status: 400},
			expected: true,
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr, ok := IsAppError(tt.err)
			if ok != tt.expected {
				t.Errorf("IsAppError() ok = %v, want %v", ok, tt.expected)
			}
			if tt.expected && appErr == nil {
				t.Error("Expected appErr to not be nil when ok is true")
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
		isApp   bool
	}{
		{
			name:    "wrap app error",
			err:     &AppError{Code: "TEST", Message: "test", Status: 400},
			message: "wrapper message",
			isApp:   true,
		},
		{
			name:    "wrap standard error",
			err:     errors.New("standard error"),
			message: "wrapper message",
			isApp:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.message)
			
			if tt.isApp {
				// Should return original AppError
				if result.Code != "TEST" {
					t.Errorf("Expected original AppError to be returned")
				}
			} else {
				// Should create new internal error
				if result.Code != "INTERNAL_ERROR" {
					t.Errorf("Expected INTERNAL_ERROR code, got %v", result.Code)
				}
				if result.Message != tt.message {
					t.Errorf("Expected message '%v', got %v", tt.message, result.Message)
				}
				if result.Status != http.StatusInternalServerError {
					t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, result.Status)
				}
				if result.Details != tt.err.Error() {
					t.Errorf("Expected details '%v', got %v", tt.err.Error(), result.Details)
				}
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    *AppError
		code   string
		status int
	}{
		{"ErrUserNotFound", ErrUserNotFound, "NOT_FOUND", http.StatusNotFound},
		{"ErrInvalidCredentials", ErrInvalidCredentials, "INVALID_CREDENTIALS", http.StatusUnauthorized},
		{"ErrUserNotVerified", ErrUserNotVerified, "USER_NOT_VERIFIED", http.StatusUnauthorized},
		{"ErrInvalidOldPassword", ErrInvalidOldPassword, "INVALID_OLD_PASSWORD", http.StatusBadRequest},
		{"ErrEmailAlreadyExists", ErrEmailAlreadyExists, "EMAIL_ALREADY_REGISTERED", http.StatusConflict},
		{"ErrPhoneAlreadyExists", ErrPhoneAlreadyExists, "PHONE_ALREADY_REGISTERED", http.StatusConflict},
		{"ErrEmailOrPhoneAlreadyRegistered", ErrEmailOrPhoneAlreadyRegistered, "EMAIL_OR_PHONE_ALREADY_REGISTERED", http.StatusConflict},
		{"ErrInvalidOTP", ErrInvalidOTP, "OTP_INVALID", http.StatusBadRequest},
		{"ErrExpiredOTP", ErrExpiredOTP, "OTP_EXPIRED", http.StatusBadRequest},
		{"ErrInvalidToken", ErrInvalidToken, "INVALID_TOKEN", http.StatusUnauthorized},
		{"ErrInvalidTokenClaims", ErrInvalidTokenClaims, "INVALID_TOKEN_CLAIMS", http.StatusUnauthorized},
		{"ErrEmailRequired", ErrEmailRequired, "EMAIL_REQUIRED", http.StatusBadRequest},
		{"ErrPhoneRequired", ErrPhoneRequired, "PHONE_REQUIRED", http.StatusBadRequest},
		{"ErrAllFieldsRequired", ErrAllFieldsRequired, "ALL_FIELD_REQUIRED", http.StatusBadRequest},
		{"ErrEmailOtpRequired", ErrEmailOtpRequired, "EMAIL_OTP_REQUIRED", http.StatusBadRequest},
		{"ErrInvalidFileFormat", ErrInvalidFileFormat, "INVALID_FILE_FORMAT", http.StatusBadRequest},
		{"ErrFileSizeExceeded", ErrFileSizeExceeded, "FILE_SIZE_EXCEEDED", http.StatusBadRequest},
		{"ErrFailedParseMultipart", ErrFailedParseMultipart, "FAILED_PARSE_MULTIPART", http.StatusBadRequest},
		{"ErrFetchFailed", ErrFetchFailed, "FETCH_FAILED", http.StatusInternalServerError},
		{"ErrInvalidId", ErrInvalidId, "INVALID_ID", http.StatusBadRequest},
		{"ErrEncryptionFailed", ErrEncryptionFailed, "ENCRYPTION_FAILED", http.StatusInternalServerError},
		{"ErrDecryptionFailed", ErrDecryptionFailed, "DECRYPTION_FAILED", http.StatusInternalServerError},
		{"ErrDatabaseOperation", ErrDatabaseOperation, "DATABASE_ERROR", http.StatusInternalServerError},
		{"ErrEmailDeliveryFailed", ErrEmailDeliveryFailed, "EMAIL_DELIVERY_FAILED", http.StatusInternalServerError},
		{"ErrCloudinaryUploadFailed", ErrCloudinaryUploadFailed, "CLOUDINARY_UPLOAD_FAILED", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Expected code '%v', got %v", tt.code, tt.err.Code)
			}
			if tt.err.Status != tt.status {
				t.Errorf("Expected status %v, got %v", tt.status, tt.err.Status)
			}
			if tt.err.Message == "" {
				t.Error("Expected non-empty message")
			}
		})
	}
}