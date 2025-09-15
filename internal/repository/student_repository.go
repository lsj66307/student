package repository

import (
	"database/sql"
	"student-management-system/internal/domain"
)

// StudentRepository 学生仓储接口
type StudentRepository interface {
	Create(student *domain.Student) error
	GetByID(id int) (*domain.Student, error)
	GetByStudentID(studentID string) (*domain.Student, error)
	Update(student *domain.Student) error
	Delete(id int) error
	List(offset, limit int) ([]*domain.Student, error)
	Count() (int, error)
	BatchCreate(students []*domain.Student) error
	BatchDelete(ids []int) error
	GetByMajor(major string) ([]*domain.Student, error)
	UpdateMajor(id int, newMajor string) error
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

	return err
}

// GetByID 根据ID获取学生
func (r *studentRepository) GetByID(id int) (*domain.Student, error) {
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students WHERE id = $1
	`

	student := &domain.Student{}
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
		return nil, nil
	}

	return student, err
}

// GetByStudentID 根据学号获取学生
func (r *studentRepository) GetByStudentID(studentID string) (*domain.Student, error) {
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students WHERE student_id = $1
	`

	student := &domain.Student{}
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
		return nil, nil
	}

	return student, err
}

// Update 更新学生信息
func (r *studentRepository) Update(student *domain.Student) error {
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

	return err
}

// Delete 删除学生
func (r *studentRepository) Delete(id int) error {
	query := `DELETE FROM students WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List 获取学生列表
func (r *studentRepository) List(offset, limit int) ([]*domain.Student, error) {
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
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
			return nil, err
		}
		students = append(students, student)
	}

	return students, rows.Err()
}

// Count 获取学生总数
func (r *studentRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM students`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// BatchCreate 批量创建学生
func (r *studentRepository) BatchCreate(students []*domain.Student) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO students (student_id, name, age, gender, phone, email, address, major, enrollment_date, graduation_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	for _, student := range students {
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
			return err
		}
	}

	return tx.Commit()
}

// BatchDelete 批量删除学生
func (r *studentRepository) BatchDelete(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := `DELETE FROM students WHERE id = ANY($1)`
	_, err := r.db.Exec(query, ids)
	return err
}

// GetByMajor 根据专业获取学生列表
func (r *studentRepository) GetByMajor(major string) ([]*domain.Student, error) {
	query := `
		SELECT id, student_id, name, age, gender, phone, email, address, major, 
		       enrollment_date, graduation_date, status, created_at, updated_at
		FROM students 
		WHERE major = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, major)
	if err != nil {
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
			return nil, err
		}
		students = append(students, student)
	}

	return students, rows.Err()
}

// UpdateMajor 更新学生专业
func (r *studentRepository) UpdateMajor(id int, newMajor string) error {
	query := `UPDATE students SET major = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, newMajor)
	return err
}
