package service

import (
	"fmt"
	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/logger"
)

// StudentService 学生服务结构
type StudentService struct {
	repo repository.StudentRepository
}

// NewStudentService 创建新的学生服务实例
func NewStudentService() *StudentService {
	return &StudentService{
		repo: repository.NewStudentRepository(repository.DB),
	}
}

// CreateStudent 创建新学生
func (s *StudentService) CreateStudent(req domain.CreateStudentRequest) (*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"name":       req.Name,
		"email":      req.Email,
		"student_id": req.StudentID,
	}).Info("Creating new student")

	student := &domain.Student{
		StudentID:      req.StudentID,
		Name:           req.Name,
		Age:            req.Age,
		Gender:         req.Gender,
		Phone:          req.Phone,
		Email:          req.Email,
		Address:        req.Address,
		Major:          req.Major,
		EnrollmentDate: req.EnrollmentDate,
		GraduationDate: req.GraduationDate,
		Status:         req.Status,
	}

	// 如果状态为空，设置默认值
	if student.Status == "" {
		student.Status = "active"
	}

	err := s.repo.Create(student)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"name":       req.Name,
			"email":      req.Email,
			"student_id": req.StudentID,
		}).Error("Failed to create student")
		return nil, fmt.Errorf("failed to create student: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"student_id": student.ID,
		"name":       student.Name,
	}).Info("Student created successfully")

	return student, nil
}

// GetStudentByID 根据ID获取学生信息
func (s *StudentService) GetStudentByID(id int) (*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"student_id": id,
	}).Info("Getting student by ID")

	student, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": id,
		}).Error("Failed to get student")
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	if student == nil {
		logger.WithFields(map[string]interface{}{
			"student_id": id,
		}).Warn("Student not found")
		return nil, fmt.Errorf("student with ID %d not found", id)
	}

	return student, nil
}

// GetAllStudents 获取所有学生信息（分页）
func (s *StudentService) GetAllStudents(page, pageSize int) ([]*domain.Student, int, error) {
	logger.WithFields(map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}).Info("Getting all students")

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取学生列表
	students, err := s.repo.List(offset, pageSize)
	if err != nil {
		logger.WithError(err).Error("Failed to get students list")
		return nil, 0, fmt.Errorf("failed to get students: %v", err)
	}

	// 获取总数
	total, err := s.repo.Count()
	if err != nil {
		logger.WithError(err).Error("Failed to get students count")
		return nil, 0, fmt.Errorf("failed to get students count: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"count": len(students),
		"total": total,
	}).Info("Students retrieved successfully")

	return students, total, nil
}

// UpdateStudent 更新学生信息
func (s *StudentService) UpdateStudent(id int, req domain.UpdateStudentRequest) (*domain.Student, error) {
	logger.WithFields(map[string]interface{}{
		"student_id": id,
		"name":       req.Name,
	}).Info("Updating student")

	// 先获取现有学生信息
	student, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": id,
		}).Error("Failed to get student for update")
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	if student == nil {
		return nil, fmt.Errorf("student with ID %d not found", id)
	}

	// 更新字段（只更新非空字段）
	if req.StudentID != "" {
		student.StudentID = req.StudentID
	}
	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Age > 0 {
		student.Age = req.Age
	}
	if req.Gender != "" {
		student.Gender = req.Gender
	}
	if req.Phone != "" {
		student.Phone = req.Phone
	}
	if req.Email != "" {
		student.Email = req.Email
	}
	if req.Address != "" {
		student.Address = req.Address
	}
	if req.Major != "" {
		student.Major = req.Major
	}
	if req.EnrollmentDate != nil {
		student.EnrollmentDate = req.EnrollmentDate
	}
	if req.GraduationDate != nil {
		student.GraduationDate = req.GraduationDate
	}
	if req.Status != "" {
		student.Status = req.Status
	}

	err = s.repo.Update(student)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": id,
		}).Error("Failed to update student")
		return nil, fmt.Errorf("failed to update student: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"student_id": id,
		"name":       student.Name,
	}).Info("Student updated successfully")

	return student, nil
}

// DeleteStudent 删除学生
func (s *StudentService) DeleteStudent(id int) error {
	logger.WithFields(map[string]interface{}{
		"student_id": id,
	}).Info("Deleting student")

	// 检查学生是否存在
	student, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": id,
		}).Error("Failed to get student for deletion")
		return fmt.Errorf("failed to get student: %v", err)
	}

	if student == nil {
		return fmt.Errorf("student with ID %d not found", id)
	}

	err = s.repo.Delete(id)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"student_id": id,
		}).Error("Failed to delete student")
		return fmt.Errorf("failed to delete student: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"student_id": id,
	}).Info("Student deleted successfully")

	return nil
}
