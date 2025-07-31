package usecase

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock company repository for testing
type mockCompanyRepository struct {
	companies map[string]*entity.Company
	nextID    int
}

func (m *mockCompanyRepository) FindAll(userID, keyword string, limit, offset int64) ([]*entity.Company, int64, error) {
	if m.companies == nil {
		return []*entity.Company{}, 0, nil
	}
	
	var result []*entity.Company
	for _, company := range m.companies {
		// Filter by user ID if provided
		if userID != "" && company.UserID != userID {
			continue
		}
		
		// Filter by keyword if provided (case-insensitive partial match)
		if keyword != "" {
			// Simple contains check for testing
			companyNameLower := company.CompanyName
			keywordLower := keyword
			if len(companyNameLower) > 0 && len(keywordLower) > 0 {
				found := false
				for i := 0; i <= len(companyNameLower)-len(keywordLower); i++ {
					if companyNameLower[i:i+len(keywordLower)] == keywordLower {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		
		result = append(result, company)
	}
	
	// Apply pagination
	total := int64(len(result))
	start := offset
	end := offset + limit
	
	if start > total {
		return []*entity.Company{}, total, nil
	}
	
	if end > total {
		end = total
	}
	
	if limit > 0 {
		result = result[start:end]
	}
	
	return result, total, nil
}

func (m *mockCompanyRepository) Create(company *entity.Company) error {
	if m.companies == nil {
		m.companies = make(map[string]*entity.Company)
	}
	
	// Check for duplicates
	for _, existing := range m.companies {
		if (company.CompanyEmail != "" && existing.CompanyEmail == company.CompanyEmail) ||
			(company.CompanyPhone != "" && existing.CompanyPhone == company.CompanyPhone) {
			return appErrors.ErrEmailOrPhoneAlreadyRegistered
		}
	}
	
	// Generate ID and set timestamp
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	
	// Use a unique key for storage
	key := company.ID.Hex()
	m.companies[key] = company
	
	return nil
}

func (m *mockCompanyRepository) FindByID(id primitive.ObjectID) (*entity.Company, error) {
	if m.companies == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	key := id.Hex()
	if company, exists := m.companies[key]; exists {
		return company, nil
	}
	
	return nil, appErrors.NewNotFoundError("Company")
}

func (m *mockCompanyRepository) FindByEmail(email string) (*entity.Company, error) {
	if m.companies == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	for _, company := range m.companies {
		if company.CompanyEmail == email {
			return company, nil
		}
	}
	
	return nil, appErrors.NewNotFoundError("Company")
}

func (m *mockCompanyRepository) FindByPhone(phone string) (*entity.Company, error) {
	if m.companies == nil {
		return nil, appErrors.NewNotFoundError("Company")
	}
	
	for _, company := range m.companies {
		if company.CompanyPhone == phone {
			return company, nil
		}
	}
	
	return nil, appErrors.NewNotFoundError("Company")
}

func (m *mockCompanyRepository) Update(company *entity.Company) error {
	if m.companies == nil {
		return appErrors.NewNotFoundError("Company")
	}
	
	key := company.ID.Hex()
	if _, exists := m.companies[key]; exists {
		m.companies[key] = company
		return nil
	}
	
	return appErrors.NewNotFoundError("Company")
}

func (m *mockCompanyRepository) Delete(id primitive.ObjectID) error {
	if m.companies == nil {
		return appErrors.NewNotFoundError("Company")
	}
	
	key := id.Hex()
	if _, exists := m.companies[key]; exists {
		delete(m.companies, key)
		return nil
	}
	
	return appErrors.NewNotFoundError("Company")
}

// Mock function to extract user ID from context
func mockUserIDFunc(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return "test-user-123"
}

func setupCompanyUsecase() *CompanyUsecase {
	return &CompanyUsecase{
		Repo:   &mockCompanyRepository{},
		UserID: mockUserIDFunc,
	}
}

func setupGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("user_id", "test-user-123")
	return c
}

func TestCompanyUsecase_GetAll_Success(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create test companies
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	
	company1 := &entity.Company{
		ID:           primitive.NewObjectID(),
		UserID:       "test-user-123",
		CompanyName:  "Test Company 1",
		CompanyEmail: "test1@company.com",
		CreatedAt:    time.Now(),
	}
	company2 := &entity.Company{
		ID:           primitive.NewObjectID(), 
		UserID:       "test-user-123",
		CompanyName:  "Another Company",
		CompanyEmail: "test2@company.com",
		CreatedAt:    time.Now(),
	}
	
	repo.companies[company1.ID.Hex()] = company1
	repo.companies[company2.ID.Hex()] = company2
	
	responses, count, err := uc.GetAll(c, "", 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if responses == nil {
		t.Fatal("Expected responses to be non-nil")
	}
	
	if len(*responses) != 2 {
		t.Errorf("Expected 2 companies, got %d", len(*responses))
	}
	
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
	
	// Check response structure
	response := (*responses)[0]
	if response.UserID == "" {
		t.Error("Expected UserID to be set")
	}
	
	if response.CompanyID.IsZero() {
		t.Error("Expected CompanyID to be set")
	}
	
	if response.CreatedAt == "" {
		t.Error("Expected CreatedAt to be formatted")
	}
}

func TestCompanyUsecase_GetAll_WithKeyword(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create test companies
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	
	company1 := &entity.Company{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-123",
		CompanyName: "Tech Solutions",
		CreatedAt:   time.Now(),
	}
	company2 := &entity.Company{
		ID:          primitive.NewObjectID(),
		UserID:      "test-user-123", 
		CompanyName: "Marketing Agency",
		CreatedAt:   time.Now(),
	}
	
	repo.companies[company1.ID.Hex()] = company1
	repo.companies[company2.ID.Hex()] = company2
	
	responses, count, err := uc.GetAll(c, "Tech", 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if count != 1 {
		t.Errorf("Expected count 1 with keyword filter, got %d", count)
	}
	
	if len(*responses) != 1 {
		t.Errorf("Expected 1 company with keyword filter, got %d", len(*responses))
	}
}

func TestCompanyUsecase_GetAll_WithPagination(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create test companies
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	
	for i := 0; i < 5; i++ {
		company := &entity.Company{
			ID:          primitive.NewObjectID(),
			UserID:      "test-user-123",
			CompanyName: "Company " + string(rune('A'+i)),
			CreatedAt:   time.Now(),
		}
		repo.companies[company.ID.Hex()] = company
	}
	
	// Test first page
	responses, count, err := uc.GetAll(c, "", 2, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if count != 5 {
		t.Errorf("Expected total count 5, got %d", count)
	}
	
	if len(*responses) != 2 {
		t.Errorf("Expected 2 companies on first page, got %d", len(*responses))
	}
	
	// Test second page
	responses, count, err = uc.GetAll(c, "", 2, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if count != 5 {
		t.Errorf("Expected total count 5, got %d", count)
	}
	
	if len(*responses) != 2 {
		t.Errorf("Expected 2 companies on second page, got %d", len(*responses))
	}
}

func TestCompanyUsecase_GetAll_EmptyResult(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	responses, count, err := uc.GetAll(c, "", 10, 0)
	if err != nil {
		t.Errorf("Expected no error for empty result, got %v", err)
	}
	
	if responses == nil {
		t.Fatal("Expected responses to be non-nil even when empty")
	}
	
	if len(*responses) != 0 {
		t.Errorf("Expected 0 companies, got %d", len(*responses))
	}
	
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestCompanyUsecase_Create_Success(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	req := dto.CompanyRequest{
		CompanyName:    "New Company",
		CompanyEmail:   "new@company.com",
		CompanyPhone:   "+1234567890",
		CompanyAddress: "123 Main St",
		CompanyLogo:    "logo.png",
	}
	
	company, err := uc.Create(c, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if company == nil {
		t.Fatal("Expected company to be created")
	}
	
	if company.CompanyName != req.CompanyName {
		t.Errorf("Expected company name %s, got %s", req.CompanyName, company.CompanyName)
	}
	
	if company.CompanyEmail != req.CompanyEmail {
		t.Errorf("Expected company email %s, got %s", req.CompanyEmail, company.CompanyEmail)
	}
	
	if company.UserID != "test-user-123" {
		t.Errorf("Expected user ID test-user-123, got %s", company.UserID)
	}
	
	if company.Verified {
		t.Error("Expected company to be unverified initially")
	}
	
	if company.ID.IsZero() {
		t.Error("Expected company ID to be set")
	}
	
	if company.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestCompanyUsecase_Create_DuplicateEmail(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create first company
	req1 := dto.CompanyRequest{
		CompanyName:  "Company 1",
		CompanyEmail: "duplicate@company.com",
	}
	
	_, err := uc.Create(c, req1)
	if err != nil {
		t.Fatalf("Expected no error creating first company, got %v", err)
	}
	
	// Try to create second company with same email
	req2 := dto.CompanyRequest{
		CompanyName:  "Company 2", 
		CompanyEmail: "duplicate@company.com",
	}
	
	_, err = uc.Create(c, req2)
	if err != appErrors.ErrEmailOrPhoneAlreadyRegistered {
		t.Errorf("Expected ErrEmailOrPhoneAlreadyRegistered, got %v", err)
	}
}

func TestCompanyUsecase_Create_DuplicatePhone(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create first company
	req1 := dto.CompanyRequest{
		CompanyName:  "Company 1",
		CompanyPhone: "+1234567890",
	}
	
	_, err := uc.Create(c, req1)
	if err != nil {
		t.Fatalf("Expected no error creating first company, got %v", err)
	}
	
	// Try to create second company with same phone
	req2 := dto.CompanyRequest{
		CompanyName:  "Company 2",
		CompanyPhone: "+1234567890",
	}
	
	_, err = uc.Create(c, req2)
	if err != appErrors.ErrEmailOrPhoneAlreadyRegistered {
		t.Errorf("Expected ErrEmailOrPhoneAlreadyRegistered, got %v", err)
	}
}

func TestCompanyUsecase_FindByID_Success(t *testing.T) {
	uc := setupCompanyUsecase()
	
	// Create a company first
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	
	originalCompany := &entity.Company{
		ID:           primitive.NewObjectID(),
		UserID:       "test-user-123",
		CompanyName:  "Test Company",
		CompanyEmail: "test@company.com",
		CreatedAt:    time.Now(),
	}
	
	repo.companies[originalCompany.ID.Hex()] = originalCompany
	
	// Find by ID
	company, err := uc.FindByID(originalCompany.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if company == nil {
		t.Fatal("Expected company to be found")
	}
	
	if company.ID != originalCompany.ID {
		t.Errorf("Expected company ID %v, got %v", originalCompany.ID, company.ID)
	}
	
	if company.CompanyName != originalCompany.CompanyName {
		t.Errorf("Expected company name %s, got %s", originalCompany.CompanyName, company.CompanyName)
	}
}

func TestCompanyUsecase_FindByID_NotFound(t *testing.T) {
	uc := setupCompanyUsecase()
	
	nonExistentID := primitive.NewObjectID()
	
	_, err := uc.FindByID(nonExistentID)
	if err == nil {
		t.Error("Expected error for non-existent company")
	}
	
	// The error should be a NotFound error
	if appErr, ok := err.(*appErrors.AppError); ok {
		if appErr.Status != 404 {
			t.Errorf("Expected 404 error status, got %d", appErr.Status)
		}
	}
}

func TestCompanyUsecase_UserIDExtraction(t *testing.T) {
	uc := setupCompanyUsecase()
	
	// Test with context that has user_id
	c := setupGinContext()
	c.Set("user_id", "custom-user-456")
	
	userID := uc.UserID(c)
	if userID != "custom-user-456" {
		t.Errorf("Expected user ID custom-user-456, got %s", userID)
	}
	
	// Test with context that doesn't have user_id (should return default)
	c2, _ := gin.CreateTestContext(nil)
	userID2 := uc.UserID(c2)
	if userID2 != "test-user-123" {
		t.Errorf("Expected default user ID test-user-123, got %s", userID2)
	}
}

func TestCompanyUsecase_ResponseMapping(t *testing.T) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Create a company with all fields
	testTime := time.Now()
	company := &entity.Company{
		ID:             primitive.NewObjectID(),
		UserID:         "test-user-123",
		CompanyName:    "Full Company",
		CompanyEmail:   "full@company.com",
		CompanyPhone:   "+1234567890",
		CompanyAddress: "123 Full Street",
		CompanyLogo:    "full-logo.png",
		Verified:       true,
		CreatedAt:      testTime,
	}
	
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	repo.companies[company.ID.Hex()] = company
	
	responses, _, err := uc.GetAll(c, "", 10, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(*responses) != 1 {
		t.Fatalf("Expected 1 response, got %d", len(*responses))
	}
	
	response := (*responses)[0]
	
	// Verify all fields are mapped correctly
	if response.UserID != company.UserID {
		t.Errorf("Expected UserID %s, got %s", company.UserID, response.UserID)
	}
	
	if response.CompanyID != company.ID {
		t.Errorf("Expected CompanyID %v, got %v", company.ID, response.CompanyID)
	}
	
	if response.CompanyName != company.CompanyName {
		t.Errorf("Expected CompanyName %s, got %s", company.CompanyName, response.CompanyName)
	}
	
	if response.CompanyEmail != company.CompanyEmail {
		t.Errorf("Expected CompanyEmail %s, got %s", company.CompanyEmail, response.CompanyEmail)
	}
	
	if response.CompanyPhone != company.CompanyPhone {
		t.Errorf("Expected CompanyPhone %s, got %s", company.CompanyPhone, response.CompanyPhone)
	}
	
	if response.CompanyAddress != company.CompanyAddress {
		t.Errorf("Expected CompanyAddress %s, got %s", company.CompanyAddress, response.CompanyAddress)
	}
	
	if response.CompanyLogo != company.CompanyLogo {
		t.Errorf("Expected CompanyLogo %s, got %s", company.CompanyLogo, response.CompanyLogo)
	}
	
	if response.Verified != company.Verified {
		t.Errorf("Expected Verified %v, got %v", company.Verified, response.Verified)
	}
	
	// Check time formatting
	expectedTime := company.CreatedAt.Format(time.RFC3339)
	if response.CreatedAt != expectedTime {
		t.Errorf("Expected CreatedAt %s, got %s", expectedTime, response.CreatedAt)
	}
}

func TestCompanyUsecase_JSONSerialization(t *testing.T) {
	// Test that DTO responses can be properly serialized to JSON
	response := dto.CompanyResponse{
		UserID:         "user123",
		CompanyID:      primitive.NewObjectID(),
		CompanyName:    "JSON Test Company",
		CompanyEmail:   "json@test.com",
		CompanyPhone:   "+1234567890",
		CompanyAddress: "JSON Street",
		CompanyLogo:    "json-logo.png",
		Verified:       true,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
	
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Expected no error marshaling to JSON, got %v", err)
	}
	
	if len(jsonData) == 0 {
		t.Error("Expected non-empty JSON data")
	}
	
	// Test unmarshaling
	var unmarshaled dto.CompanyResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Expected no error unmarshaling from JSON, got %v", err)
	}
	
	if unmarshaled.CompanyName != response.CompanyName {
		t.Errorf("Expected company name %s after JSON round-trip, got %s", response.CompanyName, unmarshaled.CompanyName)
	}
}

// Test struct initialization
func TestCompanyUsecaseStruct(t *testing.T) {
	uc := &CompanyUsecase{
		UserID: mockUserIDFunc,
	}
	
	if uc.UserID == nil {
		t.Error("Expected UserID function to be set")
	}
	
	// Test that UserID function works
	c := setupGinContext()
	userID := uc.UserID(c)
	if userID == "" {
		t.Error("Expected non-empty user ID from function")
	}
}

// Benchmark tests
func BenchmarkCompanyUsecase_GetAll(b *testing.B) {
	uc := setupCompanyUsecase()
	c := setupGinContext()
	
	// Setup test data
	repo := uc.Repo.(*mockCompanyRepository)
	repo.companies = make(map[string]*entity.Company)
	
	for i := 0; i < 100; i++ {
		company := &entity.Company{
			ID:          primitive.NewObjectID(),
			UserID:      "test-user-123",
			CompanyName: "Benchmark Company",
			CreatedAt:   time.Now(),
		}
		repo.companies[company.ID.Hex()] = company
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.GetAll(c, "", 10, 0)
	}
}

func BenchmarkCompanyUsecase_Create(b *testing.B) {
	req := dto.CompanyRequest{
		CompanyName:  "Benchmark Company",
		CompanyEmail: "benchmark@company.com",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		uc := setupCompanyUsecase()
		c := setupGinContext()
		b.StartTimer()
		
		uc.Create(c, req)
	}
}