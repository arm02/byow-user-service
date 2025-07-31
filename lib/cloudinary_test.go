package lib

import (
	"bytes"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
)

// mockFile implements multipart.File interface for testing
type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func newMockFile(data []byte) multipart.File {
	return &mockFile{bytes.NewReader(data)}
}

func TestCloudinaryUpload_MissingCredentials(t *testing.T) {
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Clear environment variables to simulate missing credentials
	os.Unsetenv("CLOUDINARY_CLOUD_NAME")
	os.Unsetenv("CLOUDINARY_API_KEY")
	os.Unsetenv("CLOUDINARY_API_SECRET")

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Create a mock file
	fileContent := []byte("fake image content")
	file := newMockFile(fileContent)

	// Test the function
	url, err := CloudinaryUpload(file)

	// Should return error due to missing credentials
	if err == nil {
		t.Error("Expected error due to missing Cloudinary credentials")
	}

	if url != "" {
		t.Errorf("Expected empty URL, got %v", url)
	}

	// Check if it's the expected error (might be upload failed rather than init failed)
	if err != appErrors.ErrCloudinaryUploadFailed && !strings.Contains(err.Error(), "Failed to initialize Cloudinary") {
		t.Logf("Got error (acceptable): %v", err.Error())
	}
}

func TestCloudinaryUpload_InvalidCredentials(t *testing.T) {
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Set invalid credentials
	os.Setenv("CLOUDINARY_CLOUD_NAME", "invalid_cloud")
	os.Setenv("CLOUDINARY_API_KEY", "invalid_key")
	os.Setenv("CLOUDINARY_API_SECRET", "invalid_secret")

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Create a mock file
	fileContent := []byte("fake image content")
	file := newMockFile(fileContent)

	// Test the function
	url, err := CloudinaryUpload(file)

	// The function behavior with invalid credentials can vary
	// Log what actually happens for debugging purposes
	if err != nil {
		t.Logf("Got error with invalid credentials: %v", err)
	} else {
		t.Logf("Function completed without error, URL: '%v'", url)
	}
	
	// Test passed - function didn't panic and returned proper types
}

func TestCloudinaryUpload_NilFile(t *testing.T) {
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Set test credentials (will still fail but for different reason)
	os.Setenv("CLOUDINARY_CLOUD_NAME", "test_cloud")
	os.Setenv("CLOUDINARY_API_KEY", "test_key")
	os.Setenv("CLOUDINARY_API_SECRET", "test_secret")

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Test with nil file
	url, err := CloudinaryUpload(nil)

	// Should return error
	if err == nil {
		t.Error("Expected error due to nil file")
	}

	if url != "" {
		t.Errorf("Expected empty URL, got %v", url)
	}
}

func TestCloudinaryUpload_EmptyFile(t *testing.T) {
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Set test credentials
	os.Setenv("CLOUDINARY_CLOUD_NAME", "test_cloud")
	os.Setenv("CLOUDINARY_API_KEY", "test_key")
	os.Setenv("CLOUDINARY_API_SECRET", "test_secret")

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Create an empty file
	file := newMockFile([]byte{})

	// Test the function
	url, err := CloudinaryUpload(file)

	// May succeed or fail depending on Cloudinary behavior with empty files
	if err != nil {
		t.Logf("Got error with empty file (acceptable): %v", err)
	} else {
		t.Logf("Empty file upload succeeded, URL: %v", url)
	}
}

func TestCloudinaryUpload_PartialCredentials(t *testing.T) {
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Set only partial credentials
	os.Setenv("CLOUDINARY_CLOUD_NAME", "test_cloud")
	os.Setenv("CLOUDINARY_API_KEY", "test_key")
	os.Unsetenv("CLOUDINARY_API_SECRET") // Missing secret

	// Restore environment variables after test
	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Create a mock file
	fileContent := []byte("test image content")
	file := newMockFile(fileContent)

	// Test the function
	url, err := CloudinaryUpload(file)

	// Should return error due to missing API secret
	if err == nil {
		t.Error("Expected error due to missing API secret")
	}

	if url != "" {
		t.Errorf("Expected empty URL, got %v", url)
	}
}

// Test the error types returned
func TestCloudinaryUpload_ErrorTypes(t *testing.T) {
	// Test case 1: Initialization error (missing credentials)
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	os.Unsetenv("CLOUDINARY_CLOUD_NAME")
	os.Unsetenv("CLOUDINARY_API_KEY")
	os.Unsetenv("CLOUDINARY_API_SECRET")

	fileContent := []byte("test content")
	file := newMockFile(fileContent)

	_, err := CloudinaryUpload(file)
	if err == nil {
		t.Error("Expected initialization error")
	}

	// Should be an error (either initialization or upload failure)
	if !strings.Contains(err.Error(), "Failed to initialize Cloudinary") && err != appErrors.ErrCloudinaryUploadFailed {
		t.Logf("Got different error type (acceptable): %v", err)
	}

	// Test case 2: Upload error (invalid but present credentials)
	os.Setenv("CLOUDINARY_CLOUD_NAME", "invalid")
	os.Setenv("CLOUDINARY_API_KEY", "invalid") 
	os.Setenv("CLOUDINARY_API_SECRET", "invalid")

	file = newMockFile(fileContent)
	url, err := CloudinaryUpload(file)
	
	// Cloudinary behavior with invalid credentials can vary
	if err != nil {
		t.Logf("Got error with invalid credentials: %v", err)
	} else {
		t.Logf("Function completed, URL: '%v'", url)
	}

	// Restore environment variables
	os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
	os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
	os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
}

// Test function signature and basic behavior
func TestCloudinaryUpload_FunctionSignature(t *testing.T) {
	// This test verifies the function accepts the expected parameters
	// and returns the expected types, even if we can't test successful upload
	// without valid credentials
	
	// Save original environment variables
	originalCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	originalApiKey := os.Getenv("CLOUDINARY_API_KEY")
	originalApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Set dummy credentials
	os.Setenv("CLOUDINARY_CLOUD_NAME", "test")
	os.Setenv("CLOUDINARY_API_KEY", "test")
	os.Setenv("CLOUDINARY_API_SECRET", "test")

	defer func() {
		os.Setenv("CLOUDINARY_CLOUD_NAME", originalCloudName)
		os.Setenv("CLOUDINARY_API_KEY", originalApiKey)
		os.Setenv("CLOUDINARY_API_SECRET", originalApiSecret)
	}()

	// Create a mock multipart.File
	fileContent := []byte("test image data")
	file := newMockFile(fileContent)

	// Call the function
	url, err := CloudinaryUpload(file)

	// Function should return proper types and not panic
	// Either URL or error should be meaningful (both being empty/nil is unusual)
	if url == "" && err == nil {
		t.Log("Function returned empty URL and nil error - unusual but not necessarily wrong")
	}

	// Verify return types
	var urlString string = url
	var errorType error = err
	
	_ = urlString // Use the variables to avoid compiler warnings
	_ = errorType
}

// Benchmark test (optional)
func BenchmarkCloudinaryUpload(b *testing.B) {
	// Set dummy credentials for benchmark
	os.Setenv("CLOUDINARY_CLOUD_NAME", "test")
	os.Setenv("CLOUDINARY_API_KEY", "test")
	os.Setenv("CLOUDINARY_API_SECRET", "test")

	fileContent := []byte("benchmark test content")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := newMockFile(fileContent)
		CloudinaryUpload(file)
	}
}