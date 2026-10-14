package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"sangpalito/backend/internal/apperror"
)

type sessionBucket struct {
	count        int
	windowStart  time.Time
	blockedUntil time.Time
}

// SessionRateLimiter limits requests to maxRequests per session, then blocks for the given timeout.
type SessionRateLimiter struct {
	maxRequests int
	timeout     time.Duration
	sessions    map[string]*sessionBucket
	mu          sync.Mutex
}

// NewSessionRateLimiter creates a new thread-safe session rate limiter.
func NewSessionRateLimiter(maxRequests int, timeout time.Duration) *SessionRateLimiter {
	limiter := &SessionRateLimiter{
		maxRequests: maxRequests,
		timeout:     timeout,
		sessions:    make(map[string]*sessionBucket),
	}

	// Background cleanup of stale session buckets every 5 minutes
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Prevent goroutine crash
			}
		}()

		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return limiter
}

func (l *SessionRateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, bucket := range l.sessions {
		// Remove entries older than 2x timeout and not currently blocked
		if now.After(bucket.blockedUntil) && now.Sub(bucket.windowStart) > (2*l.timeout) {
			delete(l.sessions, key)
		}
	}
}

// Middleware returns a Gin HandlerFunc enforcing the session rate limit.
func (l *SessionRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey := extractSessionKey(c)

		l.mu.Lock()
		now := time.Now()
		bucket, exists := l.sessions[sessionKey]

		if !exists {
			bucket = &sessionBucket{
				count:       1,
				windowStart: now,
			}
			l.sessions[sessionKey] = bucket
			l.mu.Unlock()
			c.Next()
			return
		}

		// Check if session is currently blocked in timeout
		if now.Before(bucket.blockedUntil) {
			remaining := int(bucket.blockedUntil.Sub(now).Seconds()) + 1
			l.mu.Unlock()

			c.Header("Retry-After", fmt.Sprintf("%d", remaining))
			_ = c.Error(apperror.New(
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				fmt.Sprintf("Batas request tercapai (maksimal %d kali per sesi). Silakan tunggu %d detik sebelum mencoba lagi.", l.maxRequests, remaining),
				nil,
			))
			c.Abort()
			return
		}

		// If timeout has passed, reset window
		if !bucket.blockedUntil.IsZero() && now.After(bucket.blockedUntil) {
			bucket.count = 1
			bucket.windowStart = now
			bucket.blockedUntil = time.Time{}
			l.mu.Unlock()
			c.Next()
			return
		}

		// Increment count
		bucket.count++

		// Check if limit exceeded
		if bucket.count > l.maxRequests {
			bucket.blockedUntil = now.Add(l.timeout)
			remaining := int(l.timeout.Seconds())
			l.mu.Unlock()

			c.Header("Retry-After", fmt.Sprintf("%d", remaining))
			_ = c.Error(apperror.New(
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				fmt.Sprintf("Batas request tercapai (maksimal %d kali per sesi). Silakan tunggu %d detik sebelum mencoba lagi.", l.maxRequests, remaining),
				nil,
			))
			c.Abort()
			return
		}

		l.mu.Unlock()
		c.Next()
	}
}

func extractSessionKey(c *gin.Context) string {
	// 1. Authenticated user ID in context
	if userID, exists := c.Get(CtxUserID); exists {
		if idStr, ok := userID.(string); ok && idStr != "" {
			return "user:" + idStr
		}
	}

	// 2. Authorization Bearer token header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" && parts[1] != "" {
			return "token:" + parts[1]
		}
	}

	// 3. Custom session header
	if sessHeader := c.GetHeader("X-Session-ID"); sessHeader != "" {
		return "sess:" + sessHeader
	}

	// 4. Fallback to Client IP
	return "ip:" + c.ClientIP()
}
