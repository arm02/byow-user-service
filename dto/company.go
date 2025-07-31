package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CompanyResponse struct {
	UserID         string             `json:"user_id" example:"60c72b2f9b1e8c001c8e4d3a"`
	CompanyID      primitive.ObjectID `json:"company_id" example:"60c72b2f9b1e8c001c8e4d3a"`
	CompanyName    string             `json:"company_name" example:"BuildYow"`
	CompanyEmail   string             `json:"company_email" example:"info@buildyow.com"`
	CompanyPhone   string             `json:"company_phone" example:"628112123123"`
	CompanyAddress string             `json:"company_address" example:"123 BuildYow St, Tech City"`
	CompanyLogo    string             `json:"company_logo" example:"https://assets/images/company_logo.jpg"`
	Verified       bool               `json:"verified" example:"false"`
	CreatedAt      string             `json:"created_at" example:"2023-10-01T12:00:00Z"`
}

type CompanyListResponseSwagger struct {
	Status string            `json:"status" example:"SUCCESS"`
	Code   int               `json:"code" example:"200"`
	Data   []CompanyResponse `json:"data"`
}

type CompanyRequest struct {
	CompanyName    string `json:"company_name" example:"BuildYow"`
	CompanyEmail   string `json:"company_email" example:"info@buildyow.com"`
	CompanyPhone   string `json:"company_phone" example:"628112123123"`
	CompanyAddress string `json:"company_address" example:"123 BuildYow St, Tech City"`
	CompanyLogo    string `json:"company_logo" example:"https://assets/images/company_logo.jpg"`
	Verified       bool   `json:"verified" example:"false"`
}

type CompanyRequestSwagger struct {
	Status string          `json:"status" example:"SUCCESS"`
	Code   int             `json:"code" example:"200"`
	Data   CompanyResponse `json:"data"`
}
