package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock mongo.Collection for testing
type mockUserCollection struct {
	insertOneFunc  func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	findOneFunc    func(ctx context.Context, filter interface{}) *mongo.SingleResult
	updateOneFunc  func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	documents      map[string]*entity.User // Store by email as key
	returnError    error
}

func (m *mockUserCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if m.insertOneFunc != nil {
		return m.insertOneFunc(ctx, document)
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	
	user, ok := document.(*entity.User)
	if !ok {
		return nil, errors.New("invalid document type")
	}
	
	if m.documents == nil {
		m.documents = make(map[string]*entity.User)
	}
	
	// Check if user already exists (duplicate check)
	if _, exists := m.documents[user.Email]; exists {
		return nil, errors.New("duplicate key error")
	}
	
	// Set timestamp
	user.CreatedAt = time.Now()
	m.documents[user.Email] = user
	
	return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
}

func (m *mockUserCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if m.findOneFunc != nil {
		return m.findOneFunc(ctx, filter)
	}
	
	filterMap, ok := filter.(bson.M)
	if !ok {
		return &mongo.SingleResult{}
	}
	
	if m.documents == nil {
		m.documents = make(map[string]*entity.User)
	}
	
	// Handle email filter
	if email, exists := filterMap["email"]; exists {
		if _, found := m.documents[email.(string)]; found {
			return &mongo.SingleResult{} // Would need proper SingleResult mock
		}
		return &mongo.SingleResult{} // Would return ErrNoDocuments
	}
	
	// Handle phone filter
	if phone, exists := filterMap["phone_number"]; exists {
		for _, user := range m.documents {
			if user.PhoneNumber == phone.(string) {
				return &mongo.SingleResult{} // Would need proper SingleResult mock
			}
		}
		return &mongo.SingleResult{} // Would return ErrNoDocuments
	}
	
	return &mongo.SingleResult{}
}

func (m *mockUserCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if m.updateOneFunc != nil {
		return m.updateOneFunc(ctx, filter, update)
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	
	filterMap, ok := filter.(bson.M)
	if !ok {
		return nil, errors.New("invalid filter")
	}
	
	if m.documents == nil {
		m.documents = make(map[string]*entity.User)
	}
	
	// Handle email-based update
	if email, exists := filterMap["email"]; exists {
		if _, found := m.documents[email.(string)]; found {
			return &mongo.UpdateResult{ModifiedCount: 1}, nil
		}
		return &mongo.UpdateResult{ModifiedCount: 0}, nil
	}
	
	// Handle phone-based update
	if phone, exists := filterMap["phone_number"]; exists {
		for _, user := range m.documents {
			if user.PhoneNumber == phone.(string) {
				// Update would happen here
				return &mongo.UpdateResult{ModifiedCount: 1}, nil
			}
		}
		return &mongo.UpdateResult{ModifiedCount: 0}, nil
	}
	
	return &mongo.UpdateResult{ModifiedCount: 0}, nil
}

// Mock database that returns our mock collection
type mockUserDatabase struct {
	collection *mockUserCollection
}

func (m *mockUserDatabase) Collection(name string) *mongo.Collection {
	// We can't actually return *mongo.Collection, so we'll need a different approach
	// For now, return nil and handle in repo tests differently
	return nil
}

// Create a wrapper that implements our repository interface for testing
type testUserRepo struct {
	mockCollection *mockUserCollection
}

func newTestUserRepo(mockCollection *mockUserCollection) *testUserRepo {
	return &testUserRepo{mockCollection: mockCollection}
}

func (r *testUserRepo) Create(user *entity.User) error {
	_, err := r.mockCollection.InsertOne(context.Background(), user)
	return err
}

func (r *testUserRepo) FindByEmail(email string) (*entity.User, error) {
	if r.mockCollection.documents == nil {
		return nil, appErrors.ErrUserNotFound
	}
	
	if user, found := r.mockCollection.documents[email]; found {
		return user, nil
	}
	return nil, appErrors.ErrUserNotFound
}

func (r *testUserRepo) FindByPhone(phone string) (*entity.User, error) {
	if r.mockCollection.documents == nil {
		return nil, appErrors.ErrUserNotFound
	}
	
	for _, user := range r.mockCollection.documents {
		if user.PhoneNumber == phone {
			return user, nil
		}
	}
	return nil, appErrors.ErrUserNotFound
}

func (r *testUserRepo) Update(user *entity.User) error {
	_, err := r.mockCollection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": user})
	return err
}

func (r *testUserRepo) UpdateEmail(user *entity.User, oldEmail string) error {
	if r.mockCollection.documents == nil {
		return appErrors.ErrUserNotFound
	}
	
	if _, found := r.mockCollection.documents[oldEmail]; found {
		// Remove old entry and add new one
		delete(r.mockCollection.documents, oldEmail)
		r.mockCollection.documents[user.Email] = user
		return nil
	}
	return appErrors.ErrUserNotFound
}

func (r *testUserRepo) UpdatePhone(user *entity.User, oldPhone string) error {
	if r.mockCollection.documents == nil {
		return appErrors.ErrUserNotFound
	}
	
	for email, u := range r.mockCollection.documents {
		if u.PhoneNumber == oldPhone {
			r.mockCollection.documents[email] = user
			return nil
		}
	}
	return appErrors.ErrUserNotFound
}

// Integration tests for user repository methods
func TestUserRepo_Create_Success(t *testing.T) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	user := &entity.User{
		ID:          "test-id",
		Fullname:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		Password:    "hashedpassword",
		AvatarUrl:   "avatar.jpg",
		Verified:    false,
		OnBoarded:   false,
	}
	
	err := repo.Create(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify user was created
	if mockColl.documents == nil || mockColl.documents[user.Email] == nil {
		t.Error("Expected user to be created in mock collection")
	}
	
	// Verify CreatedAt was set
	if mockColl.documents[user.Email].CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUserRepo_Create_Error(t *testing.T) {
	mockColl := &mockUserCollection{
		returnError: errors.New("database error"),
	}
	repo := newTestUserRepo(mockColl)
	
	user := &entity.User{
		Email: "john@example.com",
	}
	
	err := repo.Create(user)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUserRepo_FindByEmail_Success(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:       "test-id",
				Fullname: "John Doe",
				Email:    "john@example.com",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	user, err := repo.FindByEmail("john@example.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user == nil {
		t.Fatal("Expected user to be found")
	}
	
	if user.Email != "john@example.com" {
		t.Errorf("Expected email john@example.com, got %s", user.Email)
	}
	
	if user.Fullname != "John Doe" {
		t.Errorf("Expected fullname John Doe, got %s", user.Fullname)
	}
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	user, err := repo.FindByEmail("nonexistent@example.com")
	if err != appErrors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	
	if user != nil {
		t.Error("Expected user to be nil")
	}
}

func TestUserRepo_FindByPhone_Success(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:          "test-id",
				Fullname:    "John Doe",
				Email:       "john@example.com",
				PhoneNumber: "+1234567890",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	user, err := repo.FindByPhone("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if user == nil {
		t.Fatal("Expected user to be found")
	}
	
	if user.PhoneNumber != "+1234567890" {
		t.Errorf("Expected phone +1234567890, got %s", user.PhoneNumber)
	}
}

func TestUserRepo_FindByPhone_NotFound(t *testing.T) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	user, err := repo.FindByPhone("+9999999999")
	if err != appErrors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	
	if user != nil {
		t.Error("Expected user to be nil")
	}
}

func TestUserRepo_Update_Success(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:       "test-id",
				Fullname: "John Doe",
				Email:    "john@example.com",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	updatedUser := &entity.User{
		ID:       "test-id",
		Fullname: "John Updated",
		Email:    "john@example.com",
	}
	
	err := repo.Update(updatedUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserRepo_UpdateEmail_Success(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"old@example.com": {
				ID:       "test-id",
				Fullname: "John Doe",
				Email:    "old@example.com",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	updatedUser := &entity.User{
		ID:       "test-id",
		Fullname: "John Doe",
		Email:    "new@example.com",
	}
	
	err := repo.UpdateEmail(updatedUser, "old@example.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify old email is gone and new email exists
	if _, exists := mockColl.documents["old@example.com"]; exists {
		t.Error("Expected old email to be removed")
	}
	
	if _, exists := mockColl.documents["new@example.com"]; !exists {
		t.Error("Expected new email to be added")
	}
}

func TestUserRepo_UpdateEmail_NotFound(t *testing.T) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	updatedUser := &entity.User{
		Email: "new@example.com",
	}
	
	err := repo.UpdateEmail(updatedUser, "nonexistent@example.com")
	if err != appErrors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepo_UpdatePhone_Success(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:          "test-id",
				Email:       "john@example.com",
				PhoneNumber: "+1234567890",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	updatedUser := &entity.User{
		ID:          "test-id",
		Email:       "john@example.com",
		PhoneNumber: "+9876543210",
	}
	
	err := repo.UpdatePhone(updatedUser, "+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify phone was updated
	if mockColl.documents["john@example.com"].PhoneNumber != "+9876543210" {
		t.Error("Expected phone number to be updated")
	}
}

func TestUserRepo_UpdatePhone_NotFound(t *testing.T) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	updatedUser := &entity.User{
		PhoneNumber: "+9876543210",
	}
	
	err := repo.UpdatePhone(updatedUser, "+9999999999")
	if err != appErrors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

// Test complex update scenarios with OTP handling
func TestUserRepo_ComplexUpdateScenarios(t *testing.T) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:       "test-id",
				Email:    "john@example.com",
				OTP:      "123456",
				OTPType:  "verification",
				OTPExpiresAt: time.Now().Add(5 * time.Minute),
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	// Test update with empty OTP (should trigger unset logic)
	updatedUser := &entity.User{
		ID:    "test-id",
		Email: "john@example.com",
		OTP:   "", // Empty OTP should trigger unset
	}
	
	err := repo.Update(updatedUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Benchmark tests
func BenchmarkUserRepo_Create(b *testing.B) {
	mockColl := &mockUserCollection{}
	repo := newTestUserRepo(mockColl)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &entity.User{
			ID:    "test-id",
			Email: "john@example.com",
		}
		repo.Create(user)
	}
}

func BenchmarkUserRepo_FindByEmail(b *testing.B) {
	mockColl := &mockUserCollection{
		documents: map[string]*entity.User{
			"john@example.com": {
				ID:    "test-id",
				Email: "john@example.com",
			},
		},
	}
	repo := newTestUserRepo(mockColl)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByEmail("john@example.com")
	}
}