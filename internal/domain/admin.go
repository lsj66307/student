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
	Account   string    `json:"account" db:"account" validate:"required,min=3,max=50,nohtml,nosql"` // 账号
	Password  string    `json:"-" db:"password" validate:"required,min=6,max=100"`                  // 密码，不在JSON中显示
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=50,nohtml,nosql"`       // 用户姓名
	Phone     string    `json:"phone" db:"phone" validate:"omitempty,len=11,numeric"`               // 手机号
	Email     string    `json:"email" db:"email" validate:"omitempty,email,max=100"`                // 邮箱
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Account  string `json:"account" validate:"required,min=3,max=50,nohtml,nosql" example:"admin"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"123456"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`
	Admin     AdminInfo `json:"admin"`
}

// AdminInfo 管理员信息（不包含密码）
type AdminInfo struct {
	ID      int    `json:"id" example:"1"`
	Account string `json:"account" example:"admin"`
	Name    string `json:"name" example:"管理员"`
	Phone   string `json:"phone" example:"13800138000"`
	Email   string `json:"email" example:"admin@example.com"`
}

// JWTClaims JWT声明结构体
type JWTClaims struct {
	AdminID int    `json:"admin_id"`
	Account string `json:"account"`
	Exp     int64  `json:"exp"`
	Iat     int64  `json:"iat"`
}

// Valid 验证JWT声明是否有效
func (c *JWTClaims) Valid() error {
	if time.Now().Unix() > c.Exp {
		return ErrTokenExpired
	}
	return nil
}

// CreateAdminRequest 创建管理员请求结构体
type CreateAdminRequest struct {
	Account  string `json:"account" validate:"required,min=3,max=50,nohtml,nosql" example:"admin001"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"123456"`
	Name     string `json:"name" validate:"required,min=1,max=50,nohtml,nosql" example:"张三"`
	Phone    string `json:"phone" validate:"omitempty,len=11,numeric" example:"13800138000"`
	Email    string `json:"email" validate:"omitempty,email,max=100" example:"admin@example.com"`
}

// UpdateAdminRequest 更新管理员请求结构体
type UpdateAdminRequest struct {
	Account  string `json:"account" validate:"required,min=3,max=50,nohtml,nosql" example:"admin001"`
	Password string `json:"password" validate:"omitempty,min=6,max=100" example:"123456"` // 密码可选，不填则不更新
	Name     string `json:"name" validate:"required,min=1,max=50,nohtml,nosql" example:"张三"`
	Phone    string `json:"phone" validate:"omitempty,len=11,numeric" example:"13800138000"`
	Email    string `json:"email" validate:"omitempty,email,max=100" example:"admin@example.com"`
}

// AdminListRequest 管理员列表请求结构体
type AdminListRequest struct {
	Page     int    `json:"page" form:"page" validate:"omitempty,min=1" example:"1"`
	PageSize int    `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100" example:"10"`
	Account  string `json:"account" form:"account" validate:"omitempty,max=50" example:"admin"`
	Name     string `json:"name" form:"name" validate:"omitempty,max=50" example:"张三"`
}

// AdminListResponse 管理员列表响应结构体
type AdminListResponse struct {
	Total int         `json:"total" example:"100"`
	Page  int         `json:"page" example:"1"`
	Size  int         `json:"size" example:"10"`
	Data  []AdminInfo `json:"data"`
}
