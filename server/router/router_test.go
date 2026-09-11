package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareAllowsConfiguredOriginOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(corsMiddleware("https://panel.example"))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	allowed := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://panel.example")
	r.ServeHTTP(allowed, req)
	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://panel.example" {
		t.Fatalf("configured origin not allowed: %q", got)
	}

	blocked := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://evil.example")
	r.ServeHTTP(blocked, req)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("unexpected status for unconfigured origin: %d", blocked.Code)
	}
}
