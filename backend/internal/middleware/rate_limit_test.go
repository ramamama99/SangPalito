package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSessionRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewSessionRateLimiter(3, 200*time.Millisecond)

	r := gin.New()
	r.Use(ErrorHandler())
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Make 3 requests — all should succeed (200 OK)
	for i := 1; i <= 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Session-ID", "test-session-123")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected status 200, got %d", i, w.Code)
		}
	}

	// 4th request must be rate limited (429 Too Many Requests)
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req4.Header.Set("X-Session-ID", "test-session-123")
	r.ServeHTTP(w4, req4)

	if w4.Code != http.StatusTooManyRequests {
		t.Fatalf("request 4: expected status 429, got %d", w4.Code)
	}

	retryAfter := w4.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Error("expected Retry-After header to be present")
	}

	// Wait for timeout to expire
	time.Sleep(250 * time.Millisecond)

	// After timeout, next request should succeed again
	w5 := httptest.NewRecorder()
	req5, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req5.Header.Set("X-Session-ID", "test-session-123")
	r.ServeHTTP(w5, req5)

	if w5.Code != http.StatusOK {
		t.Fatalf("request 5 after timeout: expected status 200, got %d", w5.Code)
	}
}
