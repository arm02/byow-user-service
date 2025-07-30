package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, phone string, secret string, minutes int) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"phone": phone,
		"exp":   time.Now().Add(time.Minute * time.Duration(minutes)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
