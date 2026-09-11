package handler

import (
	"net/http"
	"stardew-panel/middleware"
	"stardew-panel/models"
	"stardew-panel/service"
	"strings"

	"github.com/gin-gonic/gin"
)

var authService *service.AuthService

// InitAuthHandler 初始化认证处理器
func InitAuthHandler() *service.AuthService {
	authService = service.NewAuthService()
	return authService
}

// Login 登录
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	// 获取客户端信息
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	resp, err := authService.Login(req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": resp.Message})
		return
	}

	// 登录成功，重置限流计数
	middleware.ResetLoginAttempt(ipAddress)

	c.JSON(http.StatusOK, resp)
}

// Logout 退出登录
func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")
		authService.Logout(token)
	}

	c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码和新密码不能为空"})
		return
	}

	// 从上下文获取用户名（由中间件设置）
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	err := authService.ChangePassword(username.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功，请重新登录"})
}

// CheckAuth 检查认证状态
func CheckAuth(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"username":      username,
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过登录接口
		if c.Request.URL.Path == "/api/v1/auth/login" {
			c.Next()
			return
		}

		// 获取 Token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 验证 Token
		session, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户名存入上下文
		c.Set("username", session.Username)
		c.Next()
	}
}
