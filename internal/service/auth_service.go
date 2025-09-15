package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"student-management-system/internal/config"
	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/utils"
)

// AuthService 认证服务
type AuthService struct {
	config *config.Config
}

// NewAuthService 创建认证服务实例
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
	}
}

// Login 管理员登录
func (s *AuthService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("login_lock:%s", req.Username)
	failCountKey := fmt.Sprintf("login_fail_count:%s", req.Username)

	// 检查是否被锁定
	locked, err := repository.RedisClient.Exists(ctx, lockKey).Result()
	if err != nil {
		return nil, fmt.Errorf("检查锁定状态失败: %w", err)
	}
	if locked > 0 {
		ttl, _ := repository.RedisClient.TTL(ctx, lockKey).Result()
		return nil, fmt.Errorf("账户已被锁定，请在 %v 后重试", ttl)
	}

	// 验证用户名和密码
	if !s.validateCredentials(req.Username, req.Password) {
		// 登录失败，增加失败次数
		failCount, err := repository.RedisClient.Incr(ctx, failCountKey).Result()
		if err != nil {
			return nil, fmt.Errorf("记录登录失败次数失败: %w", err)
		}

		// 设置失败次数的过期时间（15分钟）
		repository.RedisClient.Expire(ctx, failCountKey, 15*time.Minute)

		// 如果失败次数达到5次，锁定账户5分钟
		if failCount >= 5 {
			err = repository.RedisClient.Set(ctx, lockKey, "locked", 5*time.Minute).Err()
			if err != nil {
				return nil, fmt.Errorf("锁定账户失败: %w", err)
			}
			// 清除失败次数
			repository.RedisClient.Del(ctx, failCountKey)
			return nil, fmt.Errorf("登录失败次数过多，账户已被锁定5分钟")
		}

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
	token, expiresAt, err := utils.GenerateToken(1, req.Username, int64(expiresIn.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 构造响应
	response := &domain.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Admin: domain.AdminInfo{
			ID:       1,
			Username: req.Username,
		},
	}

	return response, nil
}

// ValidateToken 验证token
func (s *AuthService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	return utils.ValidateToken(tokenString)
}

// validateCredentials 验证用户凭据（硬编码默认账号）
func (s *AuthService) validateCredentials(username, password string) bool {
	// 默认管理员账号
	defaultUsername := "admin"
	defaultPassword := "123456"

	// 简单的用户名密码验证
	if username != defaultUsername {
		return false
	}

	// 对密码进行简单的哈希验证（实际项目中应该使用bcrypt等安全的哈希算法）
	return s.hashPassword(password) == s.hashPassword(defaultPassword)
}

// hashPassword 简单的密码哈希（仅用于演示，生产环境应使用bcrypt）
func (s *AuthService) hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password + "salt")) // 添加盐值
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetAdminInfo 根据token获取管理员信息
func (s *AuthService) GetAdminInfo(claims *domain.JWTClaims) *domain.AdminInfo {
	return &domain.AdminInfo{
		ID:       claims.AdminID,
		Username: claims.Username,
	}
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(claims *domain.JWTClaims) (*domain.LoginResponse, error) {
	// 检查token是否即将过期（剩余时间少于1小时）
	if time.Now().Unix() > claims.Exp-3600 {
		// 生成新的token
		expiresIn := s.config.JWT.ExpiresIn
		if expiresIn == 0 {
			expiresIn = 24 * time.Hour
		}

		token, expiresAt, err := utils.GenerateToken(claims.AdminID, claims.Username, int64(expiresIn.Seconds()))
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		return &domain.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			Admin: domain.AdminInfo{
				ID:       claims.AdminID,
				Username: claims.Username,
			},
		}, nil
	}

	return nil, fmt.Errorf("token does not need refresh yet")
}
