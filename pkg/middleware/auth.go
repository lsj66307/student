package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"student-management-system/internal/domain"
	"student-management-system/pkg/logger"
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
		logger.Debug("开始JWT认证")

		// 从Authorization header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("缺少Authorization header")
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
			logger.WithError(err).Warn("提取token失败")
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
			logger.WithError(err).Warn("token验证失败")
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

		logger.WithFields(logger.Fields{
			"admin_id": claims.AdminID,
			"account":  claims.Account,
		}).Info("JWT认证成功")

		// 将claims存储到上下文中，供后续处理器使用
		c.Set("claims", claims)
		c.Set("admin_id", claims.AdminID)
		c.Set("account", claims.Account)

		// 继续处理请求
		c.Next()
	}
}

// OptionalJWTAuth 可选的JWT认证中间件（不强制要求认证）
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("开始可选JWT认证")

		// 从Authorization header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Debug("无Authorization header，继续处理")
			// 没有token，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 提取token
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			logger.WithError(err).Debug("提取token失败，继续处理")
			// token格式错误，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 验证token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			logger.WithError(err).Debug("token验证失败，继续处理")
			// token无效，继续处理但不设置用户信息
			c.Next()
			return
		}

		logger.WithFields(logger.Fields{
			"admin_id": claims.AdminID,
			"account":  claims.Account,
		}).Debug("可选JWT认证成功")

		// 将claims存储到上下文中
		c.Set("claims", claims)
		c.Set("admin_id", claims.AdminID)
		c.Set("account", claims.Account)
		c.Set("authenticated", true)

		// 继续处理请求
		c.Next()
	}
}

// AdminRequired 要求管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("检查管理员权限")

		// 检查是否已通过JWT认证
		claims, exists := c.Get("claims")
		if !exists {
			logger.Warn("未通过JWT认证")
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
			logger.Error("认证数据类型错误")
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal error",
				Message: "Invalid authentication data",
			})
			c.Abort()
			return
		}

		// 检查管理员ID是否有效
		if jwtClaims.AdminID <= 0 {
			logger.WithFields(logger.Fields{
				"admin_id": jwtClaims.AdminID,
			}).Warn("无效的管理员ID")
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Forbidden",
				Message: "Invalid admin ID",
			})
			c.Abort()
			return
		}

		// 检查账号是否为空
		if jwtClaims.Account == "" {
			logger.Warn("账号为空")
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "Forbidden",
				Message: "Invalid account",
			})
			c.Abort()
			return
		}

		logger.WithFields(logger.Fields{
			"admin_id": jwtClaims.AdminID,
			"account":  jwtClaims.Account,
		}).Debug("管理员权限验证通过")

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
