package dto

import (
	"encoding/json"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	data := map[string]string{"message": "success"}
	successResp := SuccessResponse{
		Status: "SUCCESS",
		Code:   200,
		Data:   data,
	}

	if successResp.Status != "SUCCESS" {
		t.Errorf("Expected status 'SUCCESS', got %v", successResp.Status)
	}

	if successResp.Code != 200 {
		t.Errorf("Expected code 200, got %v", successResp.Code)
	}

	if successResp.Data == nil {
		t.Error("Expected non-nil data")
	}
}

func TestSuccessResponseJSON(t *testing.T) {
	data := map[string]string{"message": "success"}
	successResp := SuccessResponse{
		Status: "SUCCESS",
		Code:   200,
		Data:   data,
	}

	// Test marshaling
	jsonData, err := json.Marshal(successResp)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaled SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SuccessResponse: %v", err)
	}

	if unmarshaled.Status != successResp.Status {
		t.Errorf("Expected status %v, got %v", successResp.Status, unmarshaled.Status)
	}

	if unmarshaled.Code != successResp.Code {
		t.Errorf("Expected code %v, got %v", successResp.Code, unmarshaled.Code)
	}
}

func TestErrorResponse(t *testing.T) {
	errorData := ErrorResponseData{
		Message: "INTERNAL_SERVER_ERROR",
	}

	errorDetail := ErrorDetail{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: "Field validation error",
	}

	errorResp := ErrorResponse{
		Status: "ERROR",
		Code:   400,
		Data:   errorData,
		Error:  errorDetail,
	}

	if errorResp.Status != "ERROR" {
		t.Errorf("Expected status 'ERROR', got %v", errorResp.Status)
	}

	if errorResp.Code != 400 {
		t.Errorf("Expected code 400, got %v", errorResp.Code)
	}

	if errorResp.Data.Message != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected data message 'INTERNAL_SERVER_ERROR', got %v", errorResp.Data.Message)
	}

	if errorResp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected error code 'VALIDATION_ERROR', got %v", errorResp.Error.Code)
	}

	if errorResp.Error.Message != "Validation failed" {
		t.Errorf("Expected error message 'Validation failed', got %v", errorResp.Error.Message)
	}
}

func TestErrorResponseData(t *testing.T) {
	errorData := ErrorResponseData{
		Message: "Test error message",
	}

	if errorData.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got %v", errorData.Message)
	}
}

func TestErrorDetail(t *testing.T) {
	errorDetail := ErrorDetail{
		Code:    "TEST_ERROR",
		Message: "Test error occurred",
		Details: map[string]string{"field": "value"},
	}

	if errorDetail.Code != "TEST_ERROR" {
		t.Errorf("Expected code 'TEST_ERROR', got %v", errorDetail.Code)
	}

	if errorDetail.Message != "Test error occurred" {
		t.Errorf("Expected message 'Test error occurred', got %v", errorDetail.Message)
	}

	if errorDetail.Details == nil {
		t.Error("Expected non-nil details")
	}
}

func TestValidationErrorResponse(t *testing.T) {
	validationErrors := []ValidationError{
		{Field: "email", Message: "Invalid email format"},
		{Field: "password", Message: "Password too short"},
	}

	validationDetail := ValidationErrorDetail{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: validationErrors,
	}

	validationResp := ValidationErrorResponse{
		Status: "ERROR",
		Code:   400,
		Error:  validationDetail,
	}

	if validationResp.Status != "ERROR" {
		t.Errorf("Expected status 'ERROR', got %v", validationResp.Status)
	}

	if validationResp.Code != 400 {
		t.Errorf("Expected code 400, got %v", validationResp.Code)
	}

	if validationResp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected error code 'VALIDATION_ERROR', got %v", validationResp.Error.Code)
	}

	if len(validationResp.Error.Details) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(validationResp.Error.Details))
	}
}

func TestValidationErrorDetail(t *testing.T) {
	validationErrors := []ValidationError{
		{Field: "name", Message: "Name is required"},
	}

	validationDetail := ValidationErrorDetail{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: validationErrors,
	}

	if validationDetail.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected code 'VALIDATION_ERROR', got %v", validationDetail.Code)
	}

	if validationDetail.Message != "Validation failed" {
		t.Errorf("Expected message 'Validation failed', got %v", validationDetail.Message)
	}

	if len(validationDetail.Details) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(validationDetail.Details))
	}

	if validationDetail.Details[0].Field != "name" {
		t.Errorf("Expected field 'name', got %v", validationDetail.Details[0].Field)
	}

	if validationDetail.Details[0].Message != "Name is required" {
		t.Errorf("Expected message 'Name is required', got %v", validationDetail.Details[0].Message)
	}
}

func TestValidationError(t *testing.T) {
	validationError := ValidationError{
		Field:   "email",
		Message: "Invalid email format",
	}

	if validationError.Field != "email" {
		t.Errorf("Expected field 'email', got %v", validationError.Field)
	}

	if validationError.Message != "Invalid email format" {
		t.Errorf("Expected message 'Invalid email format', got %v", validationError.Message)
	}
}

func TestAllResponseDTOJSONSerialization(t *testing.T) {
	// Test all response DTO structs can be serialized to JSON
	structs := []interface{}{
		SuccessResponse{Status: "SUCCESS", Code: 200, Data: "test"},
		ErrorResponse{Status: "ERROR", Code: 400, Data: ErrorResponseData{Message: "error"}, Error: ErrorDetail{Code: "ERR", Message: "msg"}},
		ErrorResponseData{Message: "test message"},
		ErrorDetail{Code: "TEST", Message: "test", Details: "details"},
		ValidationErrorResponse{Status: "ERROR", Code: 400, Error: ValidationErrorDetail{Code: "VAL", Message: "msg", Details: []ValidationError{}}},
		ValidationErrorDetail{Code: "VAL", Message: "msg", Details: []ValidationError{{Field: "test", Message: "msg"}}},
		ValidationError{Field: "test", Message: "message"},
	}

	for i, s := range structs {
		_, err := json.Marshal(s)
		if err != nil {
			t.Errorf("Failed to marshal response struct at index %d: %v", i, err)
		}
	}
}

func TestEmptyResponses(t *testing.T) {
	// Test zero values
	var successResp SuccessResponse
	var errorResp ErrorResponse
	var validationResp ValidationErrorResponse

	if successResp.Status != "" || successResp.Code != 0 {
		t.Error("Expected zero values for SuccessResponse")
	}

	if errorResp.Status != "" || errorResp.Code != 0 {
		t.Error("Expected zero values for ErrorResponse")
	}

	if validationResp.Status != "" || validationResp.Code != 0 {
		t.Error("Expected zero values for ValidationErrorResponse")
	}
}