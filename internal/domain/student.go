package domain

import (
	"time"
)

// Student 学生数据模型
type Student struct {
	ID             int        `json:"id" db:"id"`
	StudentID      string     `json:"student_id" db:"student_id" validate:"required,studentid"`
	Name           string     `json:"name" db:"name" validate:"required,safename,nohtml,nosql"`
	Age            int        `json:"age" db:"age" validate:"required,min=16,max=60"`
	Gender         string     `json:"gender" db:"gender" validate:"required,oneof=男 女"`
	Phone          string     `json:"phone" db:"phone" validate:"required,phone"`
	Email          string     `json:"email" db:"email" validate:"required,email,nohtml,nosql"`
	Address        string     `json:"address" db:"address" validate:"omitempty,max=200,nohtml,nosql"`
	Major          string     `json:"major" db:"major" validate:"required,min=2,max=50,nohtml,nosql"`
	EnrollmentDate *time.Time `json:"enrollment_date" db:"enrollment_date"`
	GraduationDate *time.Time `json:"graduation_date" db:"graduation_date"`
	Status         string     `json:"status" db:"status" validate:"required,oneof=active inactive graduated"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateStudentRequest 创建学生请求结构
type CreateStudentRequest struct {
	StudentID      string     `json:"student_id" validate:"required,studentid"`
	Name           string     `json:"name" validate:"required,safename,nohtml,nosql"`
	Age            int        `json:"age" validate:"required,min=16,max=60"`
	Gender         string     `json:"gender" validate:"required,oneof=男 女"`
	Phone          string     `json:"phone" validate:"required,phone"`
	Email          string     `json:"email" validate:"required,email,nohtml,nosql"`
	Address        string     `json:"address" validate:"omitempty,max=200,nohtml,nosql"`
	Major          string     `json:"major" validate:"required,min=2,max=50,nohtml,nosql"`
	EnrollmentDate *time.Time `json:"enrollment_date"`
	GraduationDate *time.Time `json:"graduation_date"`
	Status         string     `json:"status" validate:"omitempty,oneof=active inactive graduated"`
}

// UpdateStudentRequest 更新学生请求结构
type UpdateStudentRequest struct {
	StudentID      string     `json:"student_id" validate:"omitempty,studentid"`
	Name           string     `json:"name" validate:"omitempty,safename,nohtml,nosql"`
	Age            int        `json:"age" validate:"omitempty,min=16,max=60"`
	Gender         string     `json:"gender" validate:"omitempty,oneof=男 女"`
	Phone          string     `json:"phone" validate:"omitempty,phone"`
	Email          string     `json:"email" validate:"omitempty,email,nohtml,nosql"`
	Address        string     `json:"address" validate:"omitempty,max=200,nohtml,nosql"`
	Major          string     `json:"major" validate:"omitempty,min=2,max=50,nohtml,nosql"`
	EnrollmentDate *time.Time `json:"enrollment_date"`
	GraduationDate *time.Time `json:"graduation_date"`
	Status         string     `json:"status" validate:"omitempty,oneof=active inactive graduated"`
}

// BatchCreateStudentsRequest 批量创建学生请求结构
type BatchCreateStudentsRequest struct {
	Students []CreateStudentRequest `json:"students" validate:"required,min=1,max=100,dive"`
}

// BatchDeleteStudentsRequest 批量删除学生请求结构
type BatchDeleteStudentsRequest struct {
	IDs []int `json:"ids" validate:"required,min=1,max=100,dive,min=1"`
}

// TransferStudentMajorRequest 转专业请求结构
type TransferStudentMajorRequest struct {
	NewMajor string `json:"new_major" validate:"required,min=2,max=50,nohtml,nosql"`
	Reason   string `json:"reason" validate:"required,min=10,max=500,nohtml,nosql"`
}
