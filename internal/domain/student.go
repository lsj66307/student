package domain

import (
	"time"
)

// Student 学生数据模型
type Student struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Age       int       `json:"age" db:"age"`
	Gender    string    `json:"gender" db:"gender"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Major     string    `json:"major" db:"major"`
	Grade     string    `json:"grade" db:"grade"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateStudentRequest 创建学生请求结构
type CreateStudentRequest struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required,min=1,max=150"`
	Gender string `json:"gender" binding:"required,oneof=男 女"`
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone" binding:"required"`
	Major  string `json:"major" binding:"required"`
	Grade  string `json:"grade" binding:"required"`
}

// UpdateStudentRequest 更新学生请求结构
type UpdateStudentRequest struct {
	Name   string `json:"name"`
	Age    int    `json:"age" binding:"min=1,max=150"`
	Gender string `json:"gender" binding:"oneof=男 女"`
	Email  string `json:"email" binding:"email"`
	Phone  string `json:"phone"`
	Major  string `json:"major"`
	Grade  string `json:"grade"`
}
