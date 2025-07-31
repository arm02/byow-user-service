package mailer

import (
	"strings"
	"testing"

	"github.com/buildyow/byow-user-service/constants"
)

func TestGetOTPLifetime(t *testing.T) {
	tests := []struct {
		otpType  string
		expected int
	}{
		{constants.FORGOT_PASSWORD, 10},
		{constants.EMAIL_CHANGED, 10},
		{constants.PHONE_CHANGED, 10},
		{constants.VERIFICATION, 5},
		{"unknown_type", 1},
		{"", 1},
		{"invalid", 1},
	}

	for _, tt := range tests {
		t.Run(tt.otpType, func(t *testing.T) {
			result := getOTPLifetime(tt.otpType)
			if result != tt.expected {
				t.Errorf("getOTPLifetime(%v) = %v, want %v", tt.otpType, result, tt.expected)
			}
		})
	}
}

func TestSendOTP_InvalidSMTPSettings(t *testing.T) {
	// Test with invalid SMTP settings (should fail to connect)
	email := "test@example.com"
	otp := "123456"
	host := "invalid-smtp-host"
	user := "invalid-user"
	pass := "invalid-pass"
	port := 587
	otpType := constants.VERIFICATION

	err := SendOTP(email, otp, host, user, pass, port, otpType)

	// Should return error due to invalid SMTP settings
	if err == nil {
		t.Error("Expected error with invalid SMTP settings")
	}

	// Error should be related to connection/authentication
	if !strings.Contains(strings.ToLower(err.Error()), "dial") &&
		!strings.Contains(strings.ToLower(err.Error()), "connection") &&
		!strings.Contains(strings.ToLower(err.Error()), "resolve") {
		t.Logf("Error (expected): %v", err)
	}
}

func TestSendOTP_EmptyEmail(t *testing.T) {
	// Test with empty email
	email := ""
	otp := "123456"
	host := "smtp.gmail.com"
	user := "test@gmail.com"
	pass := "password"
	port := 587
	otpType := constants.VERIFICATION

	err := SendOTP(email, otp, host, user, pass, port, otpType)

	// Should return error due to empty email
	if err == nil {
		t.Error("Expected error with empty email")
	}
}

func TestSendOTP_InvalidPort(t *testing.T) {
	// Test with invalid port
	email := "test@example.com"
	otp := "123456"
	host := "smtp.gmail.com"
	user := "test@gmail.com"
	pass := "password"
	port := -1 // Invalid port
	otpType := constants.VERIFICATION

	err := SendOTP(email, otp, host, user, pass, port, otpType)

	// Should return error due to invalid port
	if err == nil {
		t.Error("Expected error with invalid port")
	}
}

func TestSendOTP_ZeroPort(t *testing.T) {
	// Test with zero port
	email := "test@example.com"
	otp := "123456"
	host := "smtp.gmail.com"
	user := "test@gmail.com"
	pass := "password"
	port := 0 // Zero port
	otpType := constants.VERIFICATION

	err := SendOTP(email, otp, host, user, pass, port, otpType)

	// Should return error due to zero port
	if err == nil {
		t.Error("Expected error with zero port")
	}
}

func TestSendOTP_EmptyOTP(t *testing.T) {
	// Test with empty OTP (should still try to send but fail due to invalid SMTP)
	email := "test@example.com"
	otp := ""
	host := "invalid-host"
	user := "test@gmail.com"
	pass := "password"
	port := 587
	otpType := constants.VERIFICATION

	err := SendOTP(email, otp, host, user, pass, port, otpType)

	// Should return error due to invalid host (not OTP validation)
	if err == nil {
		t.Error("Expected error due to invalid SMTP host")
	}
}

func TestSendOTP_MessageContent(t *testing.T) {
	// We can't easily test the actual sending without a real SMTP server,
	// but we can test that the function doesn't panic with various inputs
	// and that it attempts to create a proper message structure

	testCases := []struct {
		name    string
		email   string
		otp     string
		otpType string
	}{
		{"verification", "user@example.com", "123456", constants.VERIFICATION},
		{"forgot_password", "user@example.com", "789012", constants.FORGOT_PASSWORD},
		{"email_changed", "user@example.com", "345678", constants.EMAIL_CHANGED},
		{"phone_changed", "user@example.com", "901234", constants.PHONE_CHANGED},
		{"unknown_type", "user@example.com", "567890", "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use invalid SMTP settings so it fails to send but doesn't panic
			err := SendOTP(tc.email, tc.otp, "invalid-host", "user", "pass", 587, tc.otpType)
			
			// We expect an error due to invalid SMTP, but no panic
			if err == nil {
				t.Error("Expected error with invalid SMTP settings")
			}
			
			// Test completed without panic
			t.Logf("Test case %s completed with expected error: %v", tc.name, err)
		})
	}
}

func TestSendOTP_DifferentOTPTypes(t *testing.T) {
	// Test that different OTP types are handled and result in appropriate lifetimes
	// being referenced in the message (though we can't verify the message content directly)
	
	otpTypes := []string{
		constants.VERIFICATION,
		constants.FORGOT_PASSWORD,
		constants.EMAIL_CHANGED,
		constants.PHONE_CHANGED,
		"unknown_type",
	}

	for _, otpType := range otpTypes {
		t.Run(otpType, func(t *testing.T) {
			// Verify that getOTPLifetime works correctly for this type
			lifetime := getOTPLifetime(otpType)
			
			expectedLifetime := 1 // default
			switch otpType {
			case constants.FORGOT_PASSWORD, constants.EMAIL_CHANGED, constants.PHONE_CHANGED:
				expectedLifetime = 10
			case constants.VERIFICATION:
				expectedLifetime = 5
			}
			
			if lifetime != expectedLifetime {
				t.Errorf("Expected lifetime %d for %s, got %d", expectedLifetime, otpType, lifetime)
			}

			// Test SendOTP with this type (will fail due to invalid SMTP but shouldn't panic)
			err := SendOTP("test@example.com", "123456", "invalid", "user", "pass", 587, otpType)
			if err == nil {
				t.Error("Expected error with invalid SMTP")
			}
		})
	}
}

func TestSendOTP_InvalidEmailFormats(t *testing.T) {
	// Test with various invalid email formats
	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"test@",
		"test@@example.com",
		"test.example.com",
		" ",
	}

	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			err := SendOTP(email, "123456", "invalid-host", "user", "pass", 587, constants.VERIFICATION)
			
			// Should return error (either due to invalid email or invalid SMTP)
			if err == nil {
				t.Errorf("Expected error with invalid email: %s", email)
			}
			
			t.Logf("Invalid email %s resulted in expected error: %v", email, err)
		})
	}
}

func TestSendOTP_LongOTP(t *testing.T) {
	// Test with very long OTP
	longOTP := strings.Repeat("1234567890", 10) // 100 characters
	
	err := SendOTP("test@example.com", longOTP, "invalid-host", "user", "pass", 587, constants.VERIFICATION)
	
	// Should still attempt to send (and fail due to invalid SMTP)
	if err == nil {
		t.Error("Expected error with invalid SMTP settings")
	}
}

func TestSendOTP_SpecialCharactersInOTP(t *testing.T) {
	// Test with special characters in OTP
	specialOTPs := []string{
		"!@#$%^",
		"123-456",
		"abc123",
		"123.456",
		"123 456",
	}

	for _, otp := range specialOTPs {
		t.Run(otp, func(t *testing.T) {
			err := SendOTP("test@example.com", otp, "invalid-host", "user", "pass", 587, constants.VERIFICATION)
			
			// Should attempt to send regardless of OTP content
			if err == nil {
				t.Error("Expected error with invalid SMTP settings")
			}
		})
	}
}

func TestSendOTP_CommonPorts(t *testing.T) {
	// Test with common SMTP ports
	commonPorts := []int{25, 587, 465, 2525}
	
	for _, port := range commonPorts {
		t.Run(string(rune(port)), func(t *testing.T) {
			err := SendOTP("test@example.com", "123456", "invalid-host", "user", "pass", port, constants.VERIFICATION)
			
			// Should fail due to invalid host, not port
			if err == nil {
				t.Error("Expected error with invalid SMTP host")
			}
		})
	}
}

// Benchmark test
func BenchmarkGetOTPLifetime(b *testing.B) {
	otpTypes := []string{
		constants.VERIFICATION,
		constants.FORGOT_PASSWORD,
		constants.EMAIL_CHANGED,
		constants.PHONE_CHANGED,
		"unknown",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		otpType := otpTypes[i%len(otpTypes)]
		getOTPLifetime(otpType)
	}
}