package jwt

import (
	"net/http"
	"os"

	"github.com/buildyow/byow-user-service/constants"
	"github.com/buildyow/byow-user-service/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Token From Cookie
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			response.Error(c, http.StatusUnauthorized, constants.ERR_INVALID_TOKEN)
			c.Abort()
			return
		}

		tokenStr := cookie.Value

		// Parse & Verification
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Method Sign
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, constants.ERR_INVALID_TOKEN)
			c.Abort()
			return
		}

		// Get Email Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims["user_id"].(string); ok {
				// Set ID to Context
				c.Set("user_id", userId)
			}
			if email, ok := claims["email"].(string); ok {
				// Set Email to Context
				c.Set("email", email)
			}
			if phone, ok := claims["phone"].(string); ok {
				// Set Phone to Context
				c.Set("phone", phone)
			}
		}

		c.Next()
	}
}
