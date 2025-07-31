package constants

import "testing"

func TestResponseConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"SUCCESS constant", SUCCESS, "SUCCESS"},
		{"ERROR constant", ERROR, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.constant)
			}
		})
	}
}

func TestSuccessMessageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"LOGOUT_SUCCESSFUL", LOGOUT_SUCCESSFUL, "LOGOUT_SUCCESSFUL"},
		{"ONBOARD_SUCCESSFUL", ONBOARD_SUCCESSFUL, "ONBOARD_SUCCESSFUL"},
		{"PASSWORD_CHANGED_SUCCESS", PASSWORD_CHANGED_SUCCESS, "PASSWORD_CHANGED_SUCCESS"},
		{"EMAIL_CHANGED_SUCCESS", EMAIL_CHANGED_SUCCESS, "EMAIL_CHANGED_SUCCESS"},
		{"PHONE_CHANGED_SUCCESS", PHONE_CHANGED_SUCCESS, "PHONE_CHANGED_SUCCESS"},
		{"OTP_VERIFIED", OTP_VERIFIED, "OTP_VERIFIED"},
		{"OTP_SENT", OTP_SENT, "OTP_SENT"},
		{"VALID_TOKEN", VALID_TOKEN, "VALID_TOKEN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.constant)
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	if DefaultPageSize != 20 {
		t.Errorf("Expected DefaultPageSize to be 20, got %v", DefaultPageSize)
	}
}

func TestOTPTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"FORGOT_PASSWORD", FORGOT_PASSWORD, "forgot_password"},
		{"VERIFICATION", VERIFICATION, "verification"},
		{"EMAIL_CHANGED", EMAIL_CHANGED, "email_changed"},
		{"PASSWORD_CHANGED", PASSWORD_CHANGED, "password_changed"},
		{"PHONE_CHANGED", PHONE_CHANGED, "phone_changed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.constant)
			}
		})
	}
}

func TestConstantsAreStrings(t *testing.T) {
	// Test that all constants are non-empty strings
	constants := []string{
		SUCCESS,
		ERROR,
		LOGOUT_SUCCESSFUL,
		ONBOARD_SUCCESSFUL,
		PASSWORD_CHANGED_SUCCESS,
		EMAIL_CHANGED_SUCCESS,
		PHONE_CHANGED_SUCCESS,
		OTP_VERIFIED,
		OTP_SENT,
		VALID_TOKEN,
		FORGOT_PASSWORD,
		VERIFICATION,
		EMAIL_CHANGED,
		PASSWORD_CHANGED,
		PHONE_CHANGED,
	}

	for i, constant := range constants {
		if constant == "" {
			t.Errorf("Constant at index %d is empty", i)
		}
	}
}

func TestDefaultPageSizeType(t *testing.T) {
	// Test that DefaultPageSize is an int and positive
	if DefaultPageSize <= 0 {
		t.Errorf("Expected DefaultPageSize to be positive, got %v", DefaultPageSize)
	}

	// Test that it's within reasonable bounds
	if DefaultPageSize > 1000 {
		t.Errorf("DefaultPageSize seems too large: %v", DefaultPageSize)
	}
}