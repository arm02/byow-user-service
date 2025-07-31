package dto

type SuccessResponse struct {
	Status string      `json:"status" example:"SUCCESS"`
	Code   int         `json:"code" example:"200"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string            `json:"status" example:"ERROR"`
	Code   int               `json:"code" example:"400"`
	Data   ErrorResponseData `json:"data,omitempty"`
	Error  ErrorDetail       `json:"error,omitempty"`
}

type ErrorResponseData struct {
	Message string `json:"message" example:"INTERNAL_SERVER_ERROR"`
}

type ErrorDetail struct {
	Code    string      `json:"code" example:"VALIDATION_ERROR"`
	Message string      `json:"message" example:"Validation failed"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationErrorResponse struct {
	Status string               `json:"status" example:"ERROR"`
	Code   int                  `json:"code" example:"400"`
	Error  ValidationErrorDetail `json:"error"`
}

type ValidationErrorDetail struct {
	Code    string             `json:"code" example:"VALIDATION_ERROR"`
	Message string             `json:"message" example:"Validation failed"`
	Details []ValidationError  `json:"details"`
}

type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
}
