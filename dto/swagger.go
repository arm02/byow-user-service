package dto

import "mime/multipart"

type SwaggerRegisterRequest struct {
	Fullname    string                `form:"full_name" example:"John Doe"`
	Email       string                `form:"email" example:"john@example.com"`
	Password    string                `form:"password" example:"supersecret"`
	PhoneNumber string                `form:"phone_number" example:"628112123123"`
	Avatar      *multipart.FileHeader `form:"avatar" swaggerignore:"false"`
}
