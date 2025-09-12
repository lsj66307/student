package models

import (
	"database/sql"
	"fmt"
	"student-management-system/database"
	"time"
)

// StudentService 学生服务结构
type StudentService struct {
	db *sql.DB
}

// NewStudentService 创建新的学生服务实例
func NewStudentService() *StudentService {
	return &StudentService{
		db: database.DB,
	}
}

// CreateStudent 创建新学生
func (s *StudentService) CreateStudent(req CreateStudentRequest) (*Student, error) {
	query := `
		INSERT INTO students (name, age, gender, email, phone, major, grade)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, age, gender, email, phone, major, grade, created_at, updated_at
	`

	student := &Student{}
	err := s.db.QueryRow(query, req.Name, req.Age, req.Gender, req.Email, req.Phone, req.Major, req.Grade).Scan(
		&student.ID, &student.Name, &student.Age, &student.Gender,
		&student.Email, &student.Phone, &student.Major, &student.Grade,
		&student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create student: %v", err)
	}

	return student, nil
}

// GetStudentByID 根据ID获取学生信息
func (s *StudentService) GetStudentByID(id int) (*Student, error) {
	query := `
		SELECT id, name, age, gender, email, phone, major, grade, created_at, updated_at
		FROM students
		WHERE id = $1
	`

	student := &Student{}
	err := s.db.QueryRow(query, id).Scan(
		&student.ID, &student.Name, &student.Age, &student.Gender,
		&student.Email, &student.Phone, &student.Major, &student.Grade,
		&student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	return student, nil
}

// GetAllStudents 获取所有学生列表
func (s *StudentService) GetAllStudents(page, pageSize int) ([]*Student, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM students"
	var total int
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count students: %v", err)
	}

	// 获取分页数据
	query := `
		SELECT id, name, age, gender, email, phone, major, grade, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query students: %v", err)
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		student := &Student{}
		err := rows.Scan(
			&student.ID, &student.Name, &student.Age, &student.Gender,
			&student.Email, &student.Phone, &student.Major, &student.Grade,
			&student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan student: %v", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %v", err)
	}

	return students, total, nil
}

// UpdateStudent 更新学生信息
func (s *StudentService) UpdateStudent(id int, req UpdateStudentRequest) (*Student, error) {
	// 构建动态更新查询
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}
	if req.Age > 0 {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, req.Age)
		argIndex++
	}
	if req.Gender != "" {
		setClauses = append(setClauses, fmt.Sprintf("gender = $%d", argIndex))
		args = append(args, req.Gender)
		argIndex++
	}
	if req.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}
	if req.Phone != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, req.Phone)
		argIndex++
	}
	if req.Major != "" {
		setClauses = append(setClauses, fmt.Sprintf("major = $%d", argIndex))
		args = append(args, req.Major)
		argIndex++
	}
	if req.Grade != "" {
		setClauses = append(setClauses, fmt.Sprintf("grade = $%d", argIndex))
		args = append(args, req.Grade)
		argIndex++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// 添加updated_at字段
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// 添加WHERE条件的ID
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE students
		SET %s
		WHERE id = $%d
		RETURNING id, name, age, gender, email, phone, major, grade, created_at, updated_at
	`, fmt.Sprintf("%s", setClauses), argIndex)

	student := &Student{}
	err := s.db.QueryRow(query, args...).Scan(
		&student.ID, &student.Name, &student.Age, &student.Gender,
		&student.Email, &student.Phone, &student.Major, &student.Grade,
		&student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to update student: %v", err)
	}

	return student, nil
}

// DeleteStudent 删除学生
func (s *StudentService) DeleteStudent(id int) error {
	query := "DELETE FROM students WHERE id = $1"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}
