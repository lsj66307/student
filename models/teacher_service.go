package models

import (
	"database/sql"
	"fmt"
	"student-management-system/database"
	"time"
)

// TeacherService 老师服务结构
type TeacherService struct {
	db *sql.DB
}

// NewTeacherService 创建新的老师服务实例
func NewTeacherService() *TeacherService {
	return &TeacherService{
		db: database.DB,
	}
}

// CreateTeacher 创建新老师
func (t *TeacherService) CreateTeacher(req CreateTeacherRequest) (*Teacher, error) {
	query := `
		INSERT INTO teachers (name, age, gender, email, phone, subject, title, department)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, age, gender, email, phone, subject, title, department, created_at, updated_at
	`

	teacher := &Teacher{}
	err := t.db.QueryRow(query, req.Name, req.Age, req.Gender, req.Email, req.Phone, req.Subject, req.Title, req.Department).Scan(
		&teacher.ID, &teacher.Name, &teacher.Age, &teacher.Gender,
		&teacher.Email, &teacher.Phone, &teacher.Subject, &teacher.Title,
		&teacher.Department, &teacher.CreatedAt, &teacher.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create teacher: %v", err)
	}

	return teacher, nil
}

// GetTeacherByID 根据ID获取老师信息
func (t *TeacherService) GetTeacherByID(id int) (*Teacher, error) {
	query := `
		SELECT id, name, age, gender, email, phone, subject, title, department, created_at, updated_at
		FROM teachers
		WHERE id = $1
	`

	teacher := &Teacher{}
	err := t.db.QueryRow(query, id).Scan(
		&teacher.ID, &teacher.Name, &teacher.Age, &teacher.Gender,
		&teacher.Email, &teacher.Phone, &teacher.Subject, &teacher.Title,
		&teacher.Department, &teacher.CreatedAt, &teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found")
		}
		return nil, fmt.Errorf("failed to get teacher: %v", err)
	}

	return teacher, nil
}

// GetAllTeachers 获取所有老师列表
func (t *TeacherService) GetAllTeachers(page, pageSize int) ([]*Teacher, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM teachers"
	var total int
	err := t.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count teachers: %v", err)
	}

	// 获取分页数据
	query := `
		SELECT id, name, age, gender, email, phone, subject, title, department, created_at, updated_at
		FROM teachers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := t.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query teachers: %v", err)
	}
	defer rows.Close()

	var teachers []*Teacher
	for rows.Next() {
		teacher := &Teacher{}
		err := rows.Scan(
			&teacher.ID, &teacher.Name, &teacher.Age, &teacher.Gender,
			&teacher.Email, &teacher.Phone, &teacher.Subject, &teacher.Title,
			&teacher.Department, &teacher.CreatedAt, &teacher.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan teacher: %v", err)
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %v", err)
	}

	return teachers, total, nil
}

// UpdateTeacher 更新老师信息
func (t *TeacherService) UpdateTeacher(id int, req UpdateTeacherRequest) (*Teacher, error) {
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
	if req.Subject != "" {
		setClauses = append(setClauses, fmt.Sprintf("subject = $%d", argIndex))
		args = append(args, req.Subject)
		argIndex++
	}
	if req.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, req.Title)
		argIndex++
	}
	if req.Department != "" {
		setClauses = append(setClauses, fmt.Sprintf("department = $%d", argIndex))
		args = append(args, req.Department)
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
		UPDATE teachers
		SET %s
		WHERE id = $%d
		RETURNING id, name, age, gender, email, phone, subject, title, department, created_at, updated_at
	`, fmt.Sprintf("%s", setClauses), argIndex)

	teacher := &Teacher{}
	err := t.db.QueryRow(query, args...).Scan(
		&teacher.ID, &teacher.Name, &teacher.Age, &teacher.Gender,
		&teacher.Email, &teacher.Phone, &teacher.Subject, &teacher.Title,
		&teacher.Department, &teacher.CreatedAt, &teacher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teacher not found")
		}
		return nil, fmt.Errorf("failed to update teacher: %v", err)
	}

	return teacher, nil
}

// DeleteTeacher 删除老师
func (t *TeacherService) DeleteTeacher(id int) error {
	query := "DELETE FROM teachers WHERE id = $1"
	result, err := t.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("teacher not found")
	}

	return nil
}