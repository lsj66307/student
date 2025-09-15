package service

import (
	"crypto/md5"
	"fmt"

	"student-management-system/internal/domain"
	"student-management-system/internal/repository"

	"github.com/sirupsen/logrus"
)

type AdminService struct {
	adminRepo *repository.AdminRepository
	logger    *logrus.Logger
}

func NewAdminService(adminRepo *repository.AdminRepository, logger *logrus.Logger) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		logger:    logger,
	}
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(req *domain.CreateAdminRequest) (*domain.AdminInfo, error) {
	// 检查账号是否已存在
	existingAdmin, err := s.adminRepo.GetAdminByAccount(req.Account)
	if err == nil && existingAdmin != nil {
		return nil, fmt.Errorf("账号已存在")
	}

	// 密码加密
	hashedPassword := s.hashPassword(req.Password)

	// 创建管理员对象
	admin := &domain.Admin{
		Account:  req.Account,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	// 保存到数据库
	err = s.adminRepo.CreateAdmin(admin)
	if err != nil {
		return nil, fmt.Errorf("创建管理员失败: %v", err)
	}

	// 返回管理员信息（不包含密码）
	adminInfo := &domain.AdminInfo{
		ID:      admin.ID,
		Account: admin.Account,
		Name:    admin.Name,
		Phone:   admin.Phone,
		Email:   admin.Email,
	}

	s.logger.WithField("admin_id", admin.ID).Info("Admin created successfully")
	return adminInfo, nil
}

// GetAdminByID 根据ID获取管理员信息
func (s *AdminService) GetAdminByID(id int) (*domain.AdminInfo, error) {
	admin, err := s.adminRepo.GetAdminByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取管理员信息失败: %v", err)
	}

	adminInfo := &domain.AdminInfo{
		ID:      admin.ID,
		Account: admin.Account,
		Name:    admin.Name,
		Phone:   admin.Phone,
		Email:   admin.Email,
	}

	return adminInfo, nil
}

// UpdateAdmin 更新管理员信息
func (s *AdminService) UpdateAdmin(id int, req *domain.UpdateAdminRequest) (*domain.AdminInfo, error) {
	// 获取现有管理员信息
	admin, err := s.adminRepo.GetAdminByID(id)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在")
	}

	// 如果要更新账号，检查新账号是否已被其他管理员使用
	if req.Account != "" && req.Account != admin.Account {
		existingAdmin, err := s.adminRepo.GetAdminByAccount(req.Account)
		if err == nil && existingAdmin != nil && existingAdmin.ID != id {
			return nil, fmt.Errorf("账号已被其他管理员使用")
		}
		admin.Account = req.Account
	}

	// 更新其他字段
	if req.Name != "" {
		admin.Name = req.Name
	}
	if req.Phone != "" {
		admin.Phone = req.Phone
	}
	if req.Email != "" {
		admin.Email = req.Email
	}

	// 如果要更新密码
	if req.Password != "" {
		admin.Password = s.hashPassword(req.Password)
	}

	// 保存更新
	err = s.adminRepo.UpdateAdmin(admin)
	if err != nil {
		return nil, fmt.Errorf("更新管理员信息失败: %v", err)
	}

	// 返回更新后的管理员信息
	adminInfo := &domain.AdminInfo{
		ID:      admin.ID,
		Account: admin.Account,
		Name:    admin.Name,
		Phone:   admin.Phone,
		Email:   admin.Email,
	}

	s.logger.WithField("admin_id", id).Info("Admin updated successfully")
	return adminInfo, nil
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(id int) error {
	// 检查管理员是否存在
	_, err := s.adminRepo.GetAdminByID(id)
	if err != nil {
		return fmt.Errorf("管理员不存在")
	}

	// 删除管理员
	err = s.adminRepo.DeleteAdmin(id)
	if err != nil {
		return fmt.Errorf("删除管理员失败: %v", err)
	}

	s.logger.WithField("admin_id", id).Info("Admin deleted successfully")
	return nil
}

// ListAdmins 获取管理员列表
func (s *AdminService) ListAdmins(req *domain.AdminListRequest) (*domain.AdminListResponse, error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 获取管理员列表
	admins, total, err := s.adminRepo.ListAdmins(req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取管理员列表失败: %v", err)
	}

	// 转换为响应格式（不包含密码）
	var adminInfos []domain.AdminInfo
	for _, admin := range admins {
		adminInfo := domain.AdminInfo{
			ID:      admin.ID,
			Account: admin.Account,
			Name:    admin.Name,
			Phone:   admin.Phone,
			Email:   admin.Email,
		}
		adminInfos = append(adminInfos, adminInfo)
	}

	response := &domain.AdminListResponse{
		Data:  adminInfos,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}

	return response, nil
}

// ValidateAdmin 验证管理员账号密码（用于登录）
func (s *AdminService) ValidateAdmin(account, password string) (*domain.Admin, error) {
	// 根据账号获取管理员
	admin, err := s.adminRepo.GetAdminByAccount(account)
	if err != nil {
		return nil, fmt.Errorf("账号或密码错误")
	}

	// 验证密码
	if s.hashPassword(password) != admin.Password {
		return nil, fmt.Errorf("账号或密码错误")
	}

	return admin, nil
}

// hashPassword MD5密码加密
func (s *AdminService) hashPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetAdminInfo 获取管理员信息（不包含密码）
func (s *AdminService) GetAdminInfo(admin *domain.Admin) *domain.AdminInfo {
	return &domain.AdminInfo{
		ID:      admin.ID,
		Account: admin.Account,
		Name:    admin.Name,
		Phone:   admin.Phone,
		Email:   admin.Email,
	}
}
