package utils

import (
	"os"
	"testing"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
)

func TestEncryptDecrypt(t *testing.T) {
	// Set up test environment variable
	originalKey := os.Getenv("DECRYPT_KEY")
	testKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", testKey)
	defer os.Setenv("DECRYPT_KEY", originalKey)

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "unicode text",
			plaintext: "Hello 世界 🌍",
		},
		{
			name:      "long text",
			plaintext: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encryption
			encrypted, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			if encrypted == "" {
				t.Error("Expected non-empty encrypted string")
			}

			if encrypted == tt.plaintext {
				t.Error("Encrypted text should be different from plaintext")
			}

			// Test decryption
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDecryptConsistency(t *testing.T) {
	// Set up test environment variable
	originalKey := os.Getenv("DECRYPT_KEY")
	testKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", testKey)
	defer os.Setenv("DECRYPT_KEY", originalKey)

	plaintext := "test message"

	// Encrypt the same message multiple times
	encrypted1, err1 := Encrypt(plaintext)
	encrypted2, err2 := Encrypt(plaintext)

	if err1 != nil || err2 != nil {
		t.Fatalf("Encrypt() errors: %v, %v", err1, err2)
	}

	// Encrypted values should be different due to random nonce
	if encrypted1 == encrypted2 {
		t.Error("Multiple encryptions of same plaintext should produce different ciphertexts")
	}

	// But both should decrypt to the same plaintext
	decrypted1, err1 := Decrypt(encrypted1)
	decrypted2, err2 := Decrypt(encrypted2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Decrypt() errors: %v, %v", err1, err2)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Errorf("Decryption failed: got %v and %v, want %v", decrypted1, decrypted2, plaintext)
	}
}

func TestEncryptWithInvalidKey(t *testing.T) {
	// Set up test environment variable with invalid key
	originalKey := os.Getenv("DECRYPT_KEY")
	os.Setenv("DECRYPT_KEY", "short")
	defer os.Setenv("DECRYPT_KEY", originalKey)

	_, err := Encrypt("test message")
	if err == nil {
		t.Error("Expected error with invalid key length")
	}
}

func TestDecryptWithInvalidKey(t *testing.T) {
	// First encrypt with valid key
	originalKey := os.Getenv("DECRYPT_KEY")
	validKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", validKey)

	encrypted, err := Encrypt("test message")
	if err != nil {
		t.Fatalf("Setup encryption failed: %v", err)
	}

	// Then try to decrypt with invalid key
	os.Setenv("DECRYPT_KEY", "short")
	defer os.Setenv("DECRYPT_KEY", originalKey)

	_, err = Decrypt(encrypted)
	if err == nil {
		t.Error("Expected error with invalid key length")
	}
}

func TestDecryptWithInvalidBase64(t *testing.T) {
	// Set up test environment variable
	originalKey := os.Getenv("DECRYPT_KEY")
	testKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", testKey)
	defer os.Setenv("DECRYPT_KEY", originalKey)

	invalidBase64 := "invalid-base64!"
	_, err := Decrypt(invalidBase64)
	if err == nil {
		t.Error("Expected error with invalid base64 input")
	}
}

func TestDecryptWithTooShortCiphertext(t *testing.T) {
	// Set up test environment variable
	originalKey := os.Getenv("DECRYPT_KEY")
	testKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", testKey)
	defer os.Setenv("DECRYPT_KEY", originalKey)

	// Create a valid base64 string that's too short
	shortCiphertext := "YWJj" // "abc" in base64, which is too short for GCM nonce
	_, err := Decrypt(shortCiphertext)
	if err == nil {
		t.Error("Expected error with too short ciphertext")
	}

	if err != appErrors.ErrDecryptionFailed {
		t.Errorf("Expected ErrDecryptionFailed, got %v", err)
	}
}

func TestDecryptWithCorruptedCiphertext(t *testing.T) {
	// Set up test environment variable
	originalKey := os.Getenv("DECRYPT_KEY")
	testKey := "12345678901234567890123456789012" // Exactly 32 bytes
	os.Setenv("DECRYPT_KEY", testKey)
	defer os.Setenv("DECRYPT_KEY", originalKey)

	// First encrypt a message
	plaintext := "test message"
	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Setup encryption failed: %v", err)
	}

	// Corrupt the encrypted message by changing one character
	corrupted := encrypted[:len(encrypted)-4] + "XXXX"
	
	_, err = Decrypt(corrupted)
	if err == nil {
		t.Error("Expected error with corrupted ciphertext")
	}
}

func TestEncryptDecryptWithEmptyKey(t *testing.T) {
	// Set up test environment variable with empty key
	originalKey := os.Getenv("DECRYPT_KEY")
	os.Setenv("DECRYPT_KEY", "")
	defer os.Setenv("DECRYPT_KEY", originalKey)

	_, err := Encrypt("test message")
	if err == nil {
		t.Error("Expected error with empty key")
	}

	_, err = Decrypt("dGVzdA==") // "test" in base64
	if err == nil {
		t.Error("Expected error with empty key")
	}
}

func TestEncryptDecryptWithMissingKey(t *testing.T) {
	// Set up test environment variable with missing key
	originalKey := os.Getenv("DECRYPT_KEY")
	os.Unsetenv("DECRYPT_KEY")
	defer os.Setenv("DECRYPT_KEY", originalKey)

	_, err := Encrypt("test message")
	if err == nil {
		t.Error("Expected error with missing key")
	}

	_, err = Decrypt("dGVzdA==") // "test" in base64
	if err == nil {
		t.Error("Expected error with missing key")
	}
}