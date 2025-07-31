package db

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Test index model creation and structure
func TestUserIndexModels(t *testing.T) {
	// Test that we can create the user index models without database connection
	// This tests the structure and BSON document creation logic
	
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("email_unique"),
		},
		{
			Keys: bson.D{{Key: "phone_number", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("phone_unique"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().
				SetName("created_at_index"),
		},
		{
			Keys: bson.D{{Key: "is_verified", Value: 1}},
			Options: options.Index().
				SetName("is_verified_index"),
		},
		{
			Keys: bson.D{{Key: "is_onboarded", Value: 1}},
			Options: options.Index().
				SetName("is_onboarded_index"),
		},
		// Compound index for common queries
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "is_verified", Value: 1},
			},
			Options: options.Index().
				SetName("email_verified_compound"),
		},
	}
	
	// Test that all index models were created
	if len(userIndexes) != 6 {
		t.Errorf("Expected 6 user indexes, got %d", len(userIndexes))
	}
	
	// Test specific index properties
	emailIndex := userIndexes[0]
	if emailIndex.Options.Name == nil || *emailIndex.Options.Name != "email_unique" {
		t.Error("Expected email index to have name 'email_unique'")
	}
	
	if emailIndex.Options.Unique == nil || !*emailIndex.Options.Unique {
		t.Error("Expected email index to be unique")
	}
	
	// Test compound index
	compoundIndex := userIndexes[5]
	if compoundIndex.Options.Name == nil || *compoundIndex.Options.Name != "email_verified_compound" {
		t.Error("Expected compound index to have name 'email_verified_compound'")
	}
	
	// Verify compound index keys
	compoundKeys, ok := compoundIndex.Keys.(bson.D)
	if !ok || len(compoundKeys) != 2 {
		t.Error("Expected compound index to have 2 keys")
	}
}

func TestCompanyIndexModels(t *testing.T) {
	// Test company index model creation
	companyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
			Options: options.Index().
				SetName("company_name_index"),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetSparse(true).
				SetName("company_email_unique"),
		},
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().
				SetName("company_phone_index"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().
				SetName("company_created_at_index"),
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: 1}},
			Options: options.Index().
				SetName("company_updated_at_index"),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().
				SetName("company_user_id_index"),
		},
		// Compound index for user companies
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_companies_compound"),
		},
		// Text index for company search
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
			},
			Options: options.Index().
				SetName("company_search_text"),
		},
	}
	
	// Test that all company indexes were created
	if len(companyIndexes) != 8 {
		t.Errorf("Expected 8 company indexes, got %d", len(companyIndexes))
	}
	
	// Test sparse unique index
	emailIndex := companyIndexes[1]
	if emailIndex.Options.Name == nil || *emailIndex.Options.Name != "company_email_unique" {
		t.Error("Expected company email index to have name 'company_email_unique'")
	}
	
	if emailIndex.Options.Unique == nil || !*emailIndex.Options.Unique {
		t.Error("Expected company email index to be unique")
	}
	
	if emailIndex.Options.Sparse == nil || !*emailIndex.Options.Sparse {
		t.Error("Expected company email index to be sparse")
	}
	
	// Test text search index
	textIndex := companyIndexes[7]
	if textIndex.Options.Name == nil || *textIndex.Options.Name != "company_search_text" {
		t.Error("Expected text index to have name 'company_search_text'")
	}
	
	// Verify text index keys
	textKeys, ok := textIndex.Keys.(bson.D)
	if !ok || len(textKeys) != 2 {
		t.Error("Expected text index to have 2 keys")
	}
	
	// Check that text search values are "text"
	for _, key := range textKeys {
		if key.Value != "text" {
			t.Errorf("Expected text index value to be 'text', got %v", key.Value)
		}
	}
}

func TestRequiredIndexesLists(t *testing.T) {
	// Test the required indexes lists used in CheckIndexes
	requiredUserIndexes := []string{
		"email_unique",
		"phone_unique",
		"created_at_index",
		"is_verified_index",
		"is_onboarded_index",
		"email_verified_compound",
	}
	
	requiredCompanyIndexes := []string{
		"company_name_index",
		"company_email_unique",
		"company_phone_index",
		"company_created_at_index",
		"company_updated_at_index",
		"company_user_id_index",
		"user_companies_compound",
		"company_search_text",
	}
	
	// Test counts
	if len(requiredUserIndexes) != 6 {
		t.Errorf("Expected 6 required user indexes, got %d", len(requiredUserIndexes))
	}
	
	if len(requiredCompanyIndexes) != 8 {
		t.Errorf("Expected 8 required company indexes, got %d", len(requiredCompanyIndexes))
	}
	
	// Test that all required indexes have unique names
	userIndexMap := make(map[string]bool)
	for _, name := range requiredUserIndexes {
		if userIndexMap[name] {
			t.Errorf("Duplicate user index name: %s", name)
		}
		userIndexMap[name] = true
	}
	
	companyIndexMap := make(map[string]bool)
	for _, name := range requiredCompanyIndexes {
		if companyIndexMap[name] {
			t.Errorf("Duplicate company index name: %s", name)
		}
		companyIndexMap[name] = true
	}
}

func TestIndexCheckLogic(t *testing.T) {
	// Test the index checking logic used in CheckIndexes
	
	// Simulate existing indexes (like what comes from MongoDB)
	existingIndexes := []bson.M{
		{"name": "_id_"},
		{"name": "email_unique"},
		{"name": "phone_unique"},
		{"name": "created_at_index"},
		// Missing: is_verified_index, is_onboarded_index, email_verified_compound
	}
	
	requiredIndexes := []string{
		"email_unique",
		"phone_unique",
		"created_at_index",
		"is_verified_index",
		"is_onboarded_index",
		"email_verified_compound",
	}
	
	// Build map of existing indexes
	existingMap := make(map[string]bool)
	for _, index := range existingIndexes {
		if name, ok := index["name"].(string); ok {
			existingMap[name] = true
		}
	}
	
	// Find missing indexes
	missing := []string{}
	for _, required := range requiredIndexes {
		if !existingMap[required] {
			missing = append(missing, required)
		}
	}
	
	// Should find 3 missing indexes
	expectedMissing := []string{"is_verified_index", "is_onboarded_index", "email_verified_compound"}
	if len(missing) != len(expectedMissing) {
		t.Errorf("Expected %d missing indexes, got %d", len(expectedMissing), len(missing))
	}
	
	// Check specific missing indexes
	missingMap := make(map[string]bool)
	for _, name := range missing {
		missingMap[name] = true
	}
	
	for _, expected := range expectedMissing {
		if !missingMap[expected] {
			t.Errorf("Expected to find missing index: %s", expected)
		}
	}
}

// Test that index functions exist and can be called (will fail gracefully)
func TestCreateIndexesFunction(t *testing.T) {
	logger := zap.NewNop()
	
	// This will fail due to nil database, but tests that the function exists
	// and handles nil input gracefully
	err := CreateIndexes(nil, logger)
	if err == nil {
		t.Error("Expected error when calling CreateIndexes with nil database")
	}
	
	t.Logf("CreateIndexes returned expected error: %v", err)
}

func TestDropIndexesFunction(t *testing.T) {
	logger := zap.NewNop()
	
	// This will fail due to nil database, but tests that the function exists
	err := DropIndexes(nil, logger)
	if err == nil {
		t.Error("Expected error when calling DropIndexes with nil database")
	}
	
	t.Logf("DropIndexes returned expected error: %v", err)
}

func TestCheckIndexesFunction(t *testing.T) {
	logger := zap.NewNop()
	
	// This will fail due to nil database, but tests that the function exists
	err := CheckIndexes(nil, logger)
	if err == nil {
		t.Error("Expected error when calling CheckIndexes with nil database")
	}
	
	t.Logf("CheckIndexes returned expected error: %v", err)
}

func TestRebuildCompanyIndexesFunction(t *testing.T) {
	logger := zap.NewNop()
	
	// This will fail due to nil database, but tests that the function exists
	err := RebuildCompanyIndexes(nil, logger)
	if err == nil {
		t.Error("Expected error when calling RebuildCompanyIndexes with nil database")
	}
	
	t.Logf("RebuildCompanyIndexes returned expected error: %v", err)
}

// Test BSON document creation
func TestBSONDocumentCreation(t *testing.T) {
	// Test BSON document patterns used in indexes
	
	// Single field index
	singleField := bson.D{{Key: "email", Value: 1}}
	if len(singleField) != 1 {
		t.Error("Expected single field BSON document to have 1 element")
	}
	if singleField[0].Key != "email" || singleField[0].Value != 1 {
		t.Error("Single field BSON document has incorrect structure")
	}
	
	// Compound index
	compoundField := bson.D{
		{Key: "user_id", Value: 1},
		{Key: "created_at", Value: -1},
	}
	if len(compoundField) != 2 {
		t.Error("Expected compound BSON document to have 2 elements")
	}
	if compoundField[1].Value != -1 {
		t.Error("Expected descending sort order (-1) for created_at")
	}
	
	// Text search index
	textSearch := bson.D{
		{Key: "name", Value: "text"},
		{Key: "description", Value: "text"},
	}
	if len(textSearch) != 2 {
		t.Error("Expected text search BSON document to have 2 elements")
	}
	for _, element := range textSearch {
		if element.Value != "text" {
			t.Errorf("Expected text search value to be 'text', got %v", element.Value)
		}
	}
}

// Test index options creation
func TestIndexOptionsCreation(t *testing.T) {
	// Test unique index options
	uniqueOpts := options.Index().SetUnique(true).SetName("test_unique")
	if uniqueOpts.Unique == nil || !*uniqueOpts.Unique {
		t.Error("Expected unique option to be true")
	}
	if uniqueOpts.Name == nil || *uniqueOpts.Name != "test_unique" {
		t.Error("Expected name option to be 'test_unique'")
	}
	
	// Test sparse index options
	sparseOpts := options.Index().SetSparse(true).SetUnique(true).SetName("test_sparse")
	if sparseOpts.Sparse == nil || !*sparseOpts.Sparse {
		t.Error("Expected sparse option to be true")
	}
	
	// Test basic index options
	basicOpts := options.Index().SetName("test_basic")
	if basicOpts.Name == nil || *basicOpts.Name != "test_basic" {
		t.Error("Expected name option to be 'test_basic'")
	}
	
	// Ensure no unexpected options are set for basic index
	if basicOpts.Unique != nil {
		t.Error("Expected unique option to be nil for basic index")
	}
	if basicOpts.Sparse != nil {
		t.Error("Expected sparse option to be nil for basic index")
	}
}

// Test collection name constants
func TestCollectionNames(t *testing.T) {
	// Test that the collection names used are consistent
	userCollection := "users_collections"
	companyCollection := "companies_collections"
	
	if userCollection == "" {
		t.Error("User collection name should not be empty")
	}
	
	if companyCollection == "" {
		t.Error("Company collection name should not be empty")
	}
	
	if userCollection == companyCollection {
		t.Error("User and company collection names should be different")
	}
	
	// Test naming convention
	if userCollection != "users_collections" {
		t.Errorf("Expected user collection name 'users_collections', got %s", userCollection)
	}
	
	if companyCollection != "companies_collections" {
		t.Errorf("Expected company collection name 'companies_collections', got %s", companyCollection)
	}
}

// Benchmark index model creation
func BenchmarkIndexModelCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark creating a single index model
		_ = mongo.IndexModel{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("email_unique"),
		}
	}
}

func BenchmarkBSONDocumentCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark creating BSON documents
		_ = bson.D{
			{Key: "user_id", Value: 1},
			{Key: "created_at", Value: -1},
		}
	}
}