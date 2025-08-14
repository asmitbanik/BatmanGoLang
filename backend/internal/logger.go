package internal

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLogger returns a zap.Logger (production config)
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// GinZapMiddleware returns a Gin middleware for logging requests
func GinZapMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		logger.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
