package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	userID := "user123"
	email := "test@example.com"
	phone := "+1234567890"
	secret := "test-secret-key"
	minutes := 30

	token, err := GenerateToken(userID, email, phone, secret, minutes)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Verify token structure (should have 3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected token to have 3 parts, got %d", len(parts))
	}
}

func TestGenerateTokenClaims(t *testing.T) {
	userID := "user123"
	email := "test@example.com"
	phone := "+1234567890"
	secret := "test-secret-key"
	minutes := 30

	token, err := GenerateToken(userID, email, phone, secret, minutes)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Parse and verify claims
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Expected MapClaims")
	}

	// Verify claims
	if claims["user_id"] != userID {
		t.Errorf("Expected user_id %v, got %v", userID, claims["user_id"])
	}

	if claims["email"] != email {
		t.Errorf("Expected email %v, got %v", email, claims["email"])
	}

	if claims["phone"] != phone {
		t.Errorf("Expected phone %v, got %v", phone, claims["phone"])
	}

	if claims["iss"] != "byow-user-service" {
		t.Errorf("Expected iss 'byow-user-service', got %v", claims["iss"])
	}

	if claims["aud"] != "byow-platform" {
		t.Errorf("Expected aud 'byow-platform', got %v", claims["aud"])
	}

	// Verify JTI exists and is non-empty
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		t.Error("Expected non-empty JTI")
	}

	// Verify issued at time
	iat, ok := claims["iat"].(float64)
	if !ok {
		t.Error("Expected iat claim")
	}

	issuedAt := time.Unix(int64(iat), 0)
	if time.Since(issuedAt) > time.Minute {
		t.Error("Token issued at time seems incorrect")
	}

	// Verify expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("Expected exp claim")
	}

	expiresAt := time.Unix(int64(exp), 0)
	expectedExpiry := issuedAt.Add(time.Duration(minutes) * time.Minute)
	
	// Allow for small time differences (within 5 seconds)
	if expiresAt.Sub(expectedExpiry).Abs() > 5*time.Second {
		t.Errorf("Expected expiry %v, got %v", expectedExpiry, expiresAt)
	}
}

func TestGenerateTokenWithDifferentExpiry(t *testing.T) {
	userID := "user123"
	email := "test@example.com"
	phone := "+1234567890"
	secret := "test-secret-key"

	tests := []int{1, 15, 60, 120, 1440} // Different minute values

	for _, minutes := range tests {
		t.Run(string(rune(minutes)), func(t *testing.T) {
			token, err := GenerateToken(userID, email, phone, secret, minutes)
			if err != nil {
				t.Fatalf("GenerateToken() error = %v", err)
			}

			// Parse and verify expiration
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				t.Fatalf("Failed to parse token: %v", err)
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Expected MapClaims")
			}

			exp, ok := claims["exp"].(float64)
			if !ok {
				t.Error("Expected exp claim")
			}

			iat, ok := claims["iat"].(float64)
			if !ok {
				t.Error("Expected iat claim")
			}

			actualDuration := time.Unix(int64(exp), 0).Sub(time.Unix(int64(iat), 0))
			expectedDuration := time.Duration(minutes) * time.Minute

			if actualDuration != expectedDuration {
				t.Errorf("Expected duration %v, got %v", expectedDuration, actualDuration)
			}
		})
	}
}

func TestGenerateTokenWithEmptySecret(t *testing.T) {
	userID := "user123"
	email := "test@example.com"
	phone := "+1234567890"
	secret := ""
	minutes := 30

	token, err := GenerateToken(userID, email, phone, secret, minutes)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Token should still be generated but will fail validation with empty secret
	if token == "" {
		t.Error("Expected non-empty token even with empty secret")
	}
}

func TestGenerateTokenWithEmptyInputs(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		email   string
		phone   string
		secret  string
		minutes int
	}{
		{
			name:    "empty user ID",
			userID:  "",
			email:   "test@example.com", 
			phone:   "+1234567890",
			secret:  "secret",
			minutes: 30,
		},
		{
			name:    "empty email",
			userID:  "user123",
			email:   "",
			phone:   "+1234567890",
			secret:  "secret",
			minutes: 30,
		},
		{
			name:    "empty phone",
			userID:  "user123",
			email:   "test@example.com",
			phone:   "",
			secret:  "secret",
			minutes: 30,
		},
		{
			name:    "zero minutes",
			userID:  "user123",
			email:   "test@example.com",
			phone:   "+1234567890",
			secret:  "secret",
			minutes: 0,
		},
		{
			name:    "negative minutes",
			userID:  "user123",
			email:   "test@example.com",
			phone:   "+1234567890",
			secret:  "secret",
			minutes: -30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.phone, tt.secret, tt.minutes)
			if err != nil {
				t.Fatalf("GenerateToken() error = %v", err)
			}

			// Should still generate token, but claims will contain empty values
			if token == "" {
				t.Error("Expected non-empty token")
			}

			// Verify the claims contain the provided values (even if empty)
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.secret), nil
			})
			// For zero/negative minutes or empty secret, we expect parsing errors
			if err != nil && tt.secret != "" && tt.minutes > 0 {
				t.Fatalf("Failed to parse token: %v", err)
			}

			if tt.secret != "" {
				// For zero or negative minutes, token will be expired, so skip validation
				if tt.minutes > 0 && parsedToken.Valid {
					claims, ok := parsedToken.Claims.(jwt.MapClaims)
					if !ok {
						t.Fatal("Expected MapClaims")
					}

					if claims["user_id"] != tt.userID {
						t.Errorf("Expected user_id %v, got %v", tt.userID, claims["user_id"])
					}

					if claims["email"] != tt.email {
						t.Errorf("Expected email %v, got %v", tt.email, claims["email"])
					}

					if claims["phone"] != tt.phone {
						t.Errorf("Expected phone %v, got %v", tt.phone, claims["phone"])
					}
				} else if tt.minutes <= 0 && parsedToken.Valid {
					t.Error("Expected expired token for zero/negative minutes")
				}
			}
		})
	}
}

func TestGenerateJTI(t *testing.T) {
	// Test generating multiple JTIs
	jti1, err := generateJTI()
	if err != nil {
		t.Fatalf("generateJTI() error = %v", err)
	}

	jti2, err := generateJTI()
	if err != nil {
		t.Fatalf("generateJTI() error = %v", err)
	}

	// JTIs should be non-empty
	if jti1 == "" || jti2 == "" {
		t.Error("Expected non-empty JTIs")
	}

	// JTIs should be different
	if jti1 == jti2 {
		t.Error("Expected different JTIs")
	}

	// JTI should be 32 characters long (16 bytes in hex)
	if len(jti1) != 32 || len(jti2) != 32 {
		t.Errorf("Expected JTI length 32, got %d and %d", len(jti1), len(jti2))
	}

	// JTI should only contain hex characters
	for _, char := range jti1 {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("JTI contains invalid hex character: %c", char)
		}
	}
}

func TestGenerateTokenMultipleTimes(t *testing.T) {
	userID := "user123"
	email := "test@example.com"
	phone := "+1234567890"
	secret := "test-secret-key"
	minutes := 30

	// Generate multiple tokens with same parameters
	tokens := make([]string, 10)
	for i := 0; i < 10; i++ {
		token, err := GenerateToken(userID, email, phone, secret, minutes)
		if err != nil {
			t.Fatalf("GenerateToken() error = %v", err)
		}
		tokens[i] = token
	}

	// All tokens should be different due to different JTI and iat values
	for i := 0; i < len(tokens); i++ {
		for j := i + 1; j < len(tokens); j++ {
			if tokens[i] == tokens[j] {
				t.Errorf("Expected different tokens, but got identical at positions %d and %d", i, j)
			}
		}
	}
}

func TestGenerateJTIError(t *testing.T) {
	// This test is mainly for coverage completion
	// Under normal circumstances, generateJTI should not fail
	// But we can test that it generates proper hex strings
	
	for i := 0; i < 100; i++ {
		jti, err := generateJTI()
		if err != nil {
			t.Fatalf("generateJTI() unexpected error = %v", err)
		}
		
		if len(jti) != 32 {
			t.Errorf("Expected JTI length 32, got %d", len(jti))
		}
		
		// Verify it's valid hex
		for _, char := range jti {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("JTI contains invalid hex character: %c", char)
			}
		}
	}
}