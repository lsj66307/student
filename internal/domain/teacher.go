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
	SubjectID  int       `json:"subject_id" db:"subject_id" validate:"required,min=1"`
	Title      string    `json:"title" db:"title" validate:"required,min=2,max=30,nohtml,nosql"`           // 职称
	Department string    `json:"department" db:"department" validate:"required,min=2,max=50,nohtml,nosql"` // 所属院系
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// 关联数据
	Subject *Subject `json:"subject,omitempty" db:"-"`

	// 扩展字段（用于关联查询）
	SubjectName string `json:"subject_name,omitempty" db:"-"`
	SubjectCode string `json:"subject_code,omitempty" db:"-"`
}

// CreateTeacherRequest 创建老师请求结构
type CreateTeacherRequest struct {
	Name       string `json:"name" validate:"required,safename,nohtml,nosql"`
	Age        int    `json:"age" validate:"required,min=22,max=70"`
	Gender     string `json:"gender" validate:"required,oneof=男 女"`
	Email      string `json:"email" validate:"required,email,nohtml,nosql"`
	Phone      string `json:"phone" validate:"required,phone"`
	SubjectID  int    `json:"subject_id" validate:"required,min=1"`
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
	SubjectID  int    `json:"subject_id" validate:"omitempty,min=1"`
	Title      string `json:"title" validate:"omitempty,min=2,max=30,nohtml,nosql"`
	Department string `json:"department" validate:"omitempty,min=2,max=50,nohtml,nosql"`
}

// TeacherListRequest 教师列表请求结构
type TeacherListRequest struct {
	Page       int    `json:"page" form:"page" validate:"omitempty,min=1"`
	Size       int    `json:"size" form:"size" validate:"omitempty,min=1,max=100"`
	Name       string `json:"name" form:"name" validate:"omitempty,max=50,nohtml,nosql"`
	SubjectID  int    `json:"subject_id" form:"subject_id" validate:"omitempty,min=1"`
	Title      string `json:"title" form:"title" validate:"omitempty,max=30,nohtml,nosql"`
	Department string `json:"department" form:"department" validate:"omitempty,max=50,nohtml,nosql"`
}

// TeacherListResponse 教师列表响应结构
type TeacherListResponse struct {
	Teachers []Teacher `json:"teachers"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Size     int       `json:"size"`
}
