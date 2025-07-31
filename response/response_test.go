package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buildyow/byow-user-service/constants"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccess(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		Success(c, 200, "test data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != constants.SUCCESS {
		t.Errorf("Expected status %v, got %v", constants.SUCCESS, response["status"])
	}

	if response["code"] != float64(200) {
		t.Errorf("Expected code 200, got %v", response["code"])
	}

	if response["response"] != "test data" {
		t.Errorf("Expected response 'test data', got %v", response["response"])
	}
}

func TestSuccessWithPagination(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		SuccessWithPagination(c, 200, []string{"item1", "item2"}, 2)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["row_count"] != float64(2) {
		t.Errorf("Expected row_count 2, got %v", response["row_count"])
	}
}

func TestSuccessWithMessage(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		SuccessWithMessage(c, 200, "Success message")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["response"] != "Success message" {
		t.Errorf("Expected response 'Success message', got %v", response["response"])
	}
}

func TestCreated(t *testing.T) {
	router := setupTestRouter()
	
	router.POST("/test", func(c *gin.Context) {
		Created(c, "created data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}
}

func TestCreatedWithMessage(t *testing.T) {
	router := setupTestRouter()
	
	router.POST("/test", func(c *gin.Context) {
		CreatedWithMessage(c, "Resource created")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}
}

func TestOK(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		OK(c, "ok data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

func TestOKWithMessage(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		OKWithMessage(c, "OK message")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

func TestSpecificSuccessResponses(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*gin.Context)
		expected string
	}{
		{"LogoutSuccess", LogoutSuccess, constants.LOGOUT_SUCCESSFUL},
		{"OnboardSuccess", OnboardSuccess, constants.ONBOARD_SUCCESSFUL},
		{"PasswordChangeSuccess", PasswordChangeSuccess, constants.PASSWORD_CHANGED_SUCCESS},
		{"EmailChangeSuccess", EmailChangeSuccess, constants.EMAIL_CHANGED_SUCCESS},
		{"PhoneChangeSuccess", PhoneChangeSuccess, constants.PHONE_CHANGED_SUCCESS},
		{"OTPVerifiedSuccess", OTPVerifiedSuccess, constants.OTP_VERIFIED},
		{"OTPSentSuccess", OTPSentSuccess, constants.OTP_SENT},
		{"ValidTokenSuccess", ValidTokenSuccess, constants.VALID_TOKEN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status code 200, got %d", w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response["response"] != tt.expected {
				t.Errorf("Expected response '%v', got %v", tt.expected, response["response"])
			}
		})
	}
}

func TestGeneral(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		General(c, 200, "Test message", "test data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	responseData := response["response"].(map[string]interface{})
	if responseData["message"] != "Test message" {
		t.Errorf("Expected message 'Test message', got %v", responseData["message"])
	}

	if responseData["data"] != "test data" {
		t.Errorf("Expected data 'test data', got %v", responseData["data"])
	}
}

func TestGeneralWithoutData(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		General(c, 200, "Test message", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	responseData := response["response"].(map[string]interface{})
	if responseData["message"] != "Test message" {
		t.Errorf("Expected message 'Test message', got %v", responseData["message"])
	}

	if _, exists := responseData["data"]; exists {
		t.Error("Expected no data field when data is nil")
	}
}

func TestGeneralHelpers(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*gin.Context)
		expected int
	}{
		{"GeneralOK", func(c *gin.Context) { GeneralOK(c, "OK message", "data") }, 200},
		{"GeneralCreated", func(c *gin.Context) { GeneralCreated(c, "Created message", "data") }, 201},
		{"GeneralMessage", func(c *gin.Context) { GeneralMessage(c, 200, "Message only") }, 200},
		{"GeneralData", func(c *gin.Context) { GeneralData(c, 200, "data only") }, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, w.Code)
			}
		})
	}
}

func TestCRUDHelpers(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(*gin.Context)
		expected int
	}{
		{"CreateSuccess", func(c *gin.Context) { CreateSuccess(c, "User", "user data") }, 201},
		{"UpdateSuccess", func(c *gin.Context) { UpdateSuccess(c, "User", "user data") }, 200},
		{"DeleteSuccess", func(c *gin.Context) { DeleteSuccess(c, "User") }, 200},
		{"FetchSuccess", func(c *gin.Context) { FetchSuccess(c, "User", "user data") }, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, w.Code)
			}
		})
	}
}

func TestListSuccess(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		ListSuccess(c, "Users", []string{"user1", "user2"}, 2)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	responseData := response["response"].(map[string]interface{})
	if responseData["row_count"] != float64(2) {
		t.Errorf("Expected row_count 2, got %v", responseData["row_count"])
	}
}

func TestError(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		Error(c, 400, "Bad request")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != constants.ERROR {
		t.Errorf("Expected status %v, got %v", constants.ERROR, response["status"])
	}

	data := response["data"].(map[string]interface{})
	if data["message"] != "Bad request" {
		t.Errorf("Expected message 'Bad request', got %v", data["message"])
	}
}

func TestErrorFromAppError(t *testing.T) {
	router := setupTestRouter()
	
	// Test with AppError
	router.GET("/test-app-error", func(c *gin.Context) {
		err := appErrors.ErrUserNotFound
		ErrorFromAppError(c, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-app-error", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	errorData := response["error"].(map[string]interface{})
	if errorData["code"] != "NOT_FOUND" {
		t.Errorf("Expected error code 'NOT_FOUND', got %v", errorData["code"])
	}
}

func TestErrorFromStandardError(t *testing.T) {
	router := setupTestRouter()
	
	// Test with standard error
	router.GET("/test-std-error", func(c *gin.Context) {
		err := errors.New("standard error") 
		ErrorFromAppError(c, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-std-error", nil)
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status code 500, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	if data["message"] != "standard error" {
		t.Errorf("Expected message 'standard error', got %v", data["message"])
	}
}

func TestValidationError(t *testing.T) {
	router := setupTestRouter()
	
	router.GET("/test", func(c *gin.Context) {
		errors := []map[string]string{
			{"field": "email", "message": "Invalid email"},
			{"field": "password", "message": "Password too short"},
		}
		ValidationError(c, errors)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	errorData := response["error"].(map[string]interface{})
	if errorData["code"] != "VALIDATION_ERROR" {
		t.Errorf("Expected error code 'VALIDATION_ERROR', got %v", errorData["code"])
	}

	if errorData["message"] != "Validation failed" {
		t.Errorf("Expected error message 'Validation failed', got %v", errorData["message"])
	}

	details := errorData["details"].([]interface{})
	if len(details) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(details))
	}
}

func TestSuccessResponse(t *testing.T) {
	// Test SuccessResponse struct
	response := SuccessResponse{
		Message: "Test message",
		Data:    "test data",
	}

	if response.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %v", response.Message)
	}

	if response.Data != "test data" {
		t.Errorf("Expected data 'test data', got %v", response.Data)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessResponse: %v", err)
	}

	var unmarshaled SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SuccessResponse: %v", err)
	}

	if unmarshaled.Message != response.Message {
		t.Errorf("Expected message '%v', got %v", response.Message, unmarshaled.Message)  
	}
}