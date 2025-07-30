package repository

import "github.com/buildyow/byow-user-service/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByPhone(phone string) (*entity.User, error)
	Update(user *entity.User) error
	UpdateEmail(user *entity.User, oldEmail string) error
	UpdatePhone(user *entity.User, oldPhone string) error
}
