package internal

import (
    "go.uber.org/zap"
)

func NewLogger(env string) *zap.Logger {
    if env == "production" {
        logger, _ := zap.NewProduction()
        return logger
    }
    logger, _ := zap.NewDevelopment()
    return logger
}

// GinZap is a tiny wrapper middleware to integrate zap with gin
func GinZap(logger *zap.Logger) gin.HandlerFunc {
    // Avoid importing gin in this file to keep separation; implement in middleware.go
    return nil
}
