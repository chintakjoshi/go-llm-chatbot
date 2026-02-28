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
	visitors    = make(map[string]*RateLimiter)
	mutex       sync.Mutex
	cleanerOnce sync.Once
)

func getVisitor(ip string, maxRequests int, windowSeconds int) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	visitor, exists := visitors[ip]
	if !exists {
		// rateLimit per windowSeconds, with a burst equal to one full window.
		requestsPerSecond := rate.Limit(float64(maxRequests) / float64(windowSeconds))
		limiter := rate.NewLimiter(requestsPerSecond, maxRequests)
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

func RateLimitMiddleware(maxRequests int, windowSeconds int) gin.HandlerFunc {
	if maxRequests <= 0 {
		maxRequests = 10
	}
	if windowSeconds <= 0 {
		windowSeconds = 60
	}

	cleanerOnce.Do(func() {
		go func() {
			for {
				time.Sleep(time.Minute)
				cleanupVisitors()
			}
		}()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip, maxRequests, windowSeconds)

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
