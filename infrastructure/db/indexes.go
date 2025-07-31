package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// CreateIndexes creates necessary database indexes for optimal performance
func CreateIndexes(db *mongo.Database, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create User indexes
	userCollection := db.Collection("users_collections")
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("email_unique"),
		},
		{
			Keys: bson.D{{Key: "phone_number", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("phone_unique"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().
				SetName("created_at_index"),
		},
		{
			Keys: bson.D{{Key: "is_verified", Value: 1}},
			Options: options.Index().
				SetName("is_verified_index"),
		},
		{
			Keys: bson.D{{Key: "is_onboarded", Value: 1}},
			Options: options.Index().
				SetName("is_onboarded_index"),
		},
		// Compound index for common queries
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "is_verified", Value: 1},
			},
			Options: options.Index().
				SetName("email_verified_compound"),
		},
	}

	// Create user indexes
	userIndexNames, err := userCollection.Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		logger.Error("Failed to create user indexes", zap.Error(err))
		return err
	}

	// Create Company indexes
	companyCollection := db.Collection("companies_collections")
	companyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
			Options: options.Index().
				SetName("company_name_index"),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetSparse(true).
				SetName("company_email_unique"),
		},
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().
				SetName("company_phone_index"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().
				SetName("company_created_at_index"),
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: 1}},
			Options: options.Index().
				SetName("company_updated_at_index"),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().
				SetName("company_user_id_index"),
		},
		// Compound index for user companies
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_companies_compound"),
		},
		// Text index for company search
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
			},
			Options: options.Index().
				SetName("company_search_text"),
		},
	}

	// Create company indexes
	companyIndexNames, err := companyCollection.Indexes().CreateMany(ctx, companyIndexes)
	if err != nil {
		logger.Error("Failed to create company indexes", zap.Error(err))
		return err
	}

	allIndexNames := append(userIndexNames, companyIndexNames...)
	logger.Info("Database indexes created successfully",
		zap.Strings("user_indexes", userIndexNames),
		zap.Strings("company_indexes", companyIndexNames),
		zap.Int("total_indexes", len(allIndexNames)))
	return nil
}

// DropIndexes drops all custom indexes (useful for testing or migration)
func DropIndexes(db *mongo.Database, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.Collection("users_collections")

	// List all indexes
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		return err
	}

	// Drop custom indexes (skip the default _id_ index)
	for _, index := range indexes {
		indexName := index["name"].(string)
		if indexName != "_id_" {
			_, err := collection.Indexes().DropOne(ctx, indexName)
			if err != nil {
				logger.Warn("Failed to drop index", zap.String("index", indexName), zap.Error(err))
			} else {
				logger.Info("Dropped index", zap.String("index", indexName))
			}
		}
	}

	return nil
}

// CheckIndexes verifies that all required indexes exist
func CheckIndexes(db *mongo.Database, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check user collection indexes
	userCollection := db.Collection("users_collections")
	userCursor, err := userCollection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer userCursor.Close(ctx)

	var userIndexes []bson.M
	if err = userCursor.All(ctx, &userIndexes); err != nil {
		return err
	}

	// Check company collection indexes
	companyCollection := db.Collection("companies_collections")
	companyCursor, err := companyCollection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer companyCursor.Close(ctx)

	var companyIndexes []bson.M
	if err = companyCursor.All(ctx, &companyIndexes); err != nil {
		return err
	}

	// Required user indexes
	requiredUserIndexes := []string{
		"email_unique",
		"phone_unique",
		"created_at_index",
		"is_verified_index",
		"is_onboarded_index",
		"email_verified_compound",
	}

	// Required company indexes
	requiredCompanyIndexes := []string{
		"company_name_index",
		"company_email_unique",
		"company_phone_index",
		"company_created_at_index",
		"company_updated_at_index",
		"company_user_id_index",
		"user_companies_compound",
		"company_search_text",
	}

	// Check user indexes
	existingUserIndexes := make(map[string]bool)
	for _, index := range userIndexes {
		if name, ok := index["name"].(string); ok {
			existingUserIndexes[name] = true
		}
	}

	// Check company indexes
	existingCompanyIndexes := make(map[string]bool)
	for _, index := range companyIndexes {
		if name, ok := index["name"].(string); ok {
			existingCompanyIndexes[name] = true
		}
	}

	// Find missing user indexes
	missingUserIndexes := []string{}
	for _, required := range requiredUserIndexes {
		if !existingUserIndexes[required] {
			missingUserIndexes = append(missingUserIndexes, required)
		}
	}

	// Find missing company indexes
	missingCompanyIndexes := []string{}
	for _, required := range requiredCompanyIndexes {
		if !existingCompanyIndexes[required] {
			missingCompanyIndexes = append(missingCompanyIndexes, required)
		}
	}

	// If any indexes are missing, recreate all
	if len(missingUserIndexes) > 0 || len(missingCompanyIndexes) > 0 {
		allMissing := append(missingUserIndexes, missingCompanyIndexes...)
		logger.Warn("Missing database indexes", zap.Strings("missing", allMissing))
		return CreateIndexes(db, logger)
	}

	logger.Info("All required database indexes are present")
	return nil
}

// RebuildCompanyIndexes rebuilds company indexes with proper sparse options
func RebuildCompanyIndexes(db *mongo.Database, logger *zap.Logger) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	companyCollection := db.Collection("companies_collections")

	// Drop existing company_email_unique index if it exists
	_, err := companyCollection.Indexes().DropOne(ctx, "company_email_unique")
	if err != nil {
		logger.Warn("Could not drop existing company_email_unique index", zap.Error(err))
	}

	// Create new sparse unique index for company email
	emailIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetSparse(true).
			SetName("company_email_unique"),
	}

	indexName, err := companyCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		logger.Error("Failed to create company email sparse index", zap.Error(err))
		return err
	}

	logger.Info("Company email index rebuilt successfully", zap.String("index", indexName))
	return nil
}
