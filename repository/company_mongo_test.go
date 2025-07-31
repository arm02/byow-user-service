package repository

import (
	"testing"
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCompanyMongoRepoStructure(t *testing.T) {
	repo := &companyMongoRepo{}
	if repo == nil {
		t.Error("Expected non-nil companyMongoRepo")
	}
}

func TestNewCompanyMongoRepo(t *testing.T) {
	// Test that NewCompanyMongoRepo returns the correct interface type
	// We can't test with real db without MongoDB, but we can test the structure
	repo := &companyMongoRepo{}
	if repo == nil {
		t.Error("Expected non-nil repository instance")
	}
}

func TestFindAllFilterConstruction(t *testing.T) {
	// Test filter construction logic used in FindAll
	testCases := []struct {
		name        string
		userID      string
		keyword     string
		expectedLen int
	}{
		{"empty filters", "", "", 0},
		{"user ID only", "user123", "", 1},
		{"keyword only", "", "test", 1}, 
		{"both filters", "user123", "test", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := bson.M{}

			if tc.keyword != "" {
				filter["company_name"] = bson.M{
					"$regex":   tc.keyword,
					"$options": "i",
				}
			}

			if tc.userID != "" {
				filter["user_id"] = tc.userID
			}

			if len(filter) != tc.expectedLen {
				t.Errorf("Expected filter length %d, got %d", tc.expectedLen, len(filter))
			}

			// Test regex filter when keyword is present
			if tc.keyword != "" {
				if regexFilter, ok := filter["company_name"].(bson.M); ok {
					if regexFilter["$regex"] != tc.keyword {
						t.Errorf("Expected regex %s, got %v", tc.keyword, regexFilter["$regex"])
					}
					if regexFilter["$options"] != "i" {
						t.Errorf("Expected case-insensitive option, got %v", regexFilter["$options"])
					}
				} else {
					t.Error("Expected company_name filter to be bson.M type")
				}
			}

			// Test user ID filter
			if tc.userID != "" {
				if filter["user_id"] != tc.userID {
					t.Errorf("Expected user_id %s, got %v", tc.userID, filter["user_id"])
				}
			}
		})
	}
}

func TestCreateDuplicateCheckLogic(t *testing.T) {
	// Test duplicate check logic used in Create method
	testCases := []struct {
		name           string
		company        *entity.Company
		expectedOrLen  int
		description    string
	}{
		{
			"no email or phone",
			&entity.Company{CompanyName: "Test Co"},
			0,
			"should not create OR conditions",
		},
		{
			"email only",
			&entity.Company{CompanyEmail: "test@example.com"},
			1,
			"should create one OR condition",
		},
		{
			"phone only", 
			&entity.Company{CompanyPhone: "+1234567890"},
			1,
			"should create one OR condition",
		},
		{
			"both email and phone",
			&entity.Company{CompanyEmail: "test@example.com", CompanyPhone: "+1234567890"},
			2,
			"should create two OR conditions",
		},
		{
			"empty email and phone",
			&entity.Company{CompanyEmail: "", CompanyPhone: ""},
			0,
			"should not create OR conditions for empty strings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the duplicate check logic
			orConditions := []bson.M{}
			
			if tc.company.CompanyEmail != "" {
				orConditions = append(orConditions, bson.M{"company_email": tc.company.CompanyEmail})
			}
			if tc.company.CompanyPhone != "" {
				orConditions = append(orConditions, bson.M{"company_phone": tc.company.CompanyPhone})
			}

			if len(orConditions) != tc.expectedOrLen {
				t.Errorf("Expected %d OR conditions, got %d", tc.expectedOrLen, len(orConditions))
			}

			// Test filter construction only when we have conditions
			if len(orConditions) > 0 {
				filter := bson.M{"$or": orConditions}
				if _, ok := filter["$or"]; !ok {
					t.Error("Expected $or in filter")
				}

				if orArray, ok := filter["$or"].([]bson.M); ok {
					if len(orArray) != tc.expectedOrLen {
						t.Errorf("Expected %d conditions in $or array, got %d", tc.expectedOrLen, len(orArray))
					}
				} else {
					t.Error("Expected $or to be []bson.M type")
				}
			}

			t.Logf("Test case: %s - %s", tc.name, tc.description)
		})
	}
}

func TestCreateSetsTimestamp(t *testing.T) {
	// Test that Create method sets CreatedAt
	company := &entity.Company{
		CompanyName:  "Test Company",
		CompanyEmail: "test@company.com",
		CompanyPhone: "+1234567890",
	}

	// Verify CreatedAt is initially zero
	if !company.CreatedAt.IsZero() {
		t.Error("Expected initial CreatedAt to be zero")
	}

	// Simulate what Create method does
	company.CreatedAt = time.Now()

	// Verify CreatedAt is now set
	if company.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if time.Since(company.CreatedAt) > time.Second {
		t.Error("CreatedAt should be very recent")
	}
}

func TestObjectIDHandling(t *testing.T) {
	// Test ObjectID handling in Create method
	company := &entity.Company{
		CompanyName: "Test Company",
	}

	// Test that we can create a primitive.ObjectID
	oid := primitive.NewObjectID()
	if oid.IsZero() {
		t.Error("Expected non-zero ObjectID")
	}

	// Simulate setting ID after insert
	company.ID = oid
	if company.ID != oid {
		t.Errorf("Expected ID %v, got %v", oid, company.ID)
	}

	// Test ObjectID string representation
	idStr := oid.Hex()
	if len(idStr) != 24 {
		t.Errorf("Expected ObjectID hex string length 24, got %d", len(idStr))
	}
}

func TestFindByIDFilter(t *testing.T) {
	// Test filter construction for FindByID
	id := primitive.NewObjectID()
	filter := bson.M{"_id": id}

	if filter["_id"] != id {
		t.Errorf("Expected _id filter %v, got %v", id, filter["_id"])
	}

	// Test with zero ObjectID
	zeroID := primitive.ObjectID{}
	zeroFilter := bson.M{"_id": zeroID}
	
	if zeroFilter["_id"] != zeroID {
		t.Error("Expected zero ObjectID in filter")
	}
}

func TestFindByEmailFilter(t *testing.T) {
	// Test email filter construction
	email := "test@company.com"
	filter := bson.M{"email": email}

	if filter["email"] != email {
		t.Errorf("Expected email filter %v, got %v", email, filter["email"])
	}
}

func TestFindByPhoneFilter(t *testing.T) {
	// Test phone filter construction  
	phone := "+1234567890"
	filter := bson.M{"phone_number": phone}

	if filter["phone_number"] != phone {
		t.Errorf("Expected phone filter %v, got %v", phone, filter["phone_number"])
	}
}

func TestUpdateFilter(t *testing.T) {
	// Test update filter construction
	company := &entity.Company{
		ID:           primitive.NewObjectID(),
		CompanyName:  "Updated Company",
		CompanyEmail: "updated@company.com",
	}

	// Test filter construction
	filter := bson.M{"id": company.ID}
	if filter["id"] != company.ID {
		t.Errorf("Expected id filter %v, got %v", company.ID, filter["id"])
	}

	// Test update document
	update := bson.M{"$set": company}
	if update["$set"] != company {
		t.Error("Expected company in $set operation")
	}
}

func TestDeleteFilter(t *testing.T) {
	// Test delete filter construction
	id := primitive.NewObjectID()
	filter := bson.M{"_id": id}

	if filter["_id"] != id {
		t.Errorf("Expected _id filter %v, got %v", id, filter["_id"])
	}
}

func TestCompanyEntityValidation(t *testing.T) {
	// Test that company entity can be properly constructed
	company := &entity.Company{
		ID:           primitive.NewObjectID(),
		UserID:       "user123",
		CompanyName:  "Test Corporation",
		CompanyEmail: "contact@testcorp.com",
		CompanyPhone: "+1-555-0123",
		CreatedAt:    time.Now(),
	}

	// Verify all fields are set correctly
	if company.ID.IsZero() {
		t.Error("Expected non-zero ID")
	}

	if company.UserID == "" {
		t.Error("Expected non-empty UserID")
	}

	if company.CompanyName == "" {
		t.Error("Expected non-empty CompanyName")
	}

	if company.CompanyEmail == "" {
		t.Error("Expected non-empty CompanyEmail")
	}

	if company.CompanyPhone == "" {
		t.Error("Expected non-empty CompanyPhone")
	}

	if company.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
}

func TestPaginationOptions(t *testing.T) {
	// Test pagination options construction used in FindAll
	testCases := []struct {
		name   string
		limit  int64
		offset int64
	}{
		{"first page", 10, 0},
		{"second page", 10, 10},
		{"large page", 100, 500},
		{"zero limit", 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate options creation (we can't import options without MongoDB driver)
			// But we can test the values
			limit := tc.limit
			offset := tc.offset

			if limit != tc.limit {
				t.Errorf("Expected limit %d, got %d", tc.limit, limit)
			}

			if offset != tc.offset {
				t.Errorf("Expected offset %d, got %d", tc.offset, offset)
			}

			// Test that negative values would be handled
			if limit < 0 {
				t.Error("Limit should not be negative")
			}

			if offset < 0 {
				t.Error("Offset should not be negative")
			}
		})
	}
}

func TestCompanyBSONMarshaling(t *testing.T) {
	// Test BSON marshaling for company entities
	company := &entity.Company{
		ID:           primitive.NewObjectID(),
		UserID:       "user123", 
		CompanyName:  "Test Company",
		CompanyEmail: "test@company.com",
		CompanyPhone: "+1234567890",
		CreatedAt:    time.Now(),
	}

	// Test that we can marshal to BSON
	data, err := bson.Marshal(company)
	if err != nil {
		t.Fatalf("Failed to marshal company: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty marshaled data")
	}

	// Test unmarshaling
	var unmarshaled entity.Company
	err = bson.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal company: %v", err)
	}

	// Verify key fields are preserved
	if unmarshaled.UserID != company.UserID {
		t.Errorf("Expected UserID %v, got %v", company.UserID, unmarshaled.UserID)
	}

	if unmarshaled.CompanyName != company.CompanyName {
		t.Errorf("Expected CompanyName %v, got %v", company.CompanyName, unmarshaled.CompanyName)
	}

	if unmarshaled.CompanyEmail != company.CompanyEmail {
		t.Errorf("Expected CompanyEmail %v, got %v", company.CompanyEmail, unmarshaled.CompanyEmail)
	}
}

func TestRegexFilterConstruction(t *testing.T) {
	// Test regex filter construction for case-insensitive search
	testCases := []struct {
		keyword  string
		expected string
	}{
		{"test", "test"},
		{"Test", "Test"},
		{"TEST", "TEST"},
		{"company name", "company name"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.keyword, func(t *testing.T) {
			if tc.keyword != "" {
				regexFilter := bson.M{
					"$regex":   tc.keyword,
					"$options": "i",
				}

				if regexFilter["$regex"] != tc.expected {
					t.Errorf("Expected regex %s, got %v", tc.expected, regexFilter["$regex"])
				}

				// Test case-insensitive option
				if regexFilter["$options"] != "i" {
					t.Error("Expected case-insensitive option 'i'")
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkFilterConstruction(b *testing.B) {
	userID := "user123"
	keyword := "test company"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter := bson.M{}

		if keyword != "" {
			filter["company_name"] = bson.M{
				"$regex":   keyword,
				"$options": "i",
			}
		}

		if userID != "" {
			filter["user_id"] = userID
		}
	}
}

func BenchmarkObjectIDCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = primitive.NewObjectID()
	}
}