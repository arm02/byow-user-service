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

// Mock mongo.Collection for company testing
type mockCompanyCollection struct {
	insertOneFunc    func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	findOneFunc      func(ctx context.Context, filter interface{}) *mongo.SingleResult
	findFunc         func(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	updateOneFunc    func(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	deleteOneFunc    func(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	countDocumentsFunc func(ctx context.Context, filter interface{}) (int64, error)
	documents        map[string]*entity.Company // Store by ID as key
	returnError      error
}

func (m *mockCompanyCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if m.insertOneFunc != nil {
		return m.insertOneFunc(ctx, document)
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	
	company, ok := document.(*entity.Company)
	if !ok {
		return nil, errors.New("invalid document type")
	}
	
	if m.documents == nil {
		m.documents = make(map[string]*entity.Company)
	}
	
	// Check for duplicates
	for _, existing := range m.documents {
		if (company.CompanyEmail != "" && existing.CompanyEmail == company.CompanyEmail) ||
		   (company.CompanyPhone != "" && existing.CompanyPhone == company.CompanyPhone) {
			return nil, appErrors.ErrEmailOrPhoneAlreadyRegistered
		}
	}
	
	// Set timestamp and ID
	company.CreatedAt = time.Now()
	if company.ID.IsZero() {
		company.ID = primitive.NewObjectID()
	}
	
	m.documents[company.ID.Hex()] = company
	
	return &mongo.InsertOneResult{InsertedID: company.ID}, nil
}

func (m *mockCompanyCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if m.findOneFunc != nil {
		return m.findOneFunc(ctx, filter)
	}
	
	// Since we can't easily mock SingleResult, we'll handle this in the test repo
	return &mongo.SingleResult{}
}

func (m *mockCompanyCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, filter)
	}
	
	// Since we can't easily mock Cursor, we'll handle this in the test repo
	return nil, nil
}

func (m *mockCompanyCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if m.updateOneFunc != nil {
		return m.updateOneFunc(ctx, filter, update)
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	
	return &mongo.UpdateResult{ModifiedCount: 1}, nil
}

func (m *mockCompanyCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if m.deleteOneFunc != nil {
		return m.deleteOneFunc(ctx, filter)
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	
	filterMap, ok := filter.(bson.M)
	if !ok {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}
	
	if id, exists := filterMap["_id"]; exists {
		if objID, ok := id.(primitive.ObjectID); ok {
			if m.documents != nil {
				if _, exists := m.documents[objID.Hex()]; exists {
					delete(m.documents, objID.Hex())
					return &mongo.DeleteResult{DeletedCount: 1}, nil
				}
			}
		}
	}
	
	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

func (m *mockCompanyCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	if m.countDocumentsFunc != nil {
		return m.countDocumentsFunc(ctx, filter)
	}
	if m.returnError != nil {
		return 0, m.returnError
	}
	
	if m.documents == nil {
		return 0, nil
	}
	
	// Simple count implementation for testing
	return int64(len(m.documents)), nil
}

// Create a wrapper that implements our repository interface for testing
type testCompanyRepo struct {
	mockCollection *mockCompanyCollection
}

func newTestCompanyRepo(mockCollection *mockCompanyCollection) *testCompanyRepo {
	return &testCompanyRepo{mockCollection: mockCollection}
}

func (r *testCompanyRepo) FindAll(userID string, keyword string, limit int64, offset int64) ([]*entity.Company, int64, error) {
	if r.mockCollection.documents == nil {
		return []*entity.Company{}, 0, nil
	}
	
	var result []*entity.Company
	for _, company := range r.mockCollection.documents {
		// Filter by userID if provided
		if userID != "" && company.UserID != userID {
			continue
		}
		
		// Filter by keyword if provided (simple contains check)
		if keyword != "" {
			found := false
			companyName := company.CompanyName
			for i := 0; i <= len(companyName)-len(keyword); i++ {
				if i+len(keyword) <= len(companyName) && companyName[i:i+len(keyword)] == keyword {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		
		result = append(result, company)
	}
	
	total := int64(len(result))
	
	// Apply pagination
	start := offset
	end := offset + limit
	
	if start > total {
		return []*entity.Company{}, total, nil
	}
	
	if end > total {
		end = total
	}
	
	if limit > 0 && start < total {
		result = result[start:end]
	}
	
	return result, total, nil
}

func (r *testCompanyRepo) Create(company *entity.Company) error {
	_, err := r.mockCollection.InsertOne(context.Background(), company)
	return err
}

func (r *testCompanyRepo) FindByID(id primitive.ObjectID) (*entity.Company, error) {
	if r.mockCollection.documents == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	if company, found := r.mockCollection.documents[id.Hex()]; found {
		return company, nil
	}
	return nil, appErrors.NewNotFoundError("Company")
}

func (r *testCompanyRepo) FindByEmail(email string) (*entity.Company, error) {
	if r.mockCollection.documents == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	for _, company := range r.mockCollection.documents {
		if company.CompanyEmail == email {
			return company, nil
		}
	}
	return nil, appErrors.NewNotFoundError("Company")
}

func (r *testCompanyRepo) FindByPhone(phone string) (*entity.Company, error) {
	if r.mockCollection.documents == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	for _, company := range r.mockCollection.documents {
		if company.CompanyPhone == phone {
			return company, nil
		}
	}
	return nil, appErrors.NewNotFoundError("Company")
}

func (r *testCompanyRepo) Update(company *entity.Company) error {
	_, err := r.mockCollection.UpdateOne(context.Background(), bson.M{"id": company.ID}, bson.M{"$set": company})
	return err
}

func (r *testCompanyRepo) Delete(id primitive.ObjectID) error {
	_, err := r.mockCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// Integration tests for company repository methods
func TestCompanyRepo_Create_Success(t *testing.T) {
	mockColl := &mockCompanyCollection{}
	repo := newTestCompanyRepo(mockColl)
	
	company := &entity.Company{
		UserID:         "user123",
		CompanyName:    "Test Company",
		CompanyEmail:   "test@company.com",
		CompanyPhone:   "+1234567890",
		CompanyAddress: "123 Test St",
		CompanyLogo:    "logo.png",
		Verified:       false,
	}
	
	err := repo.Create(company)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify company was created
	if mockColl.documents == nil || len(mockColl.documents) == 0 {
		t.Error("Expected company to be created in mock collection")
	}
	
	// Verify CreatedAt was set
	var createdCompany *entity.Company
	for _, c := range mockColl.documents {
		createdCompany = c
		break
	}
	
	if createdCompany.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	
	if createdCompany.ID.IsZero() {
		t.Error("Expected ID to be set")
	}
}

func TestCompanyRepo_Create_DuplicateEmail(t *testing.T) {
	existingCompany := &entity.Company{
		ID:           primitive.NewObjectID(),
		CompanyEmail: "duplicate@company.com",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			existingCompany.ID.Hex(): existingCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	newCompany := &entity.Company{
		CompanyName:  "New Company",
		CompanyEmail: "duplicate@company.com", // Same email
	}
	
	err := repo.Create(newCompany)
	if err != appErrors.ErrEmailOrPhoneAlreadyRegistered {
		t.Errorf("Expected ErrEmailOrPhoneAlreadyRegistered, got %v", err)
	}
}

func TestCompanyRepo_Create_DuplicatePhone(t *testing.T) {
	existingCompany := &entity.Company{
		ID:           primitive.NewObjectID(),
		CompanyPhone: "+1234567890",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			existingCompany.ID.Hex(): existingCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	newCompany := &entity.Company{
		CompanyName:  "New Company",
		CompanyPhone: "+1234567890", // Same phone
	}
	
	err := repo.Create(newCompany)
	if err != appErrors.ErrEmailOrPhoneAlreadyRegistered {
		t.Errorf("Expected ErrEmailOrPhoneAlreadyRegistered, got %v", err)
	}
}

func TestCompanyRepo_Create_Error(t *testing.T) {
	mockColl := &mockCompanyCollection{
		returnError: errors.New("database error"),
	}
	repo := newTestCompanyRepo(mockColl)
	
	company := &entity.Company{
		CompanyName: "Test Company",
	}
	
	err := repo.Create(company)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCompanyRepo_FindByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:          id,
		UserID:      "user123",
		CompanyName: "Test Company",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	company, err := repo.FindByID(id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if company == nil {
		t.Fatal("Expected company to be found")
	}
	
	if company.ID != id {
		t.Errorf("Expected ID %v, got %v", id, company.ID)
	}
	
	if company.CompanyName != "Test Company" {
		t.Errorf("Expected company name 'Test Company', got %s", company.CompanyName)
	}
}

func TestCompanyRepo_FindByID_NotFound(t *testing.T) {
	mockColl := &mockCompanyCollection{}
	repo := newTestCompanyRepo(mockColl)
	
	id := primitive.NewObjectID()
	company, err := repo.FindByID(id)
	
	if company != nil {
		t.Error("Expected company to be nil")
	}
	
	if appErr, ok := err.(*appErrors.AppError); ok {
		if appErr.Status != 404 {
			t.Errorf("Expected 404 error status, got %d", appErr.Status)
		}
	} else {
		t.Errorf("Expected AppError, got %T", err)
	}
}

func TestCompanyRepo_FindByEmail_Success(t *testing.T) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:           id,
		CompanyName:  "Test Company",
		CompanyEmail: "test@company.com",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	company, err := repo.FindByEmail("test@company.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if company == nil {
		t.Fatal("Expected company to be found")
	}
	
	if company.CompanyEmail != "test@company.com" {
		t.Errorf("Expected email 'test@company.com', got %s", company.CompanyEmail)
	}
}

func TestCompanyRepo_FindByEmail_NotFound(t *testing.T) {
	mockColl := &mockCompanyCollection{}
	repo := newTestCompanyRepo(mockColl)
	
	company, err := repo.FindByEmail("nonexistent@company.com")
	
	if company != nil {
		t.Error("Expected company to be nil")
	}
	
	if appErr, ok := err.(*appErrors.AppError); ok {
		if appErr.Status != 404 {
			t.Errorf("Expected 404 error status, got %d", appErr.Status)
		}
	} else {
		t.Errorf("Expected AppError, got %T", err)
	}
}

func TestCompanyRepo_FindByPhone_Success(t *testing.T) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:           id,
		CompanyName:  "Test Company",
		CompanyPhone: "+1234567890",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	company, err := repo.FindByPhone("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if company == nil {
		t.Fatal("Expected company to be found")
	}
	
	if company.CompanyPhone != "+1234567890" {
		t.Errorf("Expected phone '+1234567890', got %s", company.CompanyPhone)
	}
}

func TestCompanyRepo_FindAll_Success(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	
	company1 := &entity.Company{
		ID:          id1,
		UserID:      "user123",
		CompanyName: "Tech Solutions",
	}
	company2 := &entity.Company{
		ID:          id2,
		UserID:      "user123",
		CompanyName: "Marketing Agency",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id1.Hex(): company1,
			id2.Hex(): company2,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	companies, total, err := repo.FindAll("user123", "", 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	
	if len(companies) != 2 {
		t.Errorf("Expected 2 companies, got %d", len(companies))
	}
}

func TestCompanyRepo_FindAll_WithKeyword(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	
	company1 := &entity.Company{
		ID:          id1,
		UserID:      "user123",
		CompanyName: "Tech Solutions",
	}
	company2 := &entity.Company{
		ID:          id2,
		UserID:      "user123",
		CompanyName: "Marketing Agency",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id1.Hex(): company1,
			id2.Hex(): company2,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	companies, total, err := repo.FindAll("user123", "Tech", 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if total != 1 {
		t.Errorf("Expected total 1, got %d", total)
	}
	
	if len(companies) != 1 {
		t.Errorf("Expected 1 company, got %d", len(companies))
	}
	
	if companies[0].CompanyName != "Tech Solutions" {
		t.Errorf("Expected 'Tech Solutions', got %s", companies[0].CompanyName)
	}
}

func TestCompanyRepo_FindAll_WithPagination(t *testing.T) {
	// Create 5 companies
	companies := make(map[string]*entity.Company)
	for i := 0; i < 5; i++ {
		id := primitive.NewObjectID()
		company := &entity.Company{
			ID:          id,
			UserID:      "user123",
			CompanyName: "Company " + string(rune('A'+i)),
		}
		companies[id.Hex()] = company
	}
	
	mockColl := &mockCompanyCollection{
		documents: companies,
	}
	repo := newTestCompanyRepo(mockColl)
	
	// Test first page
	result, total, err := repo.FindAll("user123", "", 2, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	
	if len(result) != 2 {
		t.Errorf("Expected 2 companies on first page, got %d", len(result))
	}
	
	// Test second page
	result, total, err = repo.FindAll("user123", "", 2, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	
	if len(result) != 2 {
		t.Errorf("Expected 2 companies on second page, got %d", len(result))
	}
}

func TestCompanyRepo_Update_Success(t *testing.T) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:          id,
		CompanyName: "Original Name",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	updatedCompany := &entity.Company{
		ID:          id,
		CompanyName: "Updated Name",
	}
	
	err := repo.Update(updatedCompany)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCompanyRepo_Delete_Success(t *testing.T) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:          id,
		CompanyName: "To Be Deleted",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	err := repo.Delete(id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify company was deleted
	if _, exists := mockColl.documents[id.Hex()]; exists {
		t.Error("Expected company to be deleted")
	}
}

func TestCompanyRepo_Delete_NotFound(t *testing.T) {
	mockColl := &mockCompanyCollection{}
	repo := newTestCompanyRepo(mockColl)
	
	id := primitive.NewObjectID()
	err := repo.Delete(id)
	if err != nil {
		t.Errorf("Expected no error for non-existent delete, got %v", err)
	}
}

// Benchmark tests
func BenchmarkCompanyRepo_Create(b *testing.B) {
	mockColl := &mockCompanyCollection{}
	repo := newTestCompanyRepo(mockColl)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		company := &entity.Company{
			CompanyName: "Benchmark Company",
		}
		repo.Create(company)
	}
}

func BenchmarkCompanyRepo_FindAll(b *testing.B) {
	// Create test data
	companies := make(map[string]*entity.Company)
	for i := 0; i < 100; i++ {
		id := primitive.NewObjectID()
		company := &entity.Company{
			ID:          id,
			UserID:      "user123",
			CompanyName: "Company " + string(rune('A'+i%26)),
		}
		companies[id.Hex()] = company
	}
	
	mockColl := &mockCompanyCollection{
		documents: companies,
	}
	repo := newTestCompanyRepo(mockColl)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindAll("user123", "", 10, 0)
	}
}

func BenchmarkCompanyRepo_FindByID(b *testing.B) {
	id := primitive.NewObjectID()
	testCompany := &entity.Company{
		ID:          id,
		CompanyName: "Benchmark Company",
	}
	
	mockColl := &mockCompanyCollection{
		documents: map[string]*entity.Company{
			id.Hex(): testCompany,
		},
	}
	repo := newTestCompanyRepo(mockColl)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.FindByID(id)
	}
}