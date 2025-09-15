package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"student-management-system/internal/domain"
)

// JWTSecret JWT密钥
var JWTSecret = []byte("your-secret-key-change-this-in-production")

// GenerateToken 生成JWT token
func GenerateToken(adminID int, username string, expiresIn int64) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	// 创建JWT声明
	claims := domain.JWTClaims{
		AdminID:  adminID,
		Username: username,
		Exp:      expiresAt.Unix(),
		Iat:      now.Unix(),
	}

	// 创建header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 编码header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", time.Time{}, err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	// 编码payload
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := createSignature(message, JWTSecret)

	// 组合token
	token := message + "." + signature

	return token, expiresAt, nil
}

// ValidateToken 验证JWT token
func ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	// 分割token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, domain.ErrInvalidToken
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signatureEncoded := parts[2]

	// 验证签名
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := createSignature(message, JWTSecret)
	if signatureEncoded != expectedSignature {
		return nil, domain.ErrInvalidToken
	}

	// 解码payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// 解析claims
	var claims domain.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, domain.ErrInvalidToken
	}

	// 验证过期时间
	if err := claims.Valid(); err != nil {
		return nil, err
	}

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
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// 检查Bearer前缀
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
