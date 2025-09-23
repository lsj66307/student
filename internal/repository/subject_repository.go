package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"student-management-system/internal/domain"
	"student-management-system/pkg/logger"
)

// SubjectRepository 科目仓储接口
type SubjectRepository interface {
	Create(subject *domain.Subject) error
	GetByID(id int) (*domain.Subject, error)
	GetByCode(code string) (*domain.Subject, error)
	Update(subject *domain.Subject) error
	Delete(id int) error
	List(req *domain.SubjectListRequest) ([]*domain.Subject, int64, error)
	GetActiveSubjects() ([]*domain.Subject, error)
	ExistsByCode(code string) (bool, error)
	ExistsByCodeExcludeID(code string, id int) (bool, error)
}

// subjectRepository 科目仓储实现
type subjectRepository struct {
	db *sql.DB
}

// NewSubjectRepository 创建科目仓储实例
func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

// Create 创建科目
func (r *subjectRepository) Create(subject *domain.Subject) error {
	logger.WithFields(map[string]interface{}{
		"name": subject.Name,
		"code": subject.Code,
	}).Info("Creating subject")

	query := `
		INSERT INTO subjects (name, code, description, credits, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		subject.Name,
		subject.Code,
		subject.Description,
		subject.Credits,
		subject.Status,
	).Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)

	if err != nil {
		logger.WithError(err).Error("Failed to create subject")
		return fmt.Errorf("创建科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": subject.ID,
	}).Info("Subject created successfully")
	return nil
}

// GetByID 根据ID获取科目
func (r *subjectRepository) GetByID(id int) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Getting subject by ID")

	query := `
		SELECT id, name, code, description, credits, status, created_at, updated_at
		FROM subjects
		WHERE id = $1
	`

	subject := &domain.Subject{}
	err := r.db.QueryRow(query, id).Scan(
		&subject.ID,
		&subject.Name,
		&subject.Code,
		&subject.Description,
		&subject.Credits,
		&subject.Status,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithFields(map[string]interface{}{
				"subject_id": id,
			}).Warn("Subject not found")
			return nil, fmt.Errorf("科目不存在")
		}
		logger.WithError(err).Error("Failed to get subject by ID")
		return nil, fmt.Errorf("获取科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Subject retrieved successfully")
	return subject, nil
}

// GetByCode 根据代码获取科目
func (r *subjectRepository) GetByCode(code string) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"subject_code": code,
	}).Info("Getting subject by code")

	query := `
		SELECT id, name, code, description, credits, status, created_at, updated_at
		FROM subjects
		WHERE code = $1
	`

	subject := &domain.Subject{}
	err := r.db.QueryRow(query, code).Scan(
		&subject.ID,
		&subject.Name,
		&subject.Code,
		&subject.Description,
		&subject.Credits,
		&subject.Status,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithFields(map[string]interface{}{
				"subject_code": code,
			}).Warn("Subject not found")
			return nil, fmt.Errorf("科目不存在")
		}
		logger.WithError(err).Error("Failed to get subject by code")
		return nil, fmt.Errorf("获取科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_code": code,
	}).Info("Subject retrieved successfully")
	return subject, nil
}

// Update 更新科目
func (r *subjectRepository) Update(subject *domain.Subject) error {
	logger.WithFields(map[string]interface{}{
		"subject_id": subject.ID,
		"name":       subject.Name,
		"code":       subject.Code,
	}).Info("Updating subject")

	query := `
		UPDATE subjects 
		SET name = $1, code = $2, description = $3, credits = $4, status = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	result, err := r.db.Exec(
		query,
		subject.Name,
		subject.Code,
		subject.Description,
		subject.Credits,
		subject.Status,
		subject.ID,
	)

	if err != nil {
		logger.WithError(err).Error("Failed to update subject")
		return fmt.Errorf("更新科目失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("获取更新结果失败: %v", err)
	}

	if rowsAffected == 0 {
		logger.WithFields(map[string]interface{}{
			"subject_id": subject.ID,
		}).Warn("No rows affected during update")
		return fmt.Errorf("科目不存在或未发生变更")
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": subject.ID,
	}).Info("Subject updated successfully")
	return nil
}

// Delete 删除科目
func (r *subjectRepository) Delete(id int) error {
	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Deleting subject")

	query := `DELETE FROM subjects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.WithError(err).Error("Failed to delete subject")
		return fmt.Errorf("删除科目失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("获取删除结果失败: %v", err)
	}

	if rowsAffected == 0 {
		logger.WithFields(map[string]interface{}{
			"subject_id": id,
		}).Warn("No rows affected during delete")
		return fmt.Errorf("科目不存在")
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Subject deleted successfully")
	return nil
}

// List 获取科目列表
func (r *subjectRepository) List(req *domain.SubjectListRequest) ([]*domain.Subject, int64, error) {
	logger.WithFields(map[string]interface{}{
		"page": req.Page,
		"size": req.Size,
	}).Info("Getting subject list")

	// 构建查询条件
	var conditions []string
	var args []interface{}
	paramIndex := 1

	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name LIKE $%d", paramIndex))
		args = append(args, "%"+req.Name+"%")
		paramIndex++
	}

	if req.Code != "" {
		conditions = append(conditions, fmt.Sprintf("code LIKE $%d", paramIndex))
		args = append(args, "%"+req.Code+"%")
		paramIndex++
	}

	if req.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramIndex))
		args = append(args, req.Status)
		paramIndex++
	}

	if req.Credits > 0 {
		conditions = append(conditions, fmt.Sprintf("credits = $%d", paramIndex))
		args = append(args, req.Credits)
		paramIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subjects %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		logger.WithError(err).Error("Failed to count subjects")
		return nil, 0, fmt.Errorf("获取科目总数失败: %v", err)
	}

	// 获取列表数据
	offset := (req.Page - 1) * req.Size
	listQuery := fmt.Sprintf(`
		SELECT id, name, code, description, credits, status, created_at, updated_at
		FROM subjects %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramIndex, paramIndex+1)

	listArgs := append(args, req.Size, offset)
	rows, err := r.db.Query(listQuery, listArgs...)
	if err != nil {
		logger.WithError(err).Error("Failed to query subjects")
		return nil, 0, fmt.Errorf("查询科目列表失败: %v", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		subject := &domain.Subject{}
		err := rows.Scan(
			&subject.ID,
			&subject.Name,
			&subject.Code,
			&subject.Description,
			&subject.Credits,
			&subject.Status,
			&subject.CreatedAt,
			&subject.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("Failed to scan subject row")
			return nil, 0, fmt.Errorf("扫描科目数据失败: %v", err)
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("Error occurred during rows iteration")
		return nil, 0, fmt.Errorf("遍历科目数据失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"count": len(subjects),
		"total": total,
	}).Info("Subject list retrieved successfully")

	return subjects, total, nil
}

// GetActiveSubjects 获取所有活跃的科目
func (r *subjectRepository) GetActiveSubjects() ([]*domain.Subject, error) {
	logger.Info("Getting active subjects")

	query := `
		SELECT id, name, code, description, credits, status, created_at, updated_at
		FROM subjects
		WHERE status = 'active'
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.WithError(err).Error("Failed to query active subjects")
		return nil, fmt.Errorf("查询活跃科目失败: %v", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		subject := &domain.Subject{}
		err := rows.Scan(
			&subject.ID,
			&subject.Name,
			&subject.Code,
			&subject.Description,
			&subject.Credits,
			&subject.Status,
			&subject.CreatedAt,
			&subject.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("Failed to scan subject row")
			return nil, fmt.Errorf("扫描科目数据失败: %v", err)
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("Error occurred during rows iteration")
		return nil, fmt.Errorf("遍历科目数据失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"count": len(subjects),
	}).Info("Active subjects retrieved successfully")
	return subjects, nil
}

// ExistsByCode 检查科目代码是否存在
func (r *subjectRepository) ExistsByCode(code string) (bool, error) {
	logger.WithFields(map[string]interface{}{
		"subject_code": code,
	}).Info("Checking if subject code exists")

	query := `SELECT COUNT(*) FROM subjects WHERE code = $1`
	var count int
	err := r.db.QueryRow(query, code).Scan(&count)
	if err != nil {
		logger.WithError(err).Error("Failed to check subject code existence")
		return false, fmt.Errorf("检查科目代码失败: %v", err)
	}

	exists := count > 0
	logger.WithFields(map[string]interface{}{
		"subject_code": code,
		"exists":       exists,
	}).Info("Subject code existence check completed")

	return exists, nil
}

// ExistsByCodeExcludeID 检查科目代码是否存在（排除指定ID）
func (r *subjectRepository) ExistsByCodeExcludeID(code string, id int) (bool, error) {
	logger.WithFields(map[string]interface{}{
		"subject_code": code,
		"exclude_id":   id,
	}).Info("Checking if subject code exists excluding ID")

	query := `SELECT COUNT(*) FROM subjects WHERE code = $1 AND id != $2`
	var count int
	err := r.db.QueryRow(query, code, id).Scan(&count)
	if err != nil {
		logger.WithError(err).Error("Failed to check subject code existence")
		return false, fmt.Errorf("检查科目代码失败: %v", err)
	}

	exists := count > 0
	logger.WithFields(map[string]interface{}{
		"subject_code": code,
		"exclude_id":   id,
		"exists":       exists,
	}).Info("Subject code existence check completed")

	return exists, nil
}