package logger

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogRequestBody(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil || c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		// Read body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("failed to read request body", zap.Error(err))
			c.Next()
			return
		}

		// Restore body to the request
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Optional: skip log endpoints that don't need request body logging
		skipPaths := map[string]bool{
			"/auth/users/login":           true,
			"/auth/users/change-password": true,
			"/auth/users/register":        true,
		}
		if !skipPaths[c.FullPath()] {
			logger.Info("Request Payload",
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
				zap.ByteString("body", bodyBytes),
			)
		}

		c.Next()
	}
}
