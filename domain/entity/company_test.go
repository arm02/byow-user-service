package entity

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCompanyStruct(t *testing.T) {
	now := time.Now()
	objID := primitive.NewObjectID()
	
	company := Company{
		ID:             objID,
		UserID:         "user123",
		CompanyName:    "Test Company",
		CompanyEmail:   "info@testcompany.com",
		CompanyPhone:   "+1234567890",
		CompanyAddress: "123 Test Street, Test City",
		CompanyLogo:    "https://example.com/logo.png",
		Verified:       true,
		CreatedAt:      now,
	}

	// Test that all fields are properly set
	if company.ID != objID {
		t.Errorf("Expected ID %v, got %v", objID, company.ID)
	}

	if company.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %v", company.UserID)
	}

	if company.CompanyName != "Test Company" {
		t.Errorf("Expected CompanyName 'Test Company', got %v", company.CompanyName)
	}

	if company.CompanyEmail != "info@testcompany.com" {
		t.Errorf("Expected CompanyEmail 'info@testcompany.com', got %v", company.CompanyEmail)
	}

	if company.CompanyPhone != "+1234567890" {
		t.Errorf("Expected CompanyPhone '+1234567890', got %v", company.CompanyPhone)
	}

	if company.CompanyAddress != "123 Test Street, Test City" {
		t.Errorf("Expected CompanyAddress '123 Test Street, Test City', got %v", company.CompanyAddress)
	}

	if company.CompanyLogo != "https://example.com/logo.png" {
		t.Errorf("Expected CompanyLogo 'https://example.com/logo.png', got %v", company.CompanyLogo)
	}

	if !company.Verified {
		t.Error("Expected Verified to be true")
	}

	if company.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, company.CreatedAt)
	}
}

func TestCompanyStructZeroValues(t *testing.T) {
	var company Company

	// Test zero values
	if !company.ID.IsZero() {
		t.Errorf("Expected zero ObjectID, got %v", company.ID)
	}

	if company.UserID != "" {
		t.Errorf("Expected empty UserID, got %v", company.UserID)
	}

	if company.CompanyName != "" {
		t.Errorf("Expected empty CompanyName, got %v", company.CompanyName)
	}

	if company.CompanyEmail != "" {
		t.Errorf("Expected empty CompanyEmail, got %v", company.CompanyEmail)
	}

	if company.CompanyPhone != "" {
		t.Errorf("Expected empty CompanyPhone, got %v", company.CompanyPhone)
	}

	if company.CompanyAddress != "" {
		t.Errorf("Expected empty CompanyAddress, got %v", company.CompanyAddress)
	}

	if company.CompanyLogo != "" {
		t.Errorf("Expected empty CompanyLogo, got %v", company.CompanyLogo)
	}

	if company.Verified {
		t.Error("Expected Verified to be false")
	}

	if !company.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", company.CreatedAt)
	}
}

func TestCompanyStructPartialValues(t *testing.T) {
	now := time.Now()
	company := Company{
		CompanyName: "Partial Company",
		Verified:    false,
		CreatedAt:   now,
	}

	// Test that specified fields are set
	if company.CompanyName != "Partial Company" {
		t.Errorf("Expected CompanyName 'Partial Company', got %v", company.CompanyName)
	}

	if company.Verified {
		t.Error("Expected Verified to be false")
	}

	if company.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, company.CreatedAt)
	}

	// Test that unspecified fields have zero values
	if !company.ID.IsZero() {
		t.Errorf("Expected zero ObjectID, got %v", company.ID)
	}

	if company.UserID != "" {
		t.Errorf("Expected empty UserID, got %v", company.UserID)
	}

	if company.CompanyEmail != "" {
		t.Errorf("Expected empty CompanyEmail, got %v", company.CompanyEmail)
	}
}

func TestCompanyStructWithNewObjectID(t *testing.T) {
	// Test creating multiple companies with new ObjectIDs
	company1 := Company{ID: primitive.NewObjectID()}
	company2 := Company{ID: primitive.NewObjectID()}

	// IDs should be different
	if company1.ID == company2.ID {
		t.Error("Expected different ObjectIDs for different companies")
	}

	// IDs should not be zero
	if company1.ID.IsZero() || company2.ID.IsZero() {
		t.Error("Expected non-zero ObjectIDs")
	}
}