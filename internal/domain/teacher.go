package domain

import (
	"time"
)

// Teacher 老师数据模型
type Teacher struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name" validate:"required,safename,nohtml,nosql"`
	Age        int       `json:"age" db:"age" validate:"required,min=22,max=70"`
	Gender     string    `json:"gender" db:"gender" validate:"required,oneof=男 女"`
	Email      string    `json:"email" db:"email" validate:"required,email,nohtml,nosql"`
	Phone      string    `json:"phone" db:"phone" validate:"required,phone"`
	Subject    string    `json:"subject" db:"subject" validate:"required,min=2,max=50,nohtml,nosql"`       // 教授科目
	Title      string    `json:"title" db:"title" validate:"required,min=2,max=30,nohtml,nosql"`           // 职称
	Department string    `json:"department" db:"department" validate:"required,min=2,max=50,nohtml,nosql"` // 所属院系
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTeacherRequest 创建老师请求结构
type CreateTeacherRequest struct {
	Name       string `json:"name" validate:"required,safename,nohtml,nosql"`
	Age        int    `json:"age" validate:"required,min=22,max=70"`
	Gender     string `json:"gender" validate:"required,oneof=男 女"`
	Email      string `json:"email" validate:"required,email,nohtml,nosql"`
	Phone      string `json:"phone" validate:"required,phone"`
	Subject    string `json:"subject" validate:"required,min=2,max=50,nohtml,nosql"`
	Title      string `json:"title" validate:"required,min=2,max=30,nohtml,nosql"`
	Department string `json:"department" validate:"required,min=2,max=50,nohtml,nosql"`
}

// UpdateTeacherRequest 更新老师请求结构
type UpdateTeacherRequest struct {
	Name       string `json:"name" validate:"omitempty,safename,nohtml,nosql"`
	Age        int    `json:"age" validate:"omitempty,min=22,max=70"`
	Gender     string `json:"gender" validate:"omitempty,oneof=男 女"`
	Email      string `json:"email" validate:"omitempty,email,nohtml,nosql"`
	Phone      string `json:"phone" validate:"omitempty,phone"`
	Subject    string `json:"subject" validate:"omitempty,min=2,max=50,nohtml,nosql"`
	Title      string `json:"title" validate:"omitempty,min=2,max=30,nohtml,nosql"`
	Department string `json:"department" validate:"omitempty,min=2,max=50,nohtml,nosql"`
}
