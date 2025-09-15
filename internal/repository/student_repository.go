package repository

import (
	"database/sql"
	"student-management-system/internal/domain"
	"student-management-system/pkg/logger"
)

// StudentRepository 学生仓储接口
type StudentRepository interface {
	Create(student *domain.Student) error
	GetByID(id int) (*domain.Student, error)
	GetByStudentID(studentID string) (*domain.Student, error)
	Update(student *domain.Student) error
	UpdateMajor(studentID int, newMajor string) error
	Delete(id int) error
	List(offset, limit int) ([]*domain.Student, error)
	Count() (int, error)
	BatchCreate(students []*domain.Student) error
	BatchDelete(ids []int) error
}

// studentRepository 学生仓储实现
type studentRepository struct {
	db *sql.DB
}

// NewStudentRepository 创建学生仓储实例
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

// Create 创建学生
func (r *studentRepository) Create(student *domain.Student) error {
	logger.WithFields(map[string]interface{}{
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Creating student")

	query := `
		INSERT INTO students (student_id, name, age, gender, phone, email, address, major, enrollment_date, graduation_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		student.StudentID,
		student.Name,
		student.Age,
		student.Gender,
		student.Phone,
		student.Email,
		student.Address,
		student.Major,
		student.EnrollmentDate,
		student.GraduationDate,
		student.Status,
	).Scan(&student.ID, &student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": student.StudentID,
			"name":       student.Name,
		}).Error("Failed to create student")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"id":         student.ID,
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Student created successfully")

	return nil
}

// GetByID 根据ID获取学生
func (r *studentRepository) GetByID(id int) (*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Getting student by ID")

	student := &domain.Student{}
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students 
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.StudentID,
		&student.Name,
		&student.Age,
		&student.Gender,
		&student.Phone,
		&student.Email,
		&student.Address,
		&student.Major,
		&student.EnrollmentDate,
		&student.GraduationDate,
		&student.Status,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.WithFields(map[string]interface{}{
			"id": id,
		}).Warn("Student not found")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"id": id,
		}).Error("Failed to get student by ID")
		return nil, err
	}

	logger.WithFields(map[string]interface{}{
		"id":         student.ID,
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Student retrieved successfully")

	return student, nil
}

// GetByStudentID 根据学号获取学生
func (r *studentRepository) GetByStudentID(studentID string) (*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"student_id": studentID,
	}).Info("Getting student by student ID")

	student := &domain.Student{}
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students 
		WHERE student_id = $1
	`

	err := r.db.QueryRow(query, studentID).Scan(
		&student.ID,
		&student.StudentID,
		&student.Name,
		&student.Age,
		&student.Gender,
		&student.Phone,
		&student.Email,
		&student.Address,
		&student.Major,
		&student.EnrollmentDate,
		&student.GraduationDate,
		&student.Status,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.WithFields(map[string]interface{}{
			"student_id": studentID,
		}).Warn("Student not found by student ID")
		return nil, nil
	}

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": studentID,
		}).Error("Failed to get student by student ID")
		return nil, err
	}

	logger.WithFields(map[string]interface{}{
		"id":         student.ID,
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Student retrieved by student ID successfully")

	return student, nil
}

// Update 更新学生信息
func (r *studentRepository) Update(student *domain.Student) error {
	logger.WithFields(map[string]interface{}{
		"id":         student.ID,
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Updating student")

	query := `
		UPDATE students 
		SET student_id = $2, name = $3, age = $4, gender = $5, phone = $6, 
		    email = $7, address = $8, major = $9, enrollment_date = $10, 
		    graduation_date = $11, status = $12
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		student.ID,
		student.StudentID,
		student.Name,
		student.Age,
		student.Gender,
		student.Phone,
		student.Email,
		student.Address,
		student.Major,
		student.EnrollmentDate,
		student.GraduationDate,
		student.Status,
	)

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"id":         student.ID,
			"student_id": student.StudentID,
			"name":       student.Name,
		}).Error("Failed to update student")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"id":         student.ID,
		"student_id": student.StudentID,
		"name":       student.Name,
	}).Info("Student updated successfully")

	return nil
}

// Delete 删除学生
func (r *studentRepository) Delete(id int) error {
	logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting student")

	query := `DELETE FROM students WHERE id = $1`
	_, err := r.db.Exec(query, id)

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"id": id,
		}).Error("Failed to delete student")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Student deleted successfully")

	return nil
}

// List 获取学生列表
func (r *studentRepository) List(offset, limit int) ([]*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	}).Info("Getting student list")

	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"offset": offset,
			"limit":  limit,
		}).Error("Failed to query student list")
		return nil, err
	}
	defer rows.Close()

	var students []*domain.Student
	for rows.Next() {
		student := &domain.Student{}
		err := rows.Scan(
			&student.ID,
			&student.StudentID,
			&student.Name,
			&student.Age,
			&student.Gender,
			&student.Phone,
			&student.Email,
			&student.Address,
			&student.Major,
			&student.EnrollmentDate,
			&student.GraduationDate,
			&student.Status,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
			logger.WithError(err).Error("Failed to scan student row")
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("Error iterating student rows")
		return nil, err
	}

	logger.WithFields(map[string]interface{}{
		"count":  len(students),
		"offset": offset,
		"limit":  limit,
	}).Info("Student list retrieved successfully")

	return students, nil
}

// Count 获取学生总数
func (r *studentRepository) Count() (int, error) {
	logger.Info("Getting student count")

	query := `SELECT COUNT(*) FROM students`
	var count int
	err := r.db.QueryRow(query).Scan(&count)

	if err != nil {
		logger.WithError(err).Error("Failed to get student count")
		return 0, err
	}

	logger.WithFields(map[string]interface{}{
		"count": count,
	}).Info("Student count retrieved successfully")

	return count, nil
}

// BatchCreate 批量创建学生
func (r *studentRepository) BatchCreate(students []*domain.Student) error {
	logger.WithFields(map[string]interface{}{
		"count": len(students),
	}).Info("Batch creating students")

	tx, err := r.db.Begin()
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction for batch create")
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO students (student_id, name, age, gender, phone, email, address, major, enrollment_date, graduation_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	for i, student := range students {
		err := tx.QueryRow(
			query,
			student.StudentID,
			student.Name,
			student.Age,
			student.Gender,
			student.Phone,
			student.Email,
			student.Address,
			student.Major,
			student.EnrollmentDate,
			student.GraduationDate,
			student.Status,
		).Scan(&student.ID, &student.CreatedAt, &student.UpdatedAt)

		if err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"index":      i,
				"student_id": student.StudentID,
				"name":       student.Name,
			}).Error("Failed to create student in batch")
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit batch create transaction")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"count": len(students),
	}).Info("Batch create students completed successfully")

	return nil
}

// BatchDelete 批量删除学生
func (r *studentRepository) BatchDelete(ids []int) error {
	if len(ids) == 0 {
		logger.Warn("No IDs provided for batch delete")
		return nil
	}

	logger.WithFields(map[string]interface{}{
		"count": len(ids),
		"ids":   ids,
	}).Info("Batch deleting students")

	query := `DELETE FROM students WHERE id = ANY($1)`
	_, err := r.db.Exec(query, ids)

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"count": len(ids),
			"ids":   ids,
		}).Error("Failed to batch delete students")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"count": len(ids),
		"ids":   ids,
	}).Info("Batch delete students completed successfully")

	return nil
}

// UpdateMajor 更新学生专业
func (r *studentRepository) UpdateMajor(studentID int, newMajor string) error {
	logger.WithFields(map[string]interface{}{
		"student_id": studentID,
		"new_major":  newMajor,
	}).Info("Updating student major")

	query := `UPDATE students SET major = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(query, newMajor, studentID)

	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": studentID,
			"new_major":  newMajor,
		}).Error("Failed to update student major")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": studentID,
			"new_major":  newMajor,
		}).Error("Failed to get rows affected after updating student major")
		return err
	}

	if rowsAffected == 0 {
		logger.WithFields(map[string]interface{}{
			"student_id": studentID,
			"new_major":  newMajor,
		}).Warn("No student found with the given ID for major update")
		return sql.ErrNoRows
	}

	logger.WithFields(map[string]interface{}{
		"student_id": studentID,
		"new_major":  newMajor,
	}).Info("Student major updated successfully")

	return nil
}
