package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"student-management-system/internal/domain"
	"student-management-system/pkg/utils"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 提取token
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			var message string
			switch err {
			case domain.ErrTokenExpired:
				message = "Token has expired"
			case domain.ErrInvalidToken:
				message = "Invalid token"
			default:
				message = err.Error()
			}

			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: message,
			})
			c.Abort()
			return
		}

		// 将claims存储到上下文中，供后续处理器使用
		c.Set("claims", claims)
		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)

		// 继续处理请求
		c.Next()
	}
}

// OptionalJWTAuth 可选的JWT认证中间件（不强制要求认证）
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有token，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 提取token
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// token格式错误，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 验证token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			// token无效，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 将claims存储到上下文中
		c.Set("claims", claims)
		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)
		c.Set("authenticated", true)

		// 继续处理请求
		c.Next()
	}
}

// AdminRequired 要求管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否已通过JWT认证
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		// 验证claims类型
		jwtClaims, ok := claims.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal error",
				Message: "Invalid authentication data",
			})
			c.Abort()
			return
		}

		// 这里可以添加更多的权限检查逻辑
		// 例如检查用户角色、权限等
		_ = jwtClaims // 暂时不做额外检查，所有通过JWT认证的用户都被视为管理员

		// 继续处理请求
		c.Next()
	}
}

// GetCurrentAdmin 从上下文中获取当前管理员信息的辅助函数
func GetCurrentAdmin(c *gin.Context) (*domain.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*domain.JWTClaims)
	if !ok {
		return nil, false
	}

	return jwtClaims, true
}

// IsAuthenticated 检查用户是否已认证的辅助函数
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("claims")
	return exists
}