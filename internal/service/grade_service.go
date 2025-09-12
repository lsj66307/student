package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"student-management-system/internal/domain"
)

// GradeService 成绩服务
type GradeService struct {
	db *sql.DB
}

// NewGradeService 创建成绩服务实例
func NewGradeService(db *sql.DB) *GradeService {
	return &GradeService{db: db}
}

// CreateGrade 创建成绩
func (s *GradeService) CreateGrade(req domain.CreateGradeRequest) (*domain.Grade, error) {
	// 检查学生是否存在
	var studentCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE id = $1", req.StudentID).Scan(&studentCount)
	if err != nil {
		return nil, fmt.Errorf("检查学生存在性失败: %v", err)
	}
	if studentCount == 0 {
		return nil, fmt.Errorf("学生ID %d 不存在", req.StudentID)
	}

	// 检查是否已存在该学生的成绩记录
	var existingCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM grades WHERE student_id = $1", req.StudentID).Scan(&existingCount)
	if err != nil {
		return nil, fmt.Errorf("检查成绩记录存在性失败: %v", err)
	}
	if existingCount > 0 {
		return nil, fmt.Errorf("学生ID %d 的成绩记录已存在，请使用更新功能", req.StudentID)
	}

	// 验证老师ID的有效性
	if err := s.validateTeacherIDs(req.ChineseTeacherID, req.MathTeacherID, req.EnglishTeacherID, req.SportsTeacherID, req.MusicTeacherID); err != nil {
		return nil, err
	}

	now := time.Now()
	query := `INSERT INTO grades (student_id, chinese_score, math_score, english_score, sports_score, music_score,
			   chinese_teacher_id, math_teacher_id, english_teacher_id, sports_teacher_id, music_teacher_id, created_at, updated_at) 
			   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	var id int
	err = s.db.QueryRow(query, req.StudentID, req.ChineseScore, req.MathScore, req.EnglishScore, req.SportsScore, req.MusicScore,
		req.ChineseTeacherID, req.MathTeacherID, req.EnglishTeacherID, req.SportsTeacherID, req.MusicTeacherID, now, now).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("创建成绩失败: %v", err)
	}

	grade := &domain.Grade{
		ID:               id,
		StudentID:        req.StudentID,
		ChineseScore:     req.ChineseScore,
		MathScore:        req.MathScore,
		EnglishScore:     req.EnglishScore,
		SportsScore:      req.SportsScore,
		MusicScore:       req.MusicScore,
		ChineseTeacherID: req.ChineseTeacherID,
		MathTeacherID:    req.MathTeacherID,
		EnglishTeacherID: req.EnglishTeacherID,
		SportsTeacherID:  req.SportsTeacherID,
		MusicTeacherID:   req.MusicTeacherID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return grade, nil
}

// validateTeacherIDs 验证老师ID的有效性
func (s *GradeService) validateTeacherIDs(chineseTeacherID, mathTeacherID, englishTeacherID, sportsTeacherID, musicTeacherID *int) error {
	teacherIDs := map[string]*int{
		"体育": sportsTeacherID,
		"音乐": musicTeacherID,
		"语文": chineseTeacherID,
		"数学": mathTeacherID,
		"英语": englishTeacherID,
	}

	for subject, teacherID := range teacherIDs {
		if teacherID != nil {
			// 检查老师是否存在且教授对应科目
			var teacherSubject string
			err := s.db.QueryRow("SELECT subject FROM teachers WHERE id = $1", *teacherID).Scan(&teacherSubject)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("老师ID %d 不存在", *teacherID)
				}
				return fmt.Errorf("检查老师存在性失败: %v", err)
			}
			if teacherSubject != subject {
				return fmt.Errorf("老师ID %d 不教授%s，实际教授: %s", *teacherID, subject, teacherSubject)
			}
		}
	}
	return nil
}

// GetAllGrades 获取成绩列表（支持筛选和分页）
func (s *GradeService) GetAllGrades(params domain.GradeQueryParams) ([]domain.GradeWithDetails, int, error) {
	// 构建查询条件
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.StudentID != nil {
		conditions = append(conditions, fmt.Sprintf("g.student_id = $%d", argIndex))
		args = append(args, *params.StudentID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 获取总数
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM grades g %s`, whereClause)
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("获取成绩总数失败: %v", err)
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}

	// 计算偏移量
	offset := (params.Page - 1) * params.Size

	// 查询成绩列表
	query := fmt.Sprintf(`SELECT g.id, g.student_id, g.chinese_score, g.math_score, g.english_score, g.sports_score, g.music_score,
			 g.chinese_teacher_id, g.math_teacher_id, g.english_teacher_id,
			 g.sports_teacher_id, g.music_teacher_id,
			 g.created_at, g.updated_at, s.name as student_name,
			 ct.name as chinese_teacher_name, mt.name as math_teacher_name, et.name as english_teacher_name,
			 st.name as sports_teacher_name, mut.name as music_teacher_name
			 FROM grades g 
			 LEFT JOIN students s ON g.student_id = s.id 
			 LEFT JOIN teachers ct ON g.chinese_teacher_id = ct.id
			 LEFT JOIN teachers mt ON g.math_teacher_id = mt.id
			 LEFT JOIN teachers et ON g.english_teacher_id = et.id
			 LEFT JOIN teachers st ON g.sports_teacher_id = st.id
			 LEFT JOIN teachers mut ON g.music_teacher_id = mut.id
			 %s ORDER BY g.created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	queryArgs := append(args, params.Size, offset)
	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询成绩列表失败: %v", err)
	}
	defer rows.Close()

	var grades []domain.GradeWithDetails
	for rows.Next() {
		var grade domain.GradeWithDetails
		var chineseTeacherName, mathTeacherName, englishTeacherName, sportsTeacherName, musicTeacherName sql.NullString
		err := rows.Scan(
			&grade.ID, &grade.StudentID, &grade.ChineseScore, &grade.MathScore, &grade.EnglishScore, &grade.SportsScore, &grade.MusicScore,
			&grade.ChineseTeacherID, &grade.MathTeacherID, &grade.EnglishTeacherID,
			&grade.SportsTeacherID, &grade.MusicTeacherID,
			&grade.CreatedAt, &grade.UpdatedAt, &grade.StudentName,
			&chineseTeacherName, &mathTeacherName, &englishTeacherName,
			&sportsTeacherName, &musicTeacherName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("扫描成绩数据失败: %v", err)
		}

		// 处理可能为空的老师姓名
		if chineseTeacherName.Valid {
			grade.ChineseTeacherName = chineseTeacherName.String
		}
		if mathTeacherName.Valid {
			grade.MathTeacherName = mathTeacherName.String
		}
		if englishTeacherName.Valid {
			grade.EnglishTeacherName = englishTeacherName.String
		}
		if sportsTeacherName.Valid {
			grade.SportsTeacherName = sportsTeacherName.String
		}
		if musicTeacherName.Valid {
			grade.MusicTeacherName = musicTeacherName.String
		}

		grades = append(grades, grade)
	}

	return grades, total, nil
}

// GetGradeByID 根据ID获取成绩
func (s *GradeService) GetGradeByID(id int) (*domain.GradeWithDetails, error) {
	query := `SELECT g.id, g.student_id, g.chinese_score, g.math_score, g.english_score, g.sports_score, g.music_score,
			 g.chinese_teacher_id, g.math_teacher_id, g.english_teacher_id,
			 g.sports_teacher_id, g.music_teacher_id,
			 g.created_at, g.updated_at, s.name as student_name,
			 ch.name as chinese_teacher_name, m.name as math_teacher_name, e.name as english_teacher_name,
			 st.name as sports_teacher_name, mt.name as music_teacher_name
			 FROM grades g 
			 LEFT JOIN students s ON g.student_id = s.id 
			 LEFT JOIN teachers ch ON g.chinese_teacher_id = ch.id
			 LEFT JOIN teachers m ON g.math_teacher_id = m.id
			 LEFT JOIN teachers e ON g.english_teacher_id = e.id
			 LEFT JOIN teachers st ON g.sports_teacher_id = st.id
			 LEFT JOIN teachers mt ON g.music_teacher_id = mt.id
			 WHERE g.id = $1`

	var grade domain.GradeWithDetails
	var chineseTeacherName, mathTeacherName, englishTeacherName, sportsTeacherName, musicTeacherName sql.NullString
	err := s.db.QueryRow(query, id).Scan(
		&grade.ID, &grade.StudentID, &grade.ChineseScore, &grade.MathScore, &grade.EnglishScore, &grade.SportsScore, &grade.MusicScore,
		&grade.ChineseTeacherID, &grade.MathTeacherID, &grade.EnglishTeacherID,
		&grade.SportsTeacherID, &grade.MusicTeacherID,
		&grade.CreatedAt, &grade.UpdatedAt, &grade.StudentName,
		&chineseTeacherName, &mathTeacherName, &englishTeacherName,
		&sportsTeacherName, &musicTeacherName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("成绩不存在")
		}
		return nil, fmt.Errorf("查询成绩失败: %v", err)
	}

	// 处理可能为空的老师姓名
	if chineseTeacherName.Valid {
		grade.ChineseTeacherName = chineseTeacherName.String
	}
	if mathTeacherName.Valid {
		grade.MathTeacherName = mathTeacherName.String
	}
	if englishTeacherName.Valid {
		grade.EnglishTeacherName = englishTeacherName.String
	}
	if sportsTeacherName.Valid {
		grade.SportsTeacherName = sportsTeacherName.String
	}
	if musicTeacherName.Valid {
		grade.MusicTeacherName = musicTeacherName.String
	}

	return &grade, nil
}

// UpdateGrade 更新成绩
func (s *GradeService) UpdateGrade(id int, req domain.UpdateGradeRequest) (*domain.Grade, error) {
	// 检查成绩是否存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM grades WHERE id = $1", id).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("检查成绩存在性失败: %v", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("成绩不存在")
	}

	// 验证老师ID的有效性
	if err := s.validateTeacherIDs(req.ChineseTeacherID, req.MathTeacherID, req.EnglishTeacherID, req.SportsTeacherID, req.MusicTeacherID); err != nil {
		return nil, err
	}

	// 构建更新字段
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.ChineseScore != nil {
		setParts = append(setParts, fmt.Sprintf("chinese_score = $%d", argIndex))
		args = append(args, *req.ChineseScore)
		argIndex++
	}

	if req.MathScore != nil {
		setParts = append(setParts, fmt.Sprintf("math_score = $%d", argIndex))
		args = append(args, *req.MathScore)
		argIndex++
	}

	if req.EnglishScore != nil {
		setParts = append(setParts, fmt.Sprintf("english_score = $%d", argIndex))
		args = append(args, *req.EnglishScore)
		argIndex++
	}

	if req.SportsScore != nil {
		setParts = append(setParts, fmt.Sprintf("sports_score = $%d", argIndex))
		args = append(args, *req.SportsScore)
		argIndex++
	}

	if req.MusicScore != nil {
		setParts = append(setParts, fmt.Sprintf("music_score = $%d", argIndex))
		args = append(args, *req.MusicScore)
		argIndex++
	}

	if req.ChineseTeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("chinese_teacher_id = $%d", argIndex))
		args = append(args, *req.ChineseTeacherID)
		argIndex++
	}

	if req.MathTeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("math_teacher_id = $%d", argIndex))
		args = append(args, *req.MathTeacherID)
		argIndex++
	}

	if req.EnglishTeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("english_teacher_id = $%d", argIndex))
		args = append(args, *req.EnglishTeacherID)
		argIndex++
	}

	if req.SportsTeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("sports_teacher_id = $%d", argIndex))
		args = append(args, *req.SportsTeacherID)
		argIndex++
	}

	if req.MusicTeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("music_teacher_id = $%d", argIndex))
		args = append(args, *req.MusicTeacherID)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("没有要更新的字段")
	}

	// 添加更新时间
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++
	args = append(args, id)

	query := fmt.Sprintf("UPDATE grades SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("更新成绩失败: %v", err)
	}

	// 获取更新后的成绩
	gradeDetails, err := s.GetGradeByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取更新后的成绩失败: %v", err)
	}

	return &gradeDetails.Grade, nil
}

// DeleteGrade 删除成绩
func (s *GradeService) DeleteGrade(id int) error {
	// 检查成绩是否存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM grades WHERE id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查成绩存在性失败: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("成绩不存在")
	}

	_, err = s.db.Exec("DELETE FROM grades WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("删除成绩失败: %v", err)
	}

	return nil
}
