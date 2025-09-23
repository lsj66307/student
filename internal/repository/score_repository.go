package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"student-management-system/internal/domain"
	"student-management-system/pkg/logger"
)

// ScoreRepository 成绩仓储接口
type ScoreRepository interface {
	Create(score *domain.Score) error
	GetByID(id int) (*domain.Score, error)
	GetByStudentAndSubject(studentID, subjectID int) (*domain.Score, error)
	Update(score *domain.Score) error
	Delete(id int) error
	List(req *domain.ScoreListRequest) ([]*domain.Score, int64, error)
}

// scoreRepository 成绩仓储实现
type scoreRepository struct {
	db *sql.DB
}

// NewScoreRepository 创建成绩仓储实例
func NewScoreRepository(db *sql.DB) ScoreRepository {
	return &scoreRepository{db: db}
}

// Create 创建成绩
func (r *scoreRepository) Create(score *domain.Score) error {
	logger.Info("Creating new score")

	query := `
		INSERT INTO scores (student_id, subject_id, score, semester, exam_type, remarks, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, score.StudentID, score.SubjectID, score.Score,
		score.Semester, score.ExamType, score.Remarks).Scan(&id)
	if err != nil {
		logger.Error("Failed to create score", "error", err)
		return fmt.Errorf("failed to create score: %w", err)
	}

	score.ID = id
	logger.Info("Score created successfully", "score_id", id)
	return nil
}

// GetByID 根据ID获取成绩
func (r *scoreRepository) GetByID(id int) (*domain.Score, error) {
	query := `
		SELECT s.id, s.student_id, s.subject_id, s.score, s.semester, s.exam_type, s.remarks, s.created_at, s.updated_at,
		       st.name as student_name, st.student_id as student_code,
		       sub.name as subject_name, sub.code as subject_code
		FROM scores s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN subjects sub ON s.subject_id = sub.id
		WHERE s.id = $1
	`

	score := &domain.Score{}
	var studentName, studentCode, subjectName, subjectCode sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&score.ID, &score.StudentID, &score.SubjectID, &score.Score,
		&score.Semester, &score.ExamType, &score.Remarks, &score.CreatedAt, &score.UpdatedAt,
		&studentName, &studentCode, &subjectName, &subjectCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("score not found")
		}
		return nil, fmt.Errorf("failed to get score: %w", err)
	}

	// 设置关联数据
	if studentName.Valid {
		score.Student = &domain.Student{
			ID:        score.StudentID,
			Name:      studentName.String,
			StudentID: studentCode.String,
		}
	}

	if subjectName.Valid {
		score.Subject = &domain.Subject{
			ID:   score.SubjectID,
			Name: subjectName.String,
			Code: subjectCode.String,
		}
	}

	return score, nil
}

// GetByStudentAndSubject 根据学生ID和科目ID获取成绩
func (r *scoreRepository) GetByStudentAndSubject(studentID, subjectID int) (*domain.Score, error) {
	logger.Info("Getting score by student and subject", "student_id", studentID, "subject_id", subjectID)

	query := `
		SELECT s.id, s.student_id, s.subject_id, s.score, s.semester, s.exam_type, s.remarks, s.created_at, s.updated_at,
		       st.name as student_name, st.student_id as student_code,
		       sub.name as subject_name, sub.code as subject_code
		FROM scores s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN subjects sub ON s.subject_id = sub.id
		WHERE s.student_id = $1 AND s.subject_id = $2
	`

	score := &domain.Score{}
	var studentName, studentCode, subjectName, subjectCode sql.NullString

	err := r.db.QueryRow(query, studentID, subjectID).Scan(
		&score.ID, &score.StudentID, &score.SubjectID, &score.Score,
		&score.Semester, &score.ExamType, &score.Remarks, &score.CreatedAt, &score.UpdatedAt,
		&studentName, &studentCode, &subjectName, &subjectCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("Score not found", "student_id", studentID, "subject_id", subjectID)
			return nil, fmt.Errorf("score not found")
		}
		logger.Error("Failed to get score", "error", err, "student_id", studentID, "subject_id", subjectID)
		return nil, fmt.Errorf("failed to get score: %w", err)
	}

	// 设置关联数据
	if studentName.Valid {
		score.Student = &domain.Student{
			ID:        score.StudentID,
			Name:      studentName.String,
			StudentID: studentCode.String,
		}
	}

	if subjectName.Valid {
		score.Subject = &domain.Subject{
			ID:   score.SubjectID,
			Name: subjectName.String,
			Code: subjectCode.String,
		}
	}

	return score, nil
}

// Update 更新成绩
func (r *scoreRepository) Update(score *domain.Score) error {
	query := `
		UPDATE scores 
		SET score = $1, semester = $2, exam_type = $3, remarks = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	result, err := r.db.Exec(query, score.Score, score.Semester, score.ExamType, score.Remarks, score.ID)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("score not found")
	}

	return nil
}

// Delete 删除成绩
func (r *scoreRepository) Delete(id int) error {
	query := `DELETE FROM scores WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete score: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("score not found")
	}

	return nil
}

// List 获取成绩列表
func (r *scoreRepository) List(req *domain.ScoreListRequest) ([]*domain.Score, int64, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 构建查询条件
	var conditions []string
	var args []interface{}
	argIndex := 1

	if req.StudentID > 0 {
		conditions = append(conditions, fmt.Sprintf("s.student_id = $%d", argIndex))
		args = append(args, req.StudentID)
		argIndex++
	}

	if req.SubjectID > 0 {
		conditions = append(conditions, fmt.Sprintf("s.subject_id = $%d", argIndex))
		args = append(args, req.SubjectID)
		argIndex++
	}

	if req.Semester != "" {
		conditions = append(conditions, fmt.Sprintf("s.semester = $%d", argIndex))
		args = append(args, req.Semester)
		argIndex++
	}

	if req.ExamType != "" {
		conditions = append(conditions, fmt.Sprintf("s.exam_type = $%d", argIndex))
		args = append(args, req.ExamType)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM scores s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN subjects sub ON s.subject_id = sub.id
		%s
	`, whereClause)

	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count scores: %w", err)
	}

	// 查询数据
	offset := (req.Page - 1) * req.Size
	dataQuery := fmt.Sprintf(`
		SELECT s.id, s.student_id, s.subject_id, s.score, s.semester, s.exam_type, s.remarks, s.created_at, s.updated_at,
		       st.name as student_name, st.student_id as student_code,
		       sub.name as subject_name, sub.code as subject_code
		FROM scores s
		LEFT JOIN students st ON s.student_id = st.id
		LEFT JOIN subjects sub ON s.subject_id = sub.id
		%s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Size, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query scores: %w", err)
	}
	defer rows.Close()

	var scores []*domain.Score
	for rows.Next() {
		score := &domain.Score{}
		var studentName, studentCode, subjectName, subjectCode sql.NullString

		err := rows.Scan(
			&score.ID, &score.StudentID, &score.SubjectID, &score.Score,
			&score.Semester, &score.ExamType, &score.Remarks, &score.CreatedAt, &score.UpdatedAt,
			&studentName, &studentCode, &subjectName, &subjectCode,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan score: %w", err)
		}

		// 设置关联数据
		if studentName.Valid {
			score.Student = &domain.Student{
				ID:        score.StudentID,
				Name:      studentName.String,
				StudentID: studentCode.String,
			}
		}

		if subjectName.Valid {
			score.Subject = &domain.Subject{
				ID:   score.SubjectID,
				Name: subjectName.String,
				Code: subjectCode.String,
			}
		}

		scores = append(scores, score)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate scores: %w", err)
	}

	return scores, total, nil
}
