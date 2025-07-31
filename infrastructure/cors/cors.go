package cors

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCors() gin.HandlerFunc {
	// Get allowed origins from environment variable
	allowedOrigins := getAllowedOrigins()
	
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// getAllowedOrigins returns the list of allowed origins from environment variable
func getAllowedOrigins() []string {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	
	// Default origins for development
	defaultOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
	
	if allowedOriginsEnv == "" {
		return defaultOrigins
	}
	
	// Parse comma-separated origins from environment variable
	origins := strings.Split(allowedOriginsEnv, ",")
	var cleanOrigins []string
	
	for _, origin := range origins {
		cleanOrigin := strings.TrimSpace(origin)
		if cleanOrigin != "" {
			cleanOrigins = append(cleanOrigins, cleanOrigin)
		}
	}
	
	// If no valid origins found, return defaults
	if len(cleanOrigins) == 0 {
		return defaultOrigins
	}
	
	return cleanOrigins
}
