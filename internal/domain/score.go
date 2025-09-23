package domain

import (
	"time"
)

// Score 成绩数据模型
type Score struct {
	ID        int       `json:"id" db:"id"`
	StudentID int       `json:"student_id" db:"student_id" validate:"required,min=1"`
	SubjectID int       `json:"subject_id" db:"subject_id" validate:"required,min=1"`
	TeacherID int       `json:"teacher_id" db:"teacher_id" validate:"required,min=1"`
	Score     float64   `json:"score" db:"score" validate:"required,min=0,max=100"`
	Semester  string    `json:"semester" db:"semester" validate:"required,min=5,max=20,nohtml,nosql"`
	ExamType  string    `json:"exam_type" db:"exam_type" validate:"required,oneof=midterm final quiz assignment"`
	Remarks   string    `json:"remarks" db:"remarks" validate:"omitempty,max=200,nohtml,nosql"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// 关联数据
	Student *Student `json:"student,omitempty" db:"-"`
	Subject *Subject `json:"subject,omitempty" db:"-"`
	Teacher *Teacher `json:"teacher,omitempty" db:"-"`

	// 扩展字段（用于关联查询）
	StudentName string `json:"student_name,omitempty" db:"-"`
	SubjectName string `json:"subject_name,omitempty" db:"-"`
	SubjectCode string `json:"subject_code,omitempty" db:"-"`
}

// CreateScoreRequest 创建成绩请求结构
type CreateScoreRequest struct {
	StudentID int     `json:"student_id" validate:"required,min=1"`
	SubjectID int     `json:"subject_id" validate:"required,min=1"`
	TeacherID int     `json:"teacher_id" validate:"required,min=1"`
	Score     float64 `json:"score" validate:"required,min=0,max=100"`
	Semester  string  `json:"semester" validate:"required,min=5,max=20,nohtml,nosql"`
	ExamType  string  `json:"exam_type" validate:"required,oneof=midterm final quiz assignment"`
	Remarks   string  `json:"remarks" validate:"omitempty,max=200,nohtml,nosql"`
}

// UpdateScoreRequest 更新成绩请求结构
type UpdateScoreRequest struct {
	Score    float64 `json:"score" validate:"omitempty,min=0,max=100"`
	Semester string  `json:"semester" validate:"omitempty,min=5,max=20,nohtml,nosql"`
	ExamType string  `json:"exam_type" validate:"omitempty,oneof=midterm final quiz assignment"`
	Remarks  string  `json:"remarks" validate:"omitempty,max=200,nohtml,nosql"`
}

// ScoreListRequest 成绩列表请求结构
type ScoreListRequest struct {
	Page      int     `json:"page" form:"page" validate:"omitempty,min=1"`
	Size      int     `json:"size" form:"size" validate:"omitempty,min=1,max=100"`
	StudentID int     `json:"student_id" form:"student_id" validate:"omitempty,min=1"`
	SubjectID int     `json:"subject_id" form:"subject_id" validate:"omitempty,min=1"`
	TeacherID int     `json:"teacher_id" form:"teacher_id" validate:"omitempty,min=1"`
	Semester  string  `json:"semester" form:"semester" validate:"omitempty,max=20,nohtml,nosql"`
	ExamType  string  `json:"exam_type" form:"exam_type" validate:"omitempty,oneof=midterm final quiz assignment"`
	MinScore  float64 `json:"min_score" form:"min_score" validate:"omitempty,min=0,max=100"`
	MaxScore  float64 `json:"max_score" form:"max_score" validate:"omitempty,min=0,max=100"`
}

// ScoreListResponse 成绩列表响应结构
type ScoreListResponse struct {
	Scores []Score `json:"scores"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
}

// BatchCreateScoresRequest 批量创建成绩请求结构
type BatchCreateScoresRequest struct {
	Scores []CreateScoreRequest `json:"scores" validate:"required,min=1,max=100,dive"`
}

// StudentScoreReport 学生成绩报告
type StudentScoreReport struct {
	StudentID    int                  `json:"student_id"`
	StudentName  string               `json:"student_name"`
	Semester     string               `json:"semester"`
	Scores       []SubjectScoreDetail `json:"scores"`
	TotalScore   float64              `json:"total_score"`
	AverageScore float64              `json:"average_score"`
	GPA          float64              `json:"gpa"`
}

// SubjectScoreDetail 科目成绩详情
type SubjectScoreDetail struct {
	SubjectID   int     `json:"subject_id"`
	SubjectName string  `json:"subject_name"`
	SubjectCode string  `json:"subject_code"`
	Credits     int     `json:"credits"`
	Score       float64 `json:"score"`
	ExamType    string  `json:"exam_type"`
}

// ClassScoreStatistics 班级成绩统计
type ClassScoreStatistics struct {
	SubjectID     int     `json:"subject_id"`
	SubjectName   string  `json:"subject_name"`
	Semester      string  `json:"semester"`
	ExamType      string  `json:"exam_type"`
	StudentCount  int     `json:"student_count"`
	AverageScore  float64 `json:"average_score"`
	MaxScore      float64 `json:"max_score"`
	MinScore      float64 `json:"min_score"`
	PassRate      float64 `json:"pass_rate"`      // 及格率
	ExcellentRate float64 `json:"excellent_rate"` // 优秀率 (>=90分)
}

// SubjectScoreStatistics 科目成绩统计
type SubjectScoreStatistics struct {
	SubjectID      int     `json:"subject_id"`
	TotalCount     int     `json:"total_count"`
	AverageScore   float64 `json:"average_score"`
	MinScore       float64 `json:"min_score"`
	MaxScore       float64 `json:"max_score"`
	ExcellentCount int     `json:"excellent_count"` // >=90分
	GoodCount      int     `json:"good_count"`      // 80-89分
	FairCount      int     `json:"fair_count"`      // 70-79分
	PassCount      int     `json:"pass_count"`      // 60-69分
	FailCount      int     `json:"fail_count"`      // <60分
}
