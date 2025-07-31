package repository

import (
	"context"
	"time"

	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/domain/entity"
	"github.com/buildyow/byow-user-service/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type companyMongoRepo struct {
	collection *mongo.Collection
}

func NewCompanyMongoRepo(db *mongo.Database) repository.CompanyRepository {
	return &companyMongoRepo{
		collection: db.Collection("companies_collections"),
	}
}

func (r *companyMongoRepo) FindAll(userID string, keyword string, limit int64, offset int64) ([]*entity.Company, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if keyword != "" {
		// case-insensitive dan partial match
		filter["company_name"] = bson.M{
			"$regex":   keyword,
			"$options": "i", // case-insensitive
		}
	}

	if userID != "" {
		filter["user_id"] = userID // exact match
	}
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)

	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var companies []*entity.Company
	for cursor.Next(ctx) {
		var company entity.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, 0, err
		}
		companies = append(companies, &company)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (r *companyMongoRepo) Create(company *entity.Company) error {
	// Build filter for duplicate check, only include non-empty fields
	orConditions := []bson.M{}
	
	if company.CompanyEmail != "" {
		orConditions = append(orConditions, bson.M{"company_email": company.CompanyEmail})
	}
	if company.CompanyPhone != "" {
		orConditions = append(orConditions, bson.M{"company_phone": company.CompanyPhone})
	}
	
	// Only check for duplicates if we have fields to check
	if len(orConditions) > 0 {
		filter := bson.M{"$or": orConditions}
		
		count, err := r.collection.CountDocuments(context.Background(), filter)
		if err != nil {
			return err
		}
		if count > 0 {
			return appErrors.ErrEmailOrPhoneAlreadyRegistered
		}
	}

	company.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(context.Background(), company)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		company.ID = oid
	}
	return nil
}

func (r *companyMongoRepo) FindByID(id primitive.ObjectID) (*entity.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	var company entity.Company
	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, appErrors.NewNotFoundError("Company")
		}
		return nil, err
	}

	return &company, nil
}

func (r *companyMongoRepo) FindByEmail(email string) (*entity.Company, error) {
	var company entity.Company
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&company)
	return &company, err
}

func (r *companyMongoRepo) FindByPhone(phone string) (*entity.Company, error) {
	var company entity.Company
	err := r.collection.FindOne(context.Background(), bson.M{"phone_number": phone}).Decode(&company)
	return &company, err
}

func (r *companyMongoRepo) Update(company *entity.Company) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"id": company.ID},
		bson.M{"$set": company},
	)

	return err
}

func (r *companyMongoRepo) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
