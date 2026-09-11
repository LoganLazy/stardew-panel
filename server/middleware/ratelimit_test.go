package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginRateLimitAndReset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rateLimiter.mu.Lock()
	rateLimiter.attempts = make(map[string]*LoginAttempt)
	rateLimiter.mu.Unlock()

	r := gin.New()
	r.Use(LoginRateLimit())
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})

	request := func() int {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		r.ServeHTTP(w, req)
		return w.Code
	}
	for i := 0; i < 5; i++ {
		if status := request(); status != http.StatusUnauthorized {
			t.Fatalf("attempt %d returned %d", i+1, status)
		}
	}
	if status := request(); status != http.StatusTooManyRequests {
		t.Fatalf("sixth attempt returned %d", status)
	}

	ResetLoginAttempt("192.0.2.1")
	if status := request(); status != http.StatusUnauthorized {
		t.Fatalf("attempt after reset returned %d", status)
	}
}
