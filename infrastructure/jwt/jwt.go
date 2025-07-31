package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user_id string, email string, phone string, secret string, minutes int) (string, error) {
	// Generate unique JTI (JWT ID) for token revocation
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user_id,
		"email":   email,
		"phone":   phone,
		"jti":     jti,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Minute * time.Duration(minutes)).Unix(),
		"iss":     "byow-user-service",
		"aud":     "byow-platform",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateJTI creates a unique JWT ID for token revocation
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
