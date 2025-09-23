package service

import (
	"fmt"
	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/logger"
)

// SubjectService 科目服务结构
type SubjectService struct {
	repo repository.SubjectRepository
}

// NewSubjectService 创建新的科目服务实例
func NewSubjectService() *SubjectService {
	return &SubjectService{
		repo: repository.NewSubjectRepository(repository.DB),
	}
}

// CreateSubject 创建新科目
func (s *SubjectService) CreateSubject(req domain.CreateSubjectRequest) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"name": req.Name,
		"code": req.Code,
	}).Info("Creating new subject")

	// 检查科目代码是否已存在
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		logger.WithError(err).Error("Failed to check subject code existence")
		return nil, fmt.Errorf("检查科目代码失败: %v", err)
	}

	if exists {
		logger.WithFields(map[string]interface{}{
			"code": req.Code,
		}).Warn("Subject code already exists")
		return nil, fmt.Errorf("科目代码 %s 已存在", req.Code)
	}

	subject := &domain.Subject{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Credits:     req.Credits,
		Status:      req.Status,
	}

	// 如果状态为空，设置默认值
	if subject.Status == "" {
		subject.Status = "active"
	}

	err = s.repo.Create(subject)
	if err != nil {
		logger.WithError(err).Error("Failed to create subject")
		return nil, fmt.Errorf("创建科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": subject.ID,
		"code":       subject.Code,
		"name":       subject.Name,
	}).Info("Subject created successfully")

	return subject, nil
}

// GetSubjectByID 根据ID获取科目
func (s *SubjectService) GetSubjectByID(id int) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Getting subject by ID")

	subject, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).Error("Failed to get subject by ID")
		return nil, fmt.Errorf("获取科目失败: %v", err)
	}

	if subject == nil {
		logger.WithFields(map[string]interface{}{
			"subject_id": id,
		}).Warn("Subject not found")
		return nil, fmt.Errorf("科目不存在")
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": id,
		"name":       subject.Name,
	}).Info("Subject retrieved successfully")

	return subject, nil
}

// GetSubjectByCode 根据代码获取科目
func (s *SubjectService) GetSubjectByCode(code string) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"code": code,
	}).Info("Getting subject by code")

	subject, err := s.repo.GetByCode(code)
	if err != nil {
		logger.WithError(err).Error("Failed to get subject by code")
		return nil, fmt.Errorf("获取科目失败: %v", err)
	}

	if subject == nil {
		logger.WithFields(map[string]interface{}{
			"code": code,
		}).Warn("Subject not found")
		return nil, fmt.Errorf("科目不存在")
	}

	logger.WithFields(map[string]interface{}{
		"code": code,
		"name": subject.Name,
	}).Info("Subject retrieved successfully")

	return subject, nil
}

// GetAllSubjects 获取所有科目（分页）
func (s *SubjectService) GetAllSubjects(req domain.SubjectListRequest) ([]*domain.Subject, int64, error) {
	logger.WithFields(map[string]interface{}{
		"page": req.Page,
		"size": req.Size,
	}).Info("Getting all subjects")

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100 // 限制最大页面大小
	}

	subjects, total, err := s.repo.List(&req)
	if err != nil {
		logger.WithError(err).Error("Failed to get subjects list")
		return nil, 0, fmt.Errorf("获取科目列表失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"count": len(subjects),
		"total": total,
		"page":  req.Page,
		"size":  req.Size,
	}).Info("Subjects list retrieved successfully")

	return subjects, total, nil
}

// UpdateSubject 更新科目信息
func (s *SubjectService) UpdateSubject(id int, req domain.UpdateSubjectRequest) (*domain.Subject, error) {
	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Updating subject")

	// 先获取现有科目
	subject, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).Error("Failed to get subject for update")
		return nil, fmt.Errorf("获取科目失败: %v", err)
	}

	if subject == nil {
		logger.WithFields(map[string]interface{}{
			"subject_id": id,
		}).Warn("Subject not found for update")
		return nil, fmt.Errorf("科目不存在")
	}

	// 如果要更新科目代码，检查新代码是否已存在
	if req.Code != "" && req.Code != subject.Code {
		exists, err := s.repo.ExistsByCode(req.Code)
		if err != nil {
			logger.WithError(err).Error("Failed to check subject code existence")
			return nil, fmt.Errorf("检查科目代码失败: %v", err)
		}

		if exists {
			logger.WithFields(map[string]interface{}{
				"code": req.Code,
			}).Warn("Subject code already exists")
			return nil, fmt.Errorf("科目代码 %s 已存在", req.Code)
		}
	}

	// 更新字段
	if req.Code != "" {
		subject.Code = req.Code
	}
	if req.Name != "" {
		subject.Name = req.Name
	}
	if req.Description != "" {
		subject.Description = req.Description
	}
	if req.Credits > 0 {
		subject.Credits = req.Credits
	}
	if req.Status != "" {
		subject.Status = req.Status
	}

	err = s.repo.Update(subject)
	if err != nil {
		logger.WithError(err).Error("Failed to update subject")
		return nil, fmt.Errorf("更新科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": id,
		"name":       subject.Name,
	}).Info("Subject updated successfully")

	return subject, nil
}

// DeleteSubject 删除科目
func (s *SubjectService) DeleteSubject(id int) error {
	logger.WithFields(map[string]interface{}{
		"subject_id": id,
	}).Info("Deleting subject")

	// 先检查科目是否存在
	subject, err := s.repo.GetByID(id)
	if err != nil {
		logger.WithError(err).Error("Failed to get subject for deletion")
		return fmt.Errorf("获取科目失败: %v", err)
	}

	if subject == nil {
		logger.WithFields(map[string]interface{}{
			"subject_id": id,
		}).Warn("Subject not found for deletion")
		return fmt.Errorf("科目不存在")
	}

	// TODO: 检查是否有相关的成绩记录，如果有则不允许删除
	// 这里可以添加业务逻辑检查

	err = s.repo.Delete(id)
	if err != nil {
		logger.WithError(err).Error("Failed to delete subject")
		return fmt.Errorf("删除科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"subject_id": id,
		"name":       subject.Name,
	}).Info("Subject deleted successfully")

	return nil
}

// GetActiveSubjects 获取活跃科目
func (s *SubjectService) GetActiveSubjects() ([]*domain.Subject, error) {
	logger.Info("Getting active subjects")

	subjects, err := s.repo.GetActiveSubjects()
	if err != nil {
		logger.WithError(err).Error("Failed to get active subjects")
		return nil, fmt.Errorf("获取活跃科目失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"count": len(subjects),
	}).Info("Active subjects retrieved successfully")

	return subjects, nil
}