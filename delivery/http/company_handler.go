package http

import (
	"net/http"
	"strconv"
	"time"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/buildyow/byow-user-service/lib"
	"github.com/buildyow/byow-user-service/response"
	"github.com/buildyow/byow-user-service/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyHandler struct {
	Usecase *usecase.CompanyUsecase
}

func NewCompanyHandler(uc *usecase.CompanyUsecase) *CompanyHandler {
	return &CompanyHandler{Usecase: uc}
}

// @Summary Find All Companies
// @Tags Companies
// @Produce plain
// @Param keyword query string false "Keyword"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} dto.CompanyListResponseSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/companies/all [get]
func (h *CompanyHandler) FindAll(c *gin.Context) {
	keyword := c.Query("keyword")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var (
		limit  int64 = 10
		offset int64 = 0
	)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			offset = o
		}
	}

	companies, rowCount, err := h.Usecase.GetAll(c, keyword, limit, offset)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}

	response.ListSuccess(c, "Companies", companies, rowCount)
}

// @Summary Create Company
// @Description Register a new company
// @Tags Companies
// @Accept json
// @Produce json
// @Param company_name formData string true "Company Name" example(Cemerlang Jaya)
// @Param company_email formData string true "Company Email" example("john@company.com")
// @Param company_phone formData string true "Company Phone" example(628112123123)
// @Param company_address formData string true "Company Address" example("123 Cemerlang St, Tech City")
// @Param company_logo formData file false "Company Logo"
// @Success 201 {object} dto.CompanyRequestSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/companies/create [post]
func (h *CompanyHandler) Create(c *gin.Context) {
	var req dto.CompanyRequest
	// Bind form values to struct
	req.CompanyName = c.PostForm("company_name")
	req.CompanyEmail = c.PostForm("company_email")
	req.CompanyPhone = c.PostForm("company_phone")
	req.CompanyAddress = c.PostForm("company_address")

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		response.ErrorFromAppError(c, appErrors.ErrFailedParseMultipart)
		return
	}

	// Upload File
	file, _, err := c.Request.FormFile("company_logo")
	if err == nil {
		companyLogoUrl, err := lib.CloudinaryUpload(file)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		req.CompanyLogo = companyLogoUrl
	}

	// Call to usecase or saving to DB
	company, err := h.Usecase.Create(c, req)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}
	companyResponse := dto.CompanyResponse{
		CompanyID:      company.ID,
		CompanyName:    company.CompanyName,
		CompanyEmail:   company.CompanyEmail,
		CompanyPhone:   company.CompanyPhone,
		CompanyAddress: company.CompanyAddress,
		CompanyLogo:    company.CompanyLogo,
		UserID:         company.UserID,
		CreatedAt:      company.CreatedAt.Format(time.RFC3339),
	}
	response.CreateSuccess(c, "Company", companyResponse)
}

// @Summary Get Company By ID
// @Description Get company details by ID
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID" example("60d5ec49f1c2b14c88f3c5e5")
// @Success 200 {object} dto.CompanyRequestSwagger
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/companies/{id} [get]
func (h *CompanyHandler) FindByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		response.ErrorFromAppError(c, appErrors.ErrInvalidId)
		return
	}

	company, err := h.Usecase.FindByID(id)
	if err != nil {
		response.ErrorFromAppError(c, err)
		return
	}
	companyResponse := dto.CompanyResponse{
		CompanyID:      company.ID,
		CompanyName:    company.CompanyName,
		CompanyEmail:   company.CompanyEmail,
		CompanyPhone:   company.CompanyPhone,
		CompanyAddress: company.CompanyAddress,
		CompanyLogo:    company.CompanyLogo,
		UserID:         company.UserID,
		CreatedAt:      company.CreatedAt.Format(time.RFC3339),
	}
	response.FetchSuccess(c, "Company", companyResponse)
}
