package middleware

import (
	"github.com/gin-gonic/gin"
)

// Security adds security headers to responses
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}

// RateLimitHeaders adds rate limiting headers (for future implementation)
func RateLimitHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// These would be set by an actual rate limiter
		c.Header("X-RateLimit-Limit", "1000")
		c.Header("X-RateLimit-Remaining", "999")
		c.Header("X-RateLimit-Reset", "3600")
		
		c.Next()
	}
}