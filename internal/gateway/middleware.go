package gateway

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"student-management-system/internal/handler"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全中间件
func (g *Gateway) SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// 移除服务器标识
		c.Header("Server", "")

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func (g *Gateway) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允许的域名列表（生产环境应该配置具体域名）
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}

		// 检查是否为允许的域名
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-Response-Time")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func (g *Gateway) LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RateLimitMiddleware 限流中间件
func (g *Gateway) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单的限流实现，实际项目中可以使用Redis等
		// 这里可以根据IP、用户ID等进行限流
		clientIP := c.ClientIP()

		// 检查是否超过限流阈值
		if g.isRateLimited(clientIP) {
			c.JSON(http.StatusTooManyRequests, handler.Response{
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthMiddleware 认证中间件
func (g *Gateway) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, handler.Response{
				Code:    http.StatusUnauthorized,
				Message: "Missing authorization header",
			})
			c.Abort()
			return
		}

		// 检查Bearer token格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, handler.Response{
				Code:    http.StatusUnauthorized,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// 提取token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, handler.Response{
				Code:    http.StatusUnauthorized,
				Message: "Empty token",
			})
			c.Abort()
			return
		}

		// 验证token（这里应该调用认证服务）
		if !g.validateToken(token) {
			c.JSON(http.StatusUnauthorized, handler.Response{
				Code:    http.StatusUnauthorized,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文（从token中解析）
		userInfo := g.extractUserFromToken(token)
		c.Set("user", userInfo)
		c.Set("user_id", userInfo["user_id"])

		c.Next()
	}
}

// isRateLimited 检查是否被限流
func (g *Gateway) isRateLimited(clientIP string) bool {
	// 简单的内存限流实现
	// 实际项目中应该使用Redis等分布式缓存
	return false
}

// validateToken 验证token
func (g *Gateway) validateToken(token string) bool {
	// 这里应该调用认证服务验证token
	// 或者使用JWT库验证token

	// 简单的示例实现
	if token == "invalid" || token == "expired" {
		return false
	}

	// 实际项目中应该验证JWT签名、过期时间等
	return len(token) > 10
}

// extractUserFromToken 从token中提取用户信息
func (g *Gateway) extractUserFromToken(token string) map[string]interface{} {
	// 这里应该解析JWT token获取用户信息
	// 简单的示例实现
	return map[string]interface{}{
		"user_id":  "12345",
		"username": "test_user",
		"role":     "user",
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
}
