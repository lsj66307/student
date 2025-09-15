package domain

import (
	"time"
)

// Grade 成绩模型 - 简化为5个科目
type Grade struct {
	ID               int       `json:"id" db:"id"`
	StudentID        int       `json:"student_id" db:"student_id" validate:"required,min=1"`                            // 学生ID
	ChineseScore     *float64  `json:"chinese_score,omitempty" db:"chinese_score" validate:"omitempty,min=0,max=100"`   // 语文成绩
	MathScore        *float64  `json:"math_score,omitempty" db:"math_score" validate:"omitempty,min=0,max=100"`         // 数学成绩
	EnglishScore     *float64  `json:"english_score,omitempty" db:"english_score" validate:"omitempty,min=0,max=100"`   // 英语成绩
	SportsScore      *float64  `json:"sports_score,omitempty" db:"sports_score" validate:"omitempty,min=0,max=100"`     // 体育成绩
	MusicScore       *float64  `json:"music_score,omitempty" db:"music_score" validate:"omitempty,min=0,max=100"`       // 音乐成绩
	ChineseTeacherID *int      `json:"chinese_teacher_id,omitempty" db:"chinese_teacher_id" validate:"omitempty,min=1"` // 语文老师ID
	MathTeacherID    *int      `json:"math_teacher_id,omitempty" db:"math_teacher_id" validate:"omitempty,min=1"`       // 数学老师ID
	EnglishTeacherID *int      `json:"english_teacher_id,omitempty" db:"english_teacher_id" validate:"omitempty,min=1"` // 英语老师ID
	SportsTeacherID  *int      `json:"sports_teacher_id,omitempty" db:"sports_teacher_id" validate:"omitempty,min=1"`   // 体育老师ID
	MusicTeacherID   *int      `json:"music_teacher_id,omitempty" db:"music_teacher_id" validate:"omitempty,min=1"`     // 音乐老师ID
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// GradeWithDetails 成绩详情（包含学生和老师信息）
type GradeWithDetails struct {
	Grade
	StudentName        string `json:"student_name"`
	SportsTeacherName  string `json:"sports_teacher_name,omitempty"`
	MusicTeacherName   string `json:"music_teacher_name,omitempty"`
	ChineseTeacherName string `json:"chinese_teacher_name,omitempty"`
	MathTeacherName    string `json:"math_teacher_name,omitempty"`
	EnglishTeacherName string `json:"english_teacher_name,omitempty"`
}

// CreateGradeRequest 创建成绩请求
type CreateGradeRequest struct {
	StudentID        int      `json:"student_id" validate:"required,min=1" example:"1"`
	ChineseScore     *float64 `json:"chinese_score,omitempty" validate:"omitempty,min=0,max=100" example:"85.5"`
	MathScore        *float64 `json:"math_score,omitempty" validate:"omitempty,min=0,max=100" example:"90.0"`
	EnglishScore     *float64 `json:"english_score,omitempty" validate:"omitempty,min=0,max=100" example:"88.0"`
	SportsScore      *float64 `json:"sports_score,omitempty" validate:"omitempty,min=0,max=100" example:"92.0"`
	MusicScore       *float64 `json:"music_score,omitempty" validate:"omitempty,min=0,max=100" example:"87.0"`
	ChineseTeacherID *int     `json:"chinese_teacher_id,omitempty" validate:"omitempty,min=1" example:"1"`
	MathTeacherID    *int     `json:"math_teacher_id,omitempty" validate:"omitempty,min=1" example:"2"`
	EnglishTeacherID *int     `json:"english_teacher_id,omitempty" validate:"omitempty,min=1" example:"3"`
	SportsTeacherID  *int     `json:"sports_teacher_id,omitempty" validate:"omitempty,min=1" example:"4"`
	MusicTeacherID   *int     `json:"music_teacher_id,omitempty" validate:"omitempty,min=1" example:"5"`
}

// UpdateGradeRequest 更新成绩请求
type UpdateGradeRequest struct {
	ChineseScore     *float64 `json:"chinese_score,omitempty" validate:"omitempty,min=0,max=100" example:"85.5"`
	MathScore        *float64 `json:"math_score,omitempty" validate:"omitempty,min=0,max=100" example:"90.0"`
	EnglishScore     *float64 `json:"english_score,omitempty" validate:"omitempty,min=0,max=100" example:"88.0"`
	SportsScore      *float64 `json:"sports_score,omitempty" validate:"omitempty,min=0,max=100" example:"92.0"`
	MusicScore       *float64 `json:"music_score,omitempty" validate:"omitempty,min=0,max=100" example:"87.0"`
	ChineseTeacherID *int     `json:"chinese_teacher_id,omitempty" validate:"omitempty,min=1" example:"1"`
	MathTeacherID    *int     `json:"math_teacher_id,omitempty" validate:"omitempty,min=1" example:"2"`
	EnglishTeacherID *int     `json:"english_teacher_id,omitempty" validate:"omitempty,min=1" example:"3"`
	SportsTeacherID  *int     `json:"sports_teacher_id,omitempty" validate:"omitempty,min=1" example:"4"`
	MusicTeacherID   *int     `json:"music_teacher_id,omitempty" validate:"omitempty,min=1" example:"5"`
}

// GradeQueryParams 成绩查询参数
type GradeQueryParams struct {
	StudentID *int `form:"student_id" validate:"omitempty,min=1" example:"1"`
	Page      int  `form:"page" validate:"omitempty,min=1" example:"1"`
	Size      int  `form:"size" validate:"omitempty,min=1,max=100" example:"10"`
}
