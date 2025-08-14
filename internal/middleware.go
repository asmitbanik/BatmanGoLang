package internal

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "golang.org/x/time/rate"
)

// GinZap middleware logs requests using zap
func GinZap(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        logger.Info("request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("client_ip", c.ClientIP()),
        )
    }
}

// RateLimitMiddleware simple token bucket by IP
func RateLimitMiddleware(reqPerMin int) gin.HandlerFunc {
    // map of ip -> limiter
    var visitors = make(map[string]*rate.Limiter)

    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := getVisitorLimiter(visitors, ip, reqPerMin)
        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            return
        }
        c.Next()
    }
}

func getVisitorLimiter(visitors map[string]*rate.Limiter, ip string, reqPerMin int) *rate.Limiter {
    if limiter, exists := visitors[ip]; exists {
        return limiter
    }
    r := rate.Every(time.Minute / time.Duration(reqPerMin))
    limiter := rate.NewLimiter(r, 5)
    visitors[ip] = limiter
    return limiter
}
