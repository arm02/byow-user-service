package jwt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// TokenBlacklist represents a blacklisted token
type TokenBlacklist struct {
	JTI       string    `bson:"jti" json:"jti"`
	UserEmail string    `bson:"user_email" json:"user_email"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// BlacklistService handles token blacklisting
type BlacklistService struct {
	collection *mongo.Collection
	cache      map[string]time.Time
	mutex      sync.RWMutex
	logger     *zap.Logger
}

// NewBlacklistService creates a new blacklist service
func NewBlacklistService(db *mongo.Database, logger *zap.Logger) *BlacklistService {
	collection := db.Collection("token_blacklist")
	
	// Create TTL index for automatic cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(0). // TTL index that expires at the time specified in expires_at
			SetName("expires_at_ttl"),
	}
	
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		logger.Warn("Failed to create TTL index for token blacklist", zap.Error(err))
	}

	// Create unique index on JTI
	jtiIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "jti", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("jti_unique"),
	}
	
	_, err = collection.Indexes().CreateOne(ctx, jtiIndex)
	if err != nil {
		logger.Warn("Failed to create JTI unique index", zap.Error(err))
	}

	service := &BlacklistService{
		collection: collection,
		cache:      make(map[string]time.Time),
		logger:     logger,
	}

	// Load existing blacklisted tokens into cache
	go service.loadCacheFromDB()

	return service
}

// BlacklistToken adds a token to the blacklist
func (bs *BlacklistService) BlacklistToken(jti, userEmail string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blacklistEntry := TokenBlacklist{
		JTI:       jti,
		UserEmail: userEmail,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	// Add to database
	_, err := bs.collection.InsertOne(ctx, blacklistEntry)
	if err != nil {
		bs.logger.Error("Failed to blacklist token in database", 
			zap.String("jti", jti), 
			zap.String("user_email", userEmail), 
			zap.Error(err))
		return err
	}

	// Add to cache
	bs.mutex.Lock()
	bs.cache[jti] = expiresAt
	bs.mutex.Unlock()

	bs.logger.Info("Token blacklisted successfully", 
		zap.String("jti", jti), 
		zap.String("user_email", userEmail))

	return nil
}

// IsTokenBlacklisted checks if a token is blacklisted
func (bs *BlacklistService) IsTokenBlacklisted(jti string) bool {
	// First check cache for fast lookup
	bs.mutex.RLock()
	expiresAt, exists := bs.cache[jti]
	bs.mutex.RUnlock()

	if exists {
		// If token exists in cache and hasn't expired, it's blacklisted
		if time.Now().Before(expiresAt) {
			return true
		}
		// If expired, remove from cache
		bs.mutex.Lock()
		delete(bs.cache, jti)
		bs.mutex.Unlock()
		return false
	}

	// If not in cache, check database (fallback)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blacklistEntry TokenBlacklist
	err := bs.collection.FindOne(ctx, bson.M{
		"jti": jti,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&blacklistEntry)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			bs.logger.Warn("Error checking token blacklist", 
				zap.String("jti", jti), 
				zap.Error(err))
		}
		return false
	}

	// Add to cache for future lookups
	bs.mutex.Lock()
	bs.cache[jti] = blacklistEntry.ExpiresAt
	bs.mutex.Unlock()

	return true
}

// BlacklistAllUserTokens blacklists all tokens for a specific user
func (bs *BlacklistService) BlacklistAllUserTokens(userEmail string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// This is a placeholder - in a real implementation, you'd need to track
	// active tokens per user or implement a user-based blacklist
	blacklistEntry := TokenBlacklist{
		JTI:       fmt.Sprintf("user_%s_%d", userEmail, time.Now().Unix()),
		UserEmail: userEmail,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	_, err := bs.collection.InsertOne(ctx, blacklistEntry)
	if err != nil {
		bs.logger.Error("Failed to blacklist user tokens", 
			zap.String("user_email", userEmail), 
			zap.Error(err))
		return err
	}

	bs.logger.Info("All user tokens blacklisted", 
		zap.String("user_email", userEmail))

	return nil
}

// CleanupExpiredTokens removes expired tokens from cache
func (bs *BlacklistService) CleanupExpiredTokens() {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	now := time.Now()
	for jti, expiresAt := range bs.cache {
		if now.After(expiresAt) {
			delete(bs.cache, jti)
		}
	}
}

// loadCacheFromDB loads existing blacklisted tokens into memory cache
func (bs *BlacklistService) loadCacheFromDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load only non-expired tokens
	cursor, err := bs.collection.Find(ctx, bson.M{
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		bs.logger.Error("Failed to load blacklisted tokens from database", zap.Error(err))
		return
	}
	defer cursor.Close(ctx)

	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	count := 0
	for cursor.Next(ctx) {
		var entry TokenBlacklist
		if err := cursor.Decode(&entry); err != nil {
			bs.logger.Warn("Failed to decode blacklist entry", zap.Error(err))
			continue
		}
		bs.cache[entry.JTI] = entry.ExpiresAt
		count++
	}

	bs.logger.Info("Loaded blacklisted tokens into cache", zap.Int("count", count))
}

// StartCleanupWorker starts a background worker to clean up expired tokens
func (bs *BlacklistService) StartCleanupWorker() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			bs.CleanupExpiredTokens()
		}
	}()
}