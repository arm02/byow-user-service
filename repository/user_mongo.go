package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/buildyow/byow-user-service/domain/entity"
	appErrors "github.com/buildyow/byow-user-service/domain/errors"
	"github.com/buildyow/byow-user-service/domain/repository"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepo struct {
	collection *mongo.Collection
	ch         *amqp091.Channel
}

func NewUserMongoRepo(db *mongo.Database, ch *amqp091.Channel) repository.UserRepository {
	return &userMongoRepo{
		collection: db.Collection("users_collections"),
		ch:         ch,
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
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userMongoRepo) FindByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"phone_number": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
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

func (r *userMongoRepo) PublishOTP(otpData interface{}) error {
	body, err := json.Marshal(otpData)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.ch.PublishWithContext(ctx,
		"",
		"otp_email_queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
