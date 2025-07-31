package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"user_id"`
	CompanyName    string             `bson:"company_name"`
	CompanyEmail   string             `bson:"company_email"`
	CompanyPhone   string             `bson:"company_phone"`
	CompanyAddress string             `bson:"company_address"`
	CompanyLogo    string             `bson:"company_logo"`
	Verified       bool               `bson:"verified"`
	CreatedAt      time.Time          `bson:"created_at"`
}
