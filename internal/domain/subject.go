package domain

import (
	"time"
)

// Subject 科目数据模型
type Subject struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=50,nohtml,nosql"`
	Code        string    `json:"code" db:"code" validate:"required,min=2,max=20,nohtml,nosql"`
	Description string    `json:"description" db:"description" validate:"omitempty,max=500,nohtml,nosql"`
	Credits     int       `json:"credits" db:"credits" validate:"required,min=1,max=10"`
	Status      string    `json:"status" db:"status" validate:"required,oneof=active inactive"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSubjectRequest 创建科目请求结构
type CreateSubjectRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50,nohtml,nosql"`
	Code        string `json:"code" validate:"required,min=2,max=20,nohtml,nosql"`
	Description string `json:"description" validate:"omitempty,max=500,nohtml,nosql"`
	Credits     int    `json:"credits" validate:"required,min=1,max=10"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
}

// UpdateSubjectRequest 更新科目请求结构
type UpdateSubjectRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=50,nohtml,nosql"`
	Code        string `json:"code" validate:"omitempty,min=2,max=20,nohtml,nosql"`
	Description string `json:"description" validate:"omitempty,max=500,nohtml,nosql"`
	Credits     int    `json:"credits" validate:"omitempty,min=1,max=10"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
}

// SubjectListRequest 科目列表请求结构
type SubjectListRequest struct {
	Page    int    `json:"page" form:"page" validate:"omitempty,min=1"`
	Size    int    `json:"size" form:"size" validate:"omitempty,min=1,max=100"`
	Name    string `json:"name" form:"name" validate:"omitempty,max=50,nohtml,nosql"`
	Code    string `json:"code" form:"code" validate:"omitempty,max=20,nohtml,nosql"`
	Status  string `json:"status" form:"status" validate:"omitempty,oneof=active inactive"`
	Credits int    `json:"credits" form:"credits" validate:"omitempty,min=1,max=10"`
}

// SubjectListResponse 科目列表响应结构
type SubjectListResponse struct {
	Subjects []Subject `json:"subjects"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Size     int       `json:"size"`
}
