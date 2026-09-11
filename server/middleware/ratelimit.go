package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器
type RateLimiter struct {
	attempts map[string]*LoginAttempt
	mu       sync.RWMutex
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	Count     int
	FirstTime time.Time
	BlockedAt time.Time
}

var rateLimiter = &RateLimiter{
	attempts: make(map[string]*LoginAttempt),
}
var cleanupOnce sync.Once

// LoginRateLimit 登录限流中间件
// 限制：同一 IP 5 分钟内最多尝试 5 次
// 超过后封禁 15 分钟
func LoginRateLimit() gin.HandlerFunc {
	// 定期清理过期记录
	cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			for range ticker.C {
				rateLimiter.cleanup()
			}
		}()
	})

	return func(c *gin.Context) {
		// 只对登录接口限流
		if c.Request.URL.Path != "/api/v1/auth/login" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		rateLimiter.mu.Lock()
		attempt, exists := rateLimiter.attempts[ip]
		now := time.Now()

		if !exists {
			// 首次尝试
			rateLimiter.attempts[ip] = &LoginAttempt{
				Count:     1,
				FirstTime: now,
			}
			rateLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否在封禁期
		if !attempt.BlockedAt.IsZero() {
			if now.Sub(attempt.BlockedAt) < 15*time.Minute {
				rateLimiter.mu.Unlock()
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "登录尝试过多，该 IP 已被临时锁定 15 分钟",
				})
				c.Abort()
				return
			}
			// 解除封禁
			attempt.BlockedAt = time.Time{}
			attempt.Count = 0
			attempt.FirstTime = now
		}

		// 检查是否超过 5 分钟，重置计数
		if now.Sub(attempt.FirstTime) > 5*time.Minute {
			attempt.Count = 1
			attempt.FirstTime = now
			rateLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 增加尝试次数
		attempt.Count++

		// 超过 5 次，封禁
		if attempt.Count > 5 {
			attempt.BlockedAt = now
			rateLimiter.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "登录尝试过多，该 IP 已被临时锁定 15 分钟",
			})
			c.Abort()
			return
		}

		rateLimiter.mu.Unlock()
		c.Next()
	}
}

// ResetLoginAttempt 登录成功后重置尝试次数
func ResetLoginAttempt(ip string) {
	rateLimiter.mu.Lock()
	defer rateLimiter.mu.Unlock()
	delete(rateLimiter.attempts, ip)
}

// cleanup 清理超过 1 小时的记录
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, attempt := range rl.attempts {
		// 清理 1 小时前的记录
		if now.Sub(attempt.FirstTime) > 1*time.Hour {
			delete(rl.attempts, ip)
		}
	}
}
