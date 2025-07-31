package usecase

import (
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/domain/repository"
	"github.com/buildyow/byow-user-service/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyUsecase struct {
	Repo   repository.CompanyRepository
	UserID func(c *gin.Context) string
}

func (u *CompanyUsecase) GetAll(c *gin.Context, keyword string, limit int64, offset int64) (*[]dto.CompanyResponse, int64, error) {
	companies, rowCount, err := u.Repo.FindAll(u.UserID(c), keyword, limit, offset)
	if err != nil {
		return nil, 0, appErrors.NewNotFoundError("Companies")
	}

	var companyResponses []dto.CompanyResponse
	for _, company := range companies {
		companyResponses = append(companyResponses, dto.CompanyResponse{
			UserID:         company.UserID,
			CompanyID:      company.ID,
			CompanyName:    company.CompanyName,
			CompanyEmail:   company.CompanyEmail,
			CompanyPhone:   company.CompanyPhone,
			CompanyAddress: company.CompanyAddress,
			CompanyLogo:    company.CompanyLogo,
			Verified:       company.Verified,
			CreatedAt:      company.CreatedAt.Format(time.RFC3339),
		})
	}

	return &companyResponses, rowCount, nil
}

func (u *CompanyUsecase) Create(c *gin.Context, req dto.CompanyRequest) (*entity.Company, error) {
	company := &entity.Company{
		UserID:         u.UserID(c),
		CompanyName:    req.CompanyName,
		CompanyEmail:   req.CompanyEmail,
		CompanyPhone:   req.CompanyPhone,
		CompanyAddress: req.CompanyAddress,
		CompanyLogo:    req.CompanyLogo,
		Verified:       false,
	}
	err := u.Repo.Create(company)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (u *CompanyUsecase) FindByID(id primitive.ObjectID) (*entity.Company, error) {
	company, err := u.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return company, nil
}
