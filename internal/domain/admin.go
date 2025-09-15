package domain

import (
	"errors"
	"time"
)

// 认证相关错误
var (
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// Admin 管理员模型
type Admin struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"` // 不在JSON中显示密码
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	Admin     AdminInfo `json:"admin"`
}

// AdminInfo 管理员信息（不包含密码）
type AdminInfo struct {
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
}

// JWTClaims JWT声明结构体
type JWTClaims struct {
	AdminID  int    `json:"admin_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

// Valid 验证JWT声明是否有效
func (c *JWTClaims) Valid() error {
	if time.Now().Unix() > c.Exp {
		return ErrTokenExpired
	}
	return nil
}
