package validation

import (
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/buildyow/byow-user-service/response"
	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

// ValidatePassword validates password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters long"
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// ValidatePhoneNumber validates phone number format
func ValidatePhoneNumber(phone string) bool {
	// Remove all non-digit characters for validation
	phoneDigits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// Check if it's a valid length (8-15 digits as per E.164)
	if len(phoneDigits) < 8 || len(phoneDigits) > 15 {
		return false
	}
	
	// Check if it starts with country code or local format
	phoneRegex := regexp.MustCompile(`^(\+?[1-9]\d{1,14}|0\d{7,14})$`)
	return phoneRegex.MatchString(phone)
}

// ValidateFullName validates full name
func ValidateFullName(name string) (bool, string) {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return false, "Full name must be at least 2 characters long"
	}
	if len(name) > 100 {
		return false, "Full name must be less than 100 characters long"
	}
	
	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-'\.]+$`)
	if !nameRegex.MatchString(name) {
		return false, "Full name can only contain letters, spaces, hyphens, apostrophes, and periods"
	}
	
	return true, ""
}

// ValidateRegistrationRequest validates registration form data
func ValidateRegistrationRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errors []ValidationError

		fullName := strings.TrimSpace(c.PostForm("full_name"))
		email := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")
		phoneNumber := strings.TrimSpace(c.PostForm("phone_number"))

		// Validate full name
		if fullName == "" {
			errors = append(errors, ValidationError{Field: "full_name", Message: "Full name is required"})
		} else {
			if valid, msg := ValidateFullName(fullName); !valid {
				errors = append(errors, ValidationError{Field: "full_name", Message: msg})
			}
		}

		// Validate email
		if email == "" {
			errors = append(errors, ValidationError{Field: "email", Message: "Email is required"})
		} else if !ValidateEmail(email) {
			errors = append(errors, ValidationError{Field: "email", Message: "Invalid email format"})
		}

		// Validate password
		if password == "" {
			errors = append(errors, ValidationError{Field: "password", Message: "Password is required"})
		} else {
			if valid, msg := ValidatePassword(password); !valid {
				errors = append(errors, ValidationError{Field: "password", Message: msg})
			}
		}

		// Validate phone number
		if phoneNumber == "" {
			errors = append(errors, ValidationError{Field: "phone_number", Message: "Phone number is required"})
		} else if !ValidatePhoneNumber(phoneNumber) {
			errors = append(errors, ValidationError{Field: "phone_number", Message: "Invalid phone number format"})
		}

		if len(errors) > 0 {
			response.ValidationError(c, errors)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateLoginRequest validates login JSON data
func ValidateLoginRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid JSON format")
			c.Abort()
			return
		}

		var errors []ValidationError

		email := strings.TrimSpace(req.Email)
		password := req.Password

		// Validate email
		if email == "" {
			errors = append(errors, ValidationError{Field: "email", Message: "Email is required"})
		} else if !ValidateEmail(email) {
			errors = append(errors, ValidationError{Field: "email", Message: "Invalid email format"})
		}

		// Validate password
		if password == "" {
			errors = append(errors, ValidationError{Field: "password", Message: "Password is required"})
		}

		if len(errors) > 0 {
			response.ValidationError(c, errors)
			c.Abort()
			return
		}

		// Store validated data in context for handler
		c.Set("validated_email", email)
		c.Set("validated_password", password)

		c.Next()
	}
}

// ValidateFileUpload validates file upload constraints
func ValidateFileUpload(maxSize int64, allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("avatar")
		if err != nil {
			// File is optional, continue if no file provided
			if err == http.ErrMissingFile {
				c.Next()
				return
			}
			response.Error(c, http.StatusBadRequest, "Error processing file upload")
			c.Abort()
			return
		}
		defer file.Close()

		// Check file size
		if header.Size > maxSize {
			response.Error(c, http.StatusBadRequest, "File size exceeds maximum allowed size")
			c.Abort()
			return
		}

		// Check file type
		contentType := header.Header.Get("Content-Type")
		validType := false
		for _, allowedType := range allowedTypes {
			if strings.Contains(contentType, allowedType) {
				validType = true
				break
			}
		}

		if !validType {
			response.Error(c, http.StatusBadRequest, "Invalid file type. Only images are allowed")
			c.Abort()
			return
		}

		c.Next()
	}
}