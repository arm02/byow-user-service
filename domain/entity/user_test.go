package entity

import (
	"testing"
	"time"
)

func TestUserStruct(t *testing.T) {
	now := time.Now()
	user := User{
		ID:           "user123",
		Fullname:     "John Doe",
		Email:        "john@example.com",
		Password:     "hashedpassword",
		PhoneNumber:  "+1234567890",
		AvatarUrl:    "https://example.com/avatar.jpg",
		OnBoarded:    true,
		OTP:          "123456",
		OTPType:      "email",
		OTPExpiresAt: now.Add(time.Minute * 5),
		Verified:     true,
		CreatedAt:    now,
	}

	// Test that all fields are properly set
	if user.ID != "user123" {
		t.Errorf("Expected ID 'user123', got %v", user.ID)
	}

	if user.Fullname != "John Doe" {
		t.Errorf("Expected Fullname 'John Doe', got %v", user.Fullname)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got %v", user.Email)
	}

	if user.Password != "hashedpassword" {
		t.Errorf("Expected Password 'hashedpassword', got %v", user.Password)
	}

	if user.PhoneNumber != "+1234567890" {
		t.Errorf("Expected PhoneNumber '+1234567890', got %v", user.PhoneNumber)
	}

	if user.AvatarUrl != "https://example.com/avatar.jpg" {
		t.Errorf("Expected AvatarUrl 'https://example.com/avatar.jpg', got %v", user.AvatarUrl)
	}

	if !user.OnBoarded {
		t.Error("Expected OnBoarded to be true")
	}

	if user.OTP != "123456" {
		t.Errorf("Expected OTP '123456', got %v", user.OTP)
	}

	if user.OTPType != "email" {
		t.Errorf("Expected OTPType 'email', got %v", user.OTPType)
	}

	if !user.Verified {
		t.Error("Expected Verified to be true")
	}

	if user.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, user.CreatedAt)
	}

	// Test that OTPExpiresAt is approximately correct (within 1 second)
	expectedExpiry := now.Add(time.Minute * 5)
	if user.OTPExpiresAt.Sub(expectedExpiry).Abs() > time.Second {
		t.Errorf("Expected OTPExpiresAt around %v, got %v", expectedExpiry, user.OTPExpiresAt)
	}
}

func TestUserStructZeroValues(t *testing.T) {
	var user User

	// Test zero values
	if user.ID != "" {
		t.Errorf("Expected empty ID, got %v", user.ID)
	}

	if user.Fullname != "" {
		t.Errorf("Expected empty Fullname, got %v", user.Fullname)
	}

	if user.Email != "" {
		t.Errorf("Expected empty Email, got %v", user.Email)
	}

	if user.Password != "" {
		t.Errorf("Expected empty Password, got %v", user.Password)
	}

	if user.PhoneNumber != "" {
		t.Errorf("Expected empty PhoneNumber, got %v", user.PhoneNumber)
	}

	if user.AvatarUrl != "" {
		t.Errorf("Expected empty AvatarUrl, got %v", user.AvatarUrl)
	}

	if user.OnBoarded {
		t.Error("Expected OnBoarded to be false")
	}

	if user.OTP != "" {
		t.Errorf("Expected empty OTP, got %v", user.OTP)
	}

	if user.OTPType != "" {
		t.Errorf("Expected empty OTPType, got %v", user.OTPType)
	}

	if !user.OTPExpiresAt.IsZero() {
		t.Errorf("Expected zero OTPExpiresAt, got %v", user.OTPExpiresAt)
	}

	if user.Verified {
		t.Error("Expected Verified to be false")
	}

	if !user.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", user.CreatedAt)
	}
}

func TestUserStructPartialValues(t *testing.T) {
	now := time.Now()
	user := User{
		Email:     "test@example.com",
		Verified:  false,
		CreatedAt: now,
	}

	// Test that specified fields are set
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got %v", user.Email)
	}

	if user.Verified {
		t.Error("Expected Verified to be false")
	}

	if user.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, user.CreatedAt)
	}

	// Test that unspecified fields have zero values
	if user.ID != "" {
		t.Errorf("Expected empty ID, got %v", user.ID)
	}

	if user.Fullname != "" {
		t.Errorf("Expected empty Fullname, got %v", user.Fullname)
	}

	if user.OnBoarded {
		t.Error("Expected OnBoarded to be false")
	}
}