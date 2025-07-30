package dto

type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Code   int         `json:"code" example:"200"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string            `json:"status" example:"error"`
	Code   int               `json:"code" example:"400"`
	Data   ErrorResponseData `json:"data"`
}

type ErrorResponseData struct {
	Message string `json:"message" example:"invalid request body"`
}
