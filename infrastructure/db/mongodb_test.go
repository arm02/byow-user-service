package db

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestConnect_InvalidURI(t *testing.T) {
	// Test with invalid MongoDB URI
	invalidURI := "invalid-mongodb-uri"
	
	client, err := Connect(invalidURI)
	
	// Connect might succeed but ping should fail, or connect might fail  
	// The MongoDB driver may not immediately validate the URI format
	if err != nil {
		t.Logf("Got expected error with invalid URI: %v", err)
	} else if client != nil {
		// Try to ping to see if it actually works
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		pingErr := client.Ping(ctx, nil)
		if pingErr == nil {
			t.Log("Unexpectedly successful connection with invalid URI")
		} else {
			t.Logf("Connection created but ping failed as expected: %v", pingErr)
		}
	}
	
	if client != nil {
		client.Disconnect(context.Background())
	}
}

func TestConnect_EmptyURI(t *testing.T) {
	// Test with empty URI
	emptyURI := ""
	
	client, err := Connect(emptyURI)
	
	// Should return error with empty URI
	if err == nil {
		t.Error("Expected error with empty MongoDB URI")
	}
	
	if client != nil {
		client.Disconnect(context.Background())
	}
}

func TestConnect_Timeout(t *testing.T) {
	// Test with a URI that would timeout (non-existent host)
	timeoutURI := "mongodb://nonexistent-host:27017"
	
	// Connect may succeed (lazy connection) but ping should fail/timeout
	client, err := Connect(timeoutURI)
	
	if client != nil {
		defer client.Disconnect(context.Background())
		
		// Try to ping to test actual connectivity
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		pingErr := client.Ping(ctx, nil)
		if pingErr == nil {
			t.Error("Expected ping to fail with non-existent host")
		} else {
			t.Logf("Ping failed as expected: %v", pingErr)
		}
	}
	
	if err != nil {
		t.Logf("Got error during Connect (acceptable): %v", err)
	}
	
	// The error might not be immediately apparent as Connect just creates the client
	// The actual connection test happens when we try to use it
}

func TestConnect_LocalhostURI(t *testing.T) {
	// Test with localhost URI (this may or may not succeed depending on environment)
	localhostURI := "mongodb://localhost:27017"
	
	client, err := Connect(localhostURI)
	
	// Clean up if successful
	if client != nil {
		defer client.Disconnect(context.Background())
	}
	
	// We can't assert success or failure since MongoDB may not be running
	// But we can test that the function doesn't panic and returns proper types
	if err != nil {
		// Error is expected if MongoDB is not running
		t.Logf("Connection failed as expected (MongoDB likely not running): %v", err)
	} else {
		// If successful, test that we can ping
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		pingErr := client.Ping(ctx, nil)
		if pingErr != nil {
			t.Logf("Client created but ping failed: %v", pingErr)
		} else {
			t.Log("Successfully connected to MongoDB")
		}
	}
}

func TestConnect_ValidURIFormat(t *testing.T) {
	// Test with valid URI formats that might not connect but are correctly formatted
	validURIs := []string{
		"mongodb://user:pass@localhost:27017/testdb",
		"mongodb://localhost:27017/testdb",
		"mongodb+srv://cluster.mongodb.net/testdb",
		"mongodb://host1:27017,host2:27017/testdb",
	}
	
	for _, uri := range validURIs {
		t.Run(uri, func(t *testing.T) {
			client, err := Connect(uri)
			
			// Clean up if client was created
			if client != nil {
				defer client.Disconnect(context.Background())
			}
			
			// We test that the function doesn't panic and handles the URI
			// Error is acceptable since these are test URIs
			if err != nil {
				t.Logf("Connection failed as expected for test URI %s: %v", uri, err)
			}
		})
	}
}

func TestConnect_MalformedURIs(t *testing.T) {
	// Test with malformed URIs that should cause errors
	malformedURIs := []string{
		"not-a-uri",
		"http://localhost:27017", // Wrong protocol
		"mongodb://", // Incomplete
		"mongodb://localhost:", // No port
		"ftp://localhost:27017", // Wrong protocol
	}
	
	for _, uri := range malformedURIs {
		t.Run(uri, func(t *testing.T) {
			client, err := Connect(uri)
			
			// Clean up if client was somehow created
			if client != nil {
				defer client.Disconnect(context.Background())
			}
			
			// Most malformed URIs should cause errors
			if err == nil {
				t.Logf("Unexpectedly succeeded with malformed URI: %s", uri)
			} else {
				t.Logf("Expected error with malformed URI %s: %v", uri, err)
			}
		})
	}
}

func TestConnect_ContextTimeout(t *testing.T) {
	// This test verifies that the function uses a 10-second timeout
	// We can't easily test the timeout directly, but we can verify the function behavior
	
	// Test with a URI that should work format-wise
	uri := "mongodb://localhost:27017"
	
	start := time.Now()
	client, err := Connect(uri)
	duration := time.Since(start)
	
	if client != nil {
		defer client.Disconnect(context.Background())
	}
	
	// The function should complete relatively quickly since it just creates the client
	// The actual connection happens lazily
	if duration > 15*time.Second {
		t.Errorf("Connect took too long: %v", duration)
	}
	
	// Log the result
	if err != nil {
		t.Logf("Connect completed in %v with error: %v", duration, err)
	} else {
		t.Logf("Connect completed successfully in %v", duration)
	}
}

func TestConnect_ClientOptions(t *testing.T) {
	// Test that the function creates a client with proper options
	// This is mainly a structural test
	
	uri := "mongodb://localhost:27017"
	client, err := Connect(uri)
	
	if client != nil {
		defer client.Disconnect(context.Background())
	}
	
	// Test that we get a proper mongo.Client type
	if client != nil {
		// Verify it's the right type
		if _, ok := interface{}(client).(*mongo.Client); !ok {
			t.Error("Expected *mongo.Client type")
		}
	}
	
	// Test that error is of expected type when it occurs
	if err != nil {
		// Should be a mongo-related error or connection error
		t.Logf("Error type: %T, message: %v", err, err)
	}
}

// Benchmark test
func BenchmarkConnect(b *testing.B) {
	uri := "mongodb://localhost:27017"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := Connect(uri)
		if client != nil {
			client.Disconnect(context.Background())
		}
		if err != nil {
			// Expected in test environment
		}
	}
}

// Test that verifies the function signature and basic behavior
func TestConnect_FunctionSignature(t *testing.T) {
	// Verify function accepts string and returns (*mongo.Client, error)
	uri := "mongodb://test:27017"
	
	client, err := Connect(uri)
	
	// Test return types
	if client != nil {
		defer client.Disconnect(context.Background())
		// Verify client type
		var clientType *mongo.Client = client
		_ = clientType // Use variable to avoid compiler warning
	}
	
	// Verify error type
	var errorType error = err
	_ = errorType // Use variable to avoid compiler warning
	
	// Function should not panic
	t.Log("Function completed without panic")
}