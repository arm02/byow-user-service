package response

import (
	"fmt"

	"github.com/buildyow/byow-user-service/constants"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"status":   constants.SUCCESS,
		"code":     code,
		"response": data,
	})
}

func SuccessWithPagination(c *gin.Context, code int, data interface{}, total int64) {
	c.JSON(code, gin.H{
		"status":    constants.SUCCESS,
		"code":      code,
		"response":  data,
		"row_count": total,
	})
}

// Common success response helpers for standardized messages
func SuccessWithMessage(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":   constants.SUCCESS,
		"code":     code,
		"response": message,
	})
}

func Created(c *gin.Context, data interface{}) {
	Success(c, 201, data)
}

func CreatedWithMessage(c *gin.Context, message string) {
	SuccessWithMessage(c, 201, message)
}

func OK(c *gin.Context, data interface{}) {
	Success(c, 200, data)
}

func OKWithMessage(c *gin.Context, message string) {
	SuccessWithMessage(c, 200, message)
}

// Specific success responses using constants
func LogoutSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.LOGOUT_SUCCESSFUL)
}

func OnboardSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.ONBOARD_SUCCESSFUL)
}

func PasswordChangeSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.PASSWORD_CHANGED_SUCCESS)
}

func EmailChangeSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.EMAIL_CHANGED_SUCCESS)
}

func PhoneChangeSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.PHONE_CHANGED_SUCCESS)
}

func OTPVerifiedSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.OTP_VERIFIED)
}

func OTPSentSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.OTP_SENT)
}

func ValidTokenSuccess(c *gin.Context) {
	SuccessWithMessage(c, 200, constants.VALID_TOKEN)
}

// General Success Response Helpers - dapat digunakan untuk semua module
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// General - untuk response dengan message dan data
func General(c *gin.Context, code int, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
	}
	if data != nil {
		response.Data = data
	}

	c.JSON(code, gin.H{
		"status":   constants.SUCCESS,
		"code":     code,
		"response": response,
	})
}

// GeneralOK - untuk response 200 dengan message dan data
func GeneralOK(c *gin.Context, message string, data interface{}) {
	General(c, 200, message, data)
}

// GeneralCreated - untuk response 201 dengan message dan data
func GeneralCreated(c *gin.Context, message string, data interface{}) {
	General(c, 201, message, data)
}

// GeneralMessage - untuk response hanya dengan message
func GeneralMessage(c *gin.Context, code int, message string) {
	General(c, code, message, nil)
}

// GeneralData - untuk response hanya dengan data (message default)
func GeneralData(c *gin.Context, code int, data interface{}) {
	General(c, code, "Operation successful", data)
}

// General response helpers untuk operasi CRUD yang umum
func CreateSuccess(c *gin.Context, resourceName string, data interface{}) {
	GeneralCreated(c, fmt.Sprintf("%s created successfully", resourceName), data)
}

func UpdateSuccess(c *gin.Context, resourceName string, data interface{}) {
	GeneralOK(c, fmt.Sprintf("%s updated successfully", resourceName), data)
}

func DeleteSuccess(c *gin.Context, resourceName string) {
	GeneralOK(c, fmt.Sprintf("%s deleted successfully", resourceName), nil)
}

func FetchSuccess(c *gin.Context, resourceName string, data interface{}) {
	GeneralOK(c, fmt.Sprintf("%s retrieved successfully", resourceName), data)
}

func ListSuccess(c *gin.Context, resourceName string, data interface{}, total int64) {
	c.JSON(200, gin.H{
		"status": constants.SUCCESS,
		"code":   200,
		"response": gin.H{
			"message":   fmt.Sprintf("%s retrieved successfully", resourceName),
			"data":      data,
			"row_count": total,
		},
	})
}

func Error(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{
		"status": constants.ERROR,
		"code":   code,
		"data": gin.H{
			"message": message,
		},
	})
}

// ErrorFromAppError handles structured application errors
func ErrorFromAppError(c *gin.Context, err error) {
	if appErr, ok := appErrors.IsAppError(err); ok {
		c.JSON(appErr.Status, gin.H{
			"status": constants.ERROR,
			"code":   appErr.Status,
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	// Fallback for non-AppError types
	Error(c, 500, err.Error())
}

// ValidationError handles validation errors with multiple fields
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(400, gin.H{
		"status": constants.ERROR,
		"code":   400,
		"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "Validation failed",
			"details": errors,
		},
	})
}
