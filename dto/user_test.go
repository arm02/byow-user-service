package dto

import (
	"encoding/json"
	"testing"
)

func TestLoginRequest(t *testing.T) {
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if loginReq.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %v", loginReq.Email)
	}

	if loginReq.Password != "password123" {
		t.Errorf("Expected password 'password123', got %v", loginReq.Password)
	}
}

func TestLoginRequestJSON(t *testing.T) {
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Test marshaling
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		t.Fatalf("Failed to marshal LoginRequest: %v", err)
	}

	// Test unmarshaling
	var unmarshaled LoginRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal LoginRequest: %v", err)
	}

	if unmarshaled.Email != loginReq.Email {
		t.Errorf("Expected email %v, got %v", loginReq.Email, unmarshaled.Email)
	}

	if unmarshaled.Password != loginReq.Password {
		t.Errorf("Expected password %v, got %v", loginReq.Password, unmarshaled.Password)
	}
}

func TestRegisterRequest(t *testing.T) {
	registerReq := RegisterRequest{
		Fullname:    "John Doe",
		Email:       "john@example.com",
		Password:    "password123",
		PhoneNumber: "+1234567890",
		AvatarUrl:   "https://example.com/avatar.jpg",
	}

	if registerReq.Fullname != "John Doe" {
		t.Errorf("Expected fullname 'John Doe', got %v", registerReq.Fullname)
	}

	if registerReq.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", registerReq.Email)
	}

	if registerReq.Password != "password123" {
		t.Errorf("Expected password 'password123', got %v", registerReq.Password)
	}

	if registerReq.PhoneNumber != "+1234567890" {
		t.Errorf("Expected phone '+1234567890', got %v", registerReq.PhoneNumber)
	}

	if registerReq.AvatarUrl != "https://example.com/avatar.jpg" {
		t.Errorf("Expected avatar URL 'https://example.com/avatar.jpg', got %v", registerReq.AvatarUrl)
	}
}

func TestUserResponse(t *testing.T) {
	userResp := UserResponse{
		Fullname:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		AvatarUrl:   "https://example.com/avatar.jpg",
		Verified:    true,
		OnBoarded:   false,
		Token:       "sample-token",
		CreatedAt:   "2024-01-15T10:30:00Z",
	}

	if userResp.Fullname != "John Doe" {
		t.Errorf("Expected fullname 'John Doe', got %v", userResp.Fullname)
	}

	if userResp.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", userResp.Email)
	}

	if !userResp.Verified {
		t.Error("Expected Verified to be true")
	}

	if userResp.OnBoarded {
		t.Error("Expected OnBoarded to be false")
	}

	if userResp.Token != "sample-token" {
		t.Errorf("Expected token 'sample-token', got %v", userResp.Token)
	}
}

func TestUserResponseSwagger(t *testing.T) {
	userResp := UserResponse{
		Fullname: "John Doe",
		Email:    "john@example.com",
	}

	swaggerResp := UserResponseSwagger{
		Status: "SUCCESS",
		Code:   200,
		Data:   userResp,
	}

	if swaggerResp.Status != "SUCCESS" {
		t.Errorf("Expected status 'SUCCESS', got %v", swaggerResp.Status)
	}

	if swaggerResp.Code != 200 {
		t.Errorf("Expected code 200, got %v", swaggerResp.Code)
	}

	if swaggerResp.Data.Fullname != "John Doe" {
		t.Errorf("Expected data fullname 'John Doe', got %v", swaggerResp.Data.Fullname)
	}
}

func TestVerifyOTPRequest(t *testing.T) {
	otpReq := VerifyOTPRequest{
		Email: "john@example.com",
		OTP:   "123456",
	}

	if otpReq.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", otpReq.Email)
	}

	if otpReq.OTP != "123456" {
		t.Errorf("Expected OTP '123456', got %v", otpReq.OTP)
	}
}

func TestChangePasswordRequest(t *testing.T) {
	changeReq := ChangePasswordRequest{
		Email:    "john@example.com",
		OTP:      "123456",
		Password: "newpassword",
	}

	if changeReq.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", changeReq.Email)
	}

	if changeReq.OTP != "123456" {
		t.Errorf("Expected OTP '123456', got %v", changeReq.OTP)
	}

	if changeReq.Password != "newpassword" {
		t.Errorf("Expected password 'newpassword', got %v", changeReq.Password)
	}
}

func TestChangePasswordWithOldPasswordRequest(t *testing.T) {
	changeReq := ChangePasswordWithOldPasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	if changeReq.OldPassword != "oldpassword" {
		t.Errorf("Expected old password 'oldpassword', got %v", changeReq.OldPassword)
	}

	if changeReq.NewPassword != "newpassword" {
		t.Errorf("Expected new password 'newpassword', got %v", changeReq.NewPassword)
	}
}

func TestChangeEmailRequest(t *testing.T) {
	changeReq := ChangeEmailRequest{
		NewEmail: "newemail@example.com",
		OTP:      "123456",
	}

	if changeReq.NewEmail != "newemail@example.com" {
		t.Errorf("Expected new email 'newemail@example.com', got %v", changeReq.NewEmail)
	}

	if changeReq.OTP != "123456" {
		t.Errorf("Expected OTP '123456', got %v", changeReq.OTP)
	}
}

func TestChangePhoneRequest(t *testing.T) {
	changeReq := ChangePhoneRequest{
		NewPhone: "+9876543210",
		OTP:      "123456",
	}

	if changeReq.NewPhone != "+9876543210" {
		t.Errorf("Expected new phone '+9876543210', got %v", changeReq.NewPhone)
	}

	if changeReq.OTP != "123456" {
		t.Errorf("Expected OTP '123456', got %v", changeReq.OTP)
	}
}

func TestAllDTOJSONSerialization(t *testing.T) {
	// Test all DTO structs can be serialized to JSON
	structs := []interface{}{
		LoginRequest{Email: "test@example.com", Password: "pass"},
		RegisterRequest{Fullname: "John", Email: "john@example.com", Password: "pass", PhoneNumber: "+123"},
		UserResponse{Fullname: "John", Email: "john@example.com", Verified: true},
		UserResponseSwagger{Status: "SUCCESS", Code: 200, Data: UserResponse{}},
		VerifyOTPRequest{Email: "test@example.com", OTP: "123456"},
		ChangePasswordRequest{Email: "test@example.com", OTP: "123456", Password: "pass"},
		ChangePasswordWithOldPasswordRequest{OldPassword: "old", NewPassword: "new"},
		ChangeEmailRequest{NewEmail: "new@example.com", OTP: "123456"},
		ChangePhoneRequest{NewPhone: "+987654321", OTP: "123456"},
	}

	for i, s := range structs {
		_, err := json.Marshal(s)
		if err != nil {
			t.Errorf("Failed to marshal struct at index %d: %v", i, err)
		}
	}
}