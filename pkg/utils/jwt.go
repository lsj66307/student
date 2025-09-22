package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/errors"
	"student-management-system/pkg/logger"
)

// JWTSecret JWT密钥
var JWTSecret = []byte("your-secret-key-change-this-in-production")

// GenerateToken 生成JWT token并存储到Redis
func GenerateToken(adminID int, username string, expiresIn int64) (string, time.Time, error) {
	logger.WithFields(logger.Fields{
		"admin_id":   adminID,
		"username":   username,
		"expires_in": expiresIn,
	}).Debug("开始生成JWT token")

	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	// 创建JWT声明
	claims := domain.JWTClaims{
		AdminID: adminID,
		Account: username,
		Exp:     expiresAt.Unix(),
		Iat:     now.Unix(),
	}

	// 创建header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 编码header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		logger.WithError(err).Error("编码JWT header失败")
		return "", time.Time{}, err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	// 编码payload
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		logger.WithError(err).Error("编码JWT payload失败")
		return "", time.Time{}, err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := createSignature(message, JWTSecret)

	// 组合token
	token := message + "." + signature

	// 将token存储到Redis
	if repository.RedisClient != nil {
		ctx := context.Background()
		tokenKey := fmt.Sprintf("jwt_token:%d", adminID)

		// 存储token到Redis，设置过期时间
		err = repository.RedisClient.Set(ctx, tokenKey, token, time.Duration(expiresIn)*time.Second).Err()
		if err != nil {
			logger.WithError(err).Warn("存储token到Redis失败")
		} else {
			logger.WithFields(logger.Fields{
				"admin_id":  adminID,
				"token_key": tokenKey,
			}).Debug("Token已存储到Redis")
		}
	}

	logger.WithFields(logger.Fields{
		"admin_id":   adminID,
		"username":   username,
		"expires_at": expiresAt,
	}).Info("JWT token生成成功")

	return token, expiresAt, nil
}

// ValidateToken 验证JWT token并检查Redis存储
func ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	logger.Debug("开始验证JWT token")

	// 分割token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		logger.Warn("JWT token格式错误：分段数量不正确")
		return nil, errors.ErrInvalidToken
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	// 验证签名
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := createSignature(message, JWTSecret)
	if signatureEncoded != expectedSignature {
		logger.Warn("JWT token签名验证失败")
		return nil, errors.ErrInvalidToken
	}

	// 解码payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		logger.WithError(err).Warn("JWT token payload解码失败")
		return nil, errors.ErrInvalidToken
	}

	// 解析claims
	var claims domain.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		logger.WithError(err).Warn("JWT token claims解析失败")
		return nil, errors.ErrInvalidToken
	}

	// 验证过期时间
	if err := claims.Valid(); err != nil {
		logger.WithError(err).Warn("JWT token验证失败")
		return nil, err
	}

	// 检查Redis中的token是否存在且匹配
	if repository.RedisClient != nil {
		ctx := context.Background()
		tokenKey := fmt.Sprintf("jwt_token:%d", claims.AdminID)

		storedToken, err := repository.RedisClient.Get(ctx, tokenKey).Result()
		if err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"admin_id":  claims.AdminID,
				"token_key": tokenKey,
			}).Warn("从Redis获取token失败或token不存在")
			return nil, errors.ErrTokenExpired
		}

		if storedToken != tokenString {
			logger.WithFields(logger.Fields{
				"admin_id": claims.AdminID,
			}).Warn("Redis中的token与提供的token不匹配")
			return nil, errors.ErrInvalidToken
		}

		logger.WithFields(logger.Fields{
			"admin_id": claims.AdminID,
		}).Debug("Redis token验证成功")
	}

	logger.WithFields(logger.Fields{
		"admin_id": claims.AdminID,
		"account":  claims.Account,
	}).Debug("JWT token验证成功")

	return &claims, nil
}

// createSignature 创建HMAC-SHA256签名
func createSignature(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature)
}

// ExtractTokenFromHeader 从Authorization header中提取token
func ExtractTokenFromHeader(authHeader string) (string, error) {
	logger.Debug("从Authorization header提取token")

	if authHeader == "" {
		logger.Warn("Authorization header为空")
		return "", fmt.Errorf("authorization header is required")
	}

	// 检查Bearer前缀
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		logger.Warn("Authorization header格式错误")
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	logger.Debug("成功提取token")
	return parts[1], nil
}

// InvalidateToken 使token失效（从Redis中删除）
func InvalidateToken(adminID int) error {
	if repository.RedisClient == nil {
		logger.Warn("Redis客户端未初始化，无法使token失效")
		return nil
	}

	ctx := context.Background()
	tokenKey := fmt.Sprintf("jwt_token:%d", adminID)

	err := repository.RedisClient.Del(ctx, tokenKey).Err()
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"admin_id":  adminID,
			"token_key": tokenKey,
		}).Error("从Redis删除token失败")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"admin_id":  adminID,
		"token_key": tokenKey,
	}).Info("Token已从Redis中删除")

	return nil
}
