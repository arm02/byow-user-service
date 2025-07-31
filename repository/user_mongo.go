package repository

import (
	"context"
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	"github.com/buildyow/byow-user-service/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepo struct {
	collection *mongo.Collection
}

func NewUserMongoRepo(db *mongo.Database) repository.UserRepository {
	return &userMongoRepo{
		collection: db.Collection("user_collections"),
	}
}

func (r *userMongoRepo) Create(user *entity.User) error {
	user.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *userMongoRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *userMongoRepo) FindByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"phone_number": phone}).Decode(&user)
	return &user, err
}

func (r *userMongoRepo) Update(user *entity.User) error {
	updateData, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	var updateMap bson.M
	err = bson.Unmarshal(updateData, &updateMap)
	if err != nil {
		return err
	}

	delete(updateMap, "_id")

	unsetMap := bson.M{}
	if user.OTP == "" {
		unsetMap["otp"] = ""
		unsetMap["otp_expires_at"] = ""
		unsetMap["otp_type"] = ""
	}

	update := bson.M{}
	if len(updateMap) > 0 {
		update["$set"] = updateMap
	}
	if len(unsetMap) > 0 {
		update["$unset"] = unsetMap
	}
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": user.Email},
		update,
	)

	return err
}

func (r *userMongoRepo) UpdateEmail(user *entity.User, oldEmail string) error {
	updateData, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	var updateMap bson.M
	err = bson.Unmarshal(updateData, &updateMap)
	if err != nil {
		return err
	}

	delete(updateMap, "_id")

	unsetMap := bson.M{}
	if user.OTP == "" {
		unsetMap["otp"] = ""
		unsetMap["otp_expires_at"] = ""
		unsetMap["otp_type"] = ""
	}

	update := bson.M{}
	if len(updateMap) > 0 {
		update["$set"] = updateMap
	}
	if len(unsetMap) > 0 {
		update["$unset"] = unsetMap
	}
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"email": oldEmail},
		update,
	)

	return err
}

func (r *userMongoRepo) UpdatePhone(user *entity.User, oldPhone string) error {
	updateData, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	var updateMap bson.M
	err = bson.Unmarshal(updateData, &updateMap)
	if err != nil {
		return err
	}

	delete(updateMap, "_id")

	unsetMap := bson.M{}
	if user.OTP == "" {
		unsetMap["otp"] = ""
		unsetMap["otp_expires_at"] = ""
		unsetMap["otp_type"] = ""
	}

	update := bson.M{}
	if len(updateMap) > 0 {
		update["$set"] = updateMap
	}
	if len(unsetMap) > 0 {
		update["$unset"] = unsetMap
	}
	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"phone_number": oldPhone},
		update,
	)

	return err
}
