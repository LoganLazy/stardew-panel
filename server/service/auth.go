package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"stardew-panel/database"
	"stardew-panel/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 用户登录
func (s *AuthService) Login(username, password, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// 查询用户
	user, err := s.getUserByUsername(username)
	if err != nil {
		return &models.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, nil
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return &models.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, nil
	}

	// 生成 Token
	token, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	// 保存 Token 到数据库
	expiresAt := time.Now().Add(24 * time.Hour)
	query := `INSERT INTO sessions (token, username, expires_at, ip_address, user_agent)
	          VALUES (?, ?, ?, ?, ?)`
	_, err = database.DB.Exec(query, token, username, expiresAt, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("保存会话失败: %w", err)
	}

	return &models.LoginResponse{
		Success: true,
		Token:   token,
		Message: "登录成功",
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(username, oldPassword, newPassword string) error {
	// 查询用户
	user, err := s.getUserByUsername(username)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("旧密码错误")
	}

	// 密码强度检查
	if len(newPassword) < 6 {
		return fmt.Errorf("新密码至少需要 6 个字符")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	query := `UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE username = ?`
	_, err = database.DB.Exec(query, string(hashedPassword), username)
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// InitDefaultUser 初始化默认用户（首次启动时）
func (s *AuthService) InitDefaultUser() error {
	// 检查是否已有用户
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return fmt.Errorf("查询用户数失败: %w", err)
	}

	// 如果已有用户，跳过
	if count > 0 {
		return nil
	}

	// 创建默认用户：admin / admin123
	defaultPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	query := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	_, err = database.DB.Exec(query, "admin", string(hashedPassword))
	if err != nil {
		return fmt.Errorf("创建默认用户失败: %w", err)
	}

	fmt.Println("✅ 默认用户已创建")
	fmt.Println("   用户名: admin")
	fmt.Println("   密码: admin123")
	fmt.Println("   ⚠️  请立即登录后修改密码！")

	return nil
}

// ValidateToken 验证 Token
func (s *AuthService) ValidateToken(token string) (*Session, error) {
	// 从数据库查询会话
	session := &Session{}
	query := `SELECT token, username, created_at, expires_at FROM sessions WHERE token = ?`
	err := database.DB.QueryRow(query, token).Scan(
		&session.Token,
		&session.Username,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("无效的 Token")
	}

	// 检查是否过期
	if time.Now().After(session.ExpiresAt) {
		// 删除过期会话
		database.DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
		return nil, fmt.Errorf("Token 已过期，请重新登录")
	}

	// Token 快过期时自动刷新（< 1小时）
	if time.Until(session.ExpiresAt) < 1*time.Hour {
		newExpires := time.Now().Add(24 * time.Hour)
		database.DB.Exec(`UPDATE sessions SET expires_at = ? WHERE token = ?`, newExpires, token)
		session.ExpiresAt = newExpires
	}

	return session, nil
}

// Logout 退出登录
func (s *AuthService) Logout(token string) {
	database.DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
}

// getUserByUsername 根据用户名查询用户
func (s *AuthService) getUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?`
	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// generateToken 生成随机 Token
func (s *AuthService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Session 会话信息
type Session struct {
	Token     string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// CleanupExpiredSessions 清理过期会话（定时任务）
func (s *AuthService) CleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			database.DB.Exec(`DELETE FROM sessions WHERE expires_at < datetime('now')`)
		}
	}()
}
