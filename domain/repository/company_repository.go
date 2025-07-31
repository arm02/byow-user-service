package repository

import (
	"github.com/buildyow/byow-user-service/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyRepository interface {
	FindAll(userID string, keyword string, limit int64, offset int64) ([]*entity.Company, int64, error)
	Create(user *entity.Company) error
	FindByID(id primitive.ObjectID) (*entity.Company, error)
	FindByEmail(email string) (*entity.Company, error)
	FindByPhone(phone string) (*entity.Company, error)
	Update(user *entity.Company) error
	Delete(id primitive.ObjectID) error
}
