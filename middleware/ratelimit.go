package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter protects against abuse
type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*RateLimiter)
	mutex    sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	visitor, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(time.Minute), 10) // 10 requests per minute
		visitors[ip] = &RateLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	visitor.lastSeen = time.Now()
	return visitor.limiter
}

// CleanupVisitors removes old entries to prevent memory leaks
func cleanupVisitors() {
	mutex.Lock()
	defer mutex.Unlock()

	for ip, visitor := range visitors {
		if time.Since(visitor.lastSeen) > 5*time.Minute {
			delete(visitors, ip)
		}
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	go func() {
		for {
			time.Sleep(time.Minute)
			cleanupVisitors()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
