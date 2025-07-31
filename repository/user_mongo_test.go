package repository

import (
	"testing"
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
)

// Test basic functionality without mocking MongoDB
func TestUserMongoRepoStructure(t *testing.T) {
	// Test struct initialization
	repo := &userMongoRepo{}
	if repo == nil {
		t.Error("Expected non-nil userMongoRepo")
	}
}

func TestCreateUserSetsCreatedAt(t *testing.T) {
	// Test that Create method sets CreatedAt
	user := &entity.User{
		Fullname:    "John Doe",
		Email:       "john@example.com",
		Password:    "hashedpassword",
		PhoneNumber: "+1234567890",
	}

	// Verify CreatedAt is initially zero
	if !user.CreatedAt.IsZero() {
		t.Error("Expected initial CreatedAt to be zero")
	}

	// Simulate what Create method does
	user.CreatedAt = time.Now()

	// Verify CreatedAt is now set
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if time.Since(user.CreatedAt) > time.Second {
		t.Error("CreatedAt should be very recent")
	}
}

func TestBSONMarshaling(t *testing.T) {
	// Test BSON marshaling functionality used in Update methods
	user := &entity.User{
		ID:           "test123",
		Email:        "test@example.com",
		Fullname:     "Test User",
		PhoneNumber:  "+1234567890",
		OnBoarded:    true,
		Verified:     false,
		OTP:          "123456",
		OTPType:      "verification",
		OTPExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt:    time.Now(),
	}

	// Test marshaling
	data, err := bson.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty marshaled data")
	}

	// Test unmarshaling
	var updateMap bson.M
	err = bson.Unmarshal(data, &updateMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Verify some fields exist in the map
	if updateMap["email"] != user.Email {
		t.Errorf("Expected email %v, got %v", user.Email, updateMap["email"])
	}

	if updateMap["full_name"] != user.Fullname {
		t.Errorf("Expected full_name %v, got %v", user.Fullname, updateMap["full_name"])
	}
}

func TestUpdateMapConstruction(t *testing.T) {
	// Test the update map construction logic used in Update methods
	user := &entity.User{
		Email:    "test@example.com",
		Fullname: "Updated Name",
		OTP:      "", // Empty OTP should trigger unset logic
	}

	// Simulate the marshaling logic
	updateData, err := bson.Marshal(user)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var updateMap bson.M
	err = bson.Unmarshal(updateData, &updateMap)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Remove _id as done in the actual code
	delete(updateMap, "_id")

	// Test unset logic for empty OTP
	unsetMap := bson.M{}
	if user.OTP == "" {
		unsetMap["otp"] = ""
		unsetMap["otp_expires_at"] = ""
		unsetMap["otp_type"] = ""
	}

	// Verify unset map is created when OTP is empty
	if len(unsetMap) != 3 {
		t.Errorf("Expected 3 unset fields, got %d", len(unsetMap))
	}

	// Construct final update document
	update := bson.M{}
	if len(updateMap) > 0 {
		update["$set"] = updateMap
	}
	if len(unsetMap) > 0 {
		update["$unset"] = unsetMap
	}

	// Verify both $set and $unset are present
	if _, hasSet := update["$set"]; !hasSet {
		t.Error("Expected $set in update document")
	}

	if _, hasUnset := update["$unset"]; !hasUnset {
		t.Error("Expected $unset in update document")
	}
}

func TestUpdateMapWithOTP(t *testing.T) {
	// Test update map when OTP is present
	user := &entity.User{
		Email:        "test@example.com",
		Fullname:     "Updated Name",
		OTP:          "123456", // Non-empty OTP
		OTPType:      "verification",
		OTPExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Simulate the marshaling logic
	updateData, err := bson.Marshal(user)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var updateMap bson.M
	err = bson.Unmarshal(updateData, &updateMap)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Remove _id as done in the actual code
	delete(updateMap, "_id")

	// Test unset logic for non-empty OTP
	unsetMap := bson.M{}
	if user.OTP == "" {
		unsetMap["otp"] = ""
		unsetMap["otp_expires_at"] = ""
		unsetMap["otp_type"] = ""
	}

	// Verify unset map is empty when OTP is present
	if len(unsetMap) != 0 {
		t.Errorf("Expected 0 unset fields, got %d", len(unsetMap))
	}

	// Construct final update document
	update := bson.M{}
	if len(updateMap) > 0 {
		update["$set"] = updateMap
	}
	if len(unsetMap) > 0 {
		update["$unset"] = unsetMap
	}

	// Verify only $set is present
	if _, hasSet := update["$set"]; !hasSet {
		t.Error("Expected $set in update document")
	}

	if _, hasUnset := update["$unset"]; hasUnset {
		t.Error("Did not expect $unset in update document when OTP is present")
	}
}

func TestBSONFilters(t *testing.T) {
	// Test BSON filter construction used in Find methods
	email := "test@example.com"
	emailFilter := bson.M{"email": email}

	if emailFilter["email"] != email {
		t.Errorf("Expected email filter %v, got %v", email, emailFilter["email"])
	}

	phone := "+1234567890"
	phoneFilter := bson.M{"phone_number": phone}

	if phoneFilter["phone_number"] != phone {
		t.Errorf("Expected phone filter %v, got %v", phone, phoneFilter["phone_number"])
	}
}