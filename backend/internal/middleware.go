package internal

import (
	"net/http"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware limits requests per IP per minute
func RateLimitMiddleware(maxPerMin int) gin.HandlerFunc {
	var mu sync.Mutex
	var visitors = make(map[string][]time.Time)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		times := visitors[ip]
		// Remove old timestamps
		var fresh []time.Time
		for _, t := range times {
			if now.Sub(t) < time.Minute {
				fresh = append(fresh, t)
			}
		}
		if len(fresh) >= maxPerMin {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		fresh = append(fresh, now)
		visitors[ip] = fresh
		mu.Unlock()
		c.Next()
	}
}
