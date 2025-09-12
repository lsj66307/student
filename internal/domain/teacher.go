package domain

import (
	"time"
)

// Teacher 老师数据模型
type Teacher struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Age        int       `json:"age" db:"age"`
	Gender     string    `json:"gender" db:"gender"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	Subject    string    `json:"subject" db:"subject"`       // 教授科目
	Title      string    `json:"title" db:"title"`           // 职称
	Department string    `json:"department" db:"department"` // 所属院系
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTeacherRequest 创建老师请求结构
type CreateTeacherRequest struct {
	Name       string `json:"name" binding:"required"`
	Age        int    `json:"age" binding:"required,min=22,max=70"`
	Gender     string `json:"gender" binding:"required,oneof=男 女"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"required"`
	Subject    string `json:"subject" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Department string `json:"department" binding:"required"`
}

// UpdateTeacherRequest 更新老师请求结构
type UpdateTeacherRequest struct {
	Name       string `json:"name"`
	Age        int    `json:"age" binding:"min=22,max=70"`
	Gender     string `json:"gender" binding:"oneof=男 女"`
	Email      string `json:"email" binding:"email"`
	Phone      string `json:"phone"`
	Subject    string `json:"subject"`
	Title      string `json:"title"`
	Department string `json:"department"`
}
