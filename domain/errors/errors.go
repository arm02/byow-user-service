package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error constructors for common scenarios
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    "UNAUTHORIZED",
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    "CONFLICT",
		Message: message,
		Status:  http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// Specific business logic errors - synced with constants.go format
var (
	// User authentication errors
	ErrUserNotFound           = &AppError{Code: "NOT_FOUND", Message: "User not found", Status: http.StatusNotFound}
	ErrInvalidCredentials     = &AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password", Status: http.StatusUnauthorized}
	ErrUserNotVerified        = &AppError{Code: "USER_NOT_VERIFIED", Message: "User account not verified", Status: http.StatusUnauthorized}
	ErrInvalidOldPassword     = &AppError{Code: "INVALID_OLD_PASSWORD", Message: "Invalid old password", Status: http.StatusBadRequest}
	
	// Registration errors
	ErrEmailAlreadyExists           = &AppError{Code: "EMAIL_ALREADY_REGISTERED", Message: "Email already registered", Status: http.StatusConflict}
	ErrPhoneAlreadyExists           = &AppError{Code: "PHONE_ALREADY_REGISTERED", Message: "Phone already registered", Status: http.StatusConflict}
	ErrEmailOrPhoneAlreadyRegistered = &AppError{Code: "EMAIL_OR_PHONE_ALREADY_REGISTERED", Message: "Email or phone already registered", Status: http.StatusConflict}
	
	// OTP errors
	ErrInvalidOTP             = &AppError{Code: "OTP_INVALID", Message: "Invalid OTP", Status: http.StatusBadRequest}
	ErrExpiredOTP             = &AppError{Code: "OTP_EXPIRED", Message: "OTP expired", Status: http.StatusBadRequest}
	
	// Token errors
	ErrInvalidToken           = &AppError{Code: "INVALID_TOKEN", Message: "Invalid or expired token", Status: http.StatusUnauthorized}
	ErrInvalidTokenClaims     = &AppError{Code: "INVALID_TOKEN_CLAIMS", Message: "Invalid token claims", Status: http.StatusUnauthorized}
	
	// Validation errors
	ErrEmailRequired          = &AppError{Code: "EMAIL_REQUIRED", Message: "Email is required", Status: http.StatusBadRequest}
	ErrPhoneRequired          = &AppError{Code: "PHONE_REQUIRED", Message: "Phone number is required", Status: http.StatusBadRequest}
	ErrAllFieldsRequired      = &AppError{Code: "ALL_FIELD_REQUIRED", Message: "All fields are required", Status: http.StatusBadRequest}
	ErrEmailOtpRequired       = &AppError{Code: "EMAIL_OTP_REQUIRED", Message: "Email and OTP are required", Status: http.StatusBadRequest}
	
	// File upload errors
	ErrInvalidFileFormat      = &AppError{Code: "INVALID_FILE_FORMAT", Message: "Invalid file format", Status: http.StatusBadRequest}
	ErrFileSizeExceeded       = &AppError{Code: "FILE_SIZE_EXCEEDED", Message: "File size exceeds limit", Status: http.StatusBadRequest}
	ErrFailedParseMultipart   = &AppError{Code: "FAILED_PARSE_MULTIPART", Message: "Failed to parse multipart form", Status: http.StatusBadRequest}
	
	// General errors
	ErrFetchFailed            = &AppError{Code: "FETCH_FAILED", Message: "Failed to fetch data", Status: http.StatusInternalServerError}
	ErrInvalidId              = &AppError{Code: "INVALID_ID", Message: "Invalid ID format", Status: http.StatusBadRequest}
	ErrEncryptionFailed       = &AppError{Code: "ENCRYPTION_FAILED", Message: "Encryption operation failed", Status: http.StatusInternalServerError}
	ErrDecryptionFailed       = &AppError{Code: "DECRYPTION_FAILED", Message: "Decryption operation failed", Status: http.StatusInternalServerError}
	ErrDatabaseOperation      = &AppError{Code: "DATABASE_ERROR", Message: "Database operation failed", Status: http.StatusInternalServerError}
	ErrEmailDeliveryFailed    = &AppError{Code: "EMAIL_DELIVERY_FAILED", Message: "Email delivery failed", Status: http.StatusInternalServerError}
	ErrCloudinaryUploadFailed = &AppError{Code: "CLOUDINARY_UPLOAD_FAILED", Message: "File upload failed", Status: http.StatusInternalServerError}
)

// Helper function to check if error is of specific type
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// Helper function to convert any error to AppError
func WrapError(err error, message string) *AppError {
	if appErr, ok := IsAppError(err); ok {
		return appErr
	}
	
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
		Details: err.Error(),
	}
}