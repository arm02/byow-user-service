package dto

type SuccessResponse struct {
	Status string      `json:"status" example:"SUCCESS"`
	Code   int         `json:"code" example:"200"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status string            `json:"status" example:"ERROR"`
	Code   int               `json:"code" example:"400"`
	Data   ErrorResponseData `json:"data"`
}

type ErrorResponseData struct {
	Message string `json:"message" example:"INTERNAL_SERVER_ERROR"`
}
