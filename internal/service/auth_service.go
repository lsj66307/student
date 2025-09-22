package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"student-management-system/internal/config"
	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/errors"
	"student-management-system/pkg/logger"
	"student-management-system/pkg/utils"
)

// AuthService 认证服务
type AuthService struct {
	config    *config.Config
	adminRepo *repository.AdminRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(cfg *config.Config, adminRepo *repository.AdminRepository) *AuthService {
	return &AuthService{
		config:    cfg,
		adminRepo: adminRepo,
	}
}

// Login 管理员登录
func (s *AuthService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("login_lock:%s", req.Account)
	failCountKey := fmt.Sprintf("login_fail_count:%s", req.Account)

	logger.WithFields(map[string]interface{}{
		"account": req.Account,
	}).Info("Admin login attempt")

	// 检查是否被锁定
	locked, err := repository.RedisClient.Exists(ctx, lockKey).Result()
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"account": req.Account,
		}).Error("Failed to check lock status")
		return nil, fmt.Errorf("检查锁定状态失败: %w", err)
	}
	if locked > 0 {
		ttl, _ := repository.RedisClient.TTL(ctx, lockKey).Result()
		logger.WithFields(map[string]interface{}{
			"account": req.Account,
			"ttl":     ttl,
		}).Warn("Account is locked")
		return nil, fmt.Errorf("账户已被锁定，请在 %v 后重试", ttl)
	}

	// 验证用户名和密码
	admin, err := s.validateCredentials(req.Account, req.Password)
	if err != nil {
		// 登录失败，增加失败次数
		failCount, err := repository.RedisClient.Incr(ctx, failCountKey).Result()
		if err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"account": req.Account,
			}).Error("Failed to record login failure count")
			return nil, fmt.Errorf("记录登录失败次数失败: %w", err)
		}

		// 设置失败次数的过期时间（15分钟）
		repository.RedisClient.Expire(ctx, failCountKey, 15*time.Minute)

		// 如果失败次数达到5次，锁定账户5分钟
		if failCount >= 5 {
			err = repository.RedisClient.Set(ctx, lockKey, "locked", 5*time.Minute).Err()
			if err != nil {
				logger.WithError(err).WithFields(map[string]interface{}{
					"account": req.Account,
				}).Error("Failed to lock account")
				return nil, fmt.Errorf("锁定账户失败: %w", err)
			}
			// 清除失败次数
			repository.RedisClient.Del(ctx, failCountKey)
			logger.WithFields(map[string]interface{}{
				"account":    req.Account,
				"fail_count": failCount,
			}).Warn("Account locked due to too many failed attempts")
			return nil, fmt.Errorf("登录失败次数过多，账户已被锁定5分钟")
		}

		logger.WithFields(map[string]interface{}{
			"account":    req.Account,
			"fail_count": failCount,
		}).Warn("Login failed - invalid credentials")
		return nil, fmt.Errorf("用户名或密码错误，还可尝试 %d 次", 5-failCount)
	}

	// 登录成功，清除失败次数
	repository.RedisClient.Del(ctx, failCountKey)

	// 获取JWT过期时间配置
	expiresIn := s.config.JWT.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 24 * time.Hour // 默认24小时
	}

	// 生成JWT token
	token, expiresAt, err := utils.GenerateToken(admin.ID, admin.Account, int64(expiresIn.Seconds()))
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"account": req.Account,
		}).Error("Failed to generate token")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 构造响应
	response := &domain.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Admin: domain.AdminInfo{
			ID:      admin.ID,
			Account: admin.Account,
			Name:    admin.Name,
			Phone:   admin.Phone,
			Email:   admin.Email,
		},
	}

	logger.WithFields(map[string]interface{}{
		"account":    req.Account,
		"expires_at": expiresAt,
	}).Info("Admin login successful")

	return response, nil
}

// ValidateToken 验证token
func (s *AuthService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	return utils.ValidateToken(tokenString)
}

// validateCredentials 验证用户凭据
func (s *AuthService) validateCredentials(account, password string) (*domain.Admin, error) {
	// 从数据库中查询管理员信息
	admin, err := s.adminRepo.GetAdminByAccount(account)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"account": account,
		}).Error("Failed to get admin by account")
		return nil, errors.ErrInvalidCredentials
	}

	// 验证密码
	if s.hashPassword(password) != admin.Password {
		logger.WithFields(map[string]interface{}{
			"account": account,
		}).Warn("Password verification failed")
		return nil, errors.ErrInvalidCredentials
	}

	return admin, nil
}

// hashPassword 简单的密码哈希（使用MD5）
func (s *AuthService) hashPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetAdminInfo 根据token获取管理员信息
func (s *AuthService) GetAdminInfo(claims *domain.JWTClaims) *domain.AdminInfo {
	// 从数据库中获取完整的管理员信息
	admin, err := s.adminRepo.GetAdminByID(claims.AdminID)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"admin_id": claims.AdminID,
		}).Error("Failed to get admin info")
		// 如果获取失败，返回基本信息
		return &domain.AdminInfo{
			ID:      claims.AdminID,
			Account: claims.Account,
		}
	}

	return &domain.AdminInfo{
		ID:      admin.ID,
		Account: admin.Account,
		Name:    admin.Name,
		Phone:   admin.Phone,
		Email:   admin.Email,
	}
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(claims *domain.JWTClaims) (*domain.LoginResponse, error) {
	// 检查token是否即将过期（剩余时间少于1小时）
	if time.Now().Unix() > claims.Exp-3600 {
		// 先使旧token失效
		if err := utils.InvalidateToken(claims.AdminID); err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"admin_id": claims.AdminID,
				"account":  claims.Account,
			}).Warn("Failed to invalidate old token during refresh")
		}

		// 生成新的token
		expiresIn := s.config.JWT.ExpiresIn
		if expiresIn == 0 {
			expiresIn = 12 * time.Hour // 默认12小时
		}

		token, expiresAt, err := utils.GenerateToken(claims.AdminID, claims.Account, int64(expiresIn.Seconds()))
		if err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"admin_id": claims.AdminID,
				"account":  claims.Account,
			}).Error("Failed to refresh token")
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		logger.WithFields(map[string]interface{}{
			"admin_id":   claims.AdminID,
			"account":    claims.Account,
			"expires_at": expiresAt,
		}).Info("Token refreshed successfully")

		// 获取完整的管理员信息
		adminInfo := s.GetAdminInfo(claims)

		return &domain.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			Admin:     *adminInfo,
		}, nil
	}

	return nil, fmt.Errorf("token does not need refresh yet")
}

// Logout 用户登出，使token失效
func (s *AuthService) Logout(adminID int) error {
	logger.WithFields(map[string]interface{}{
		"admin_id": adminID,
	}).Info("User logout initiated")

	// 使token失效
	err := utils.InvalidateToken(adminID)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"admin_id": adminID,
		}).Error("Failed to invalidate token during logout")
		return fmt.Errorf("failed to logout: %w", err)
	}

	logger.WithFields(map[string]interface{}{
		"admin_id": adminID,
	}).Info("User logout successful")

	return nil
}
