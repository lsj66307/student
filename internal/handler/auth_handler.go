package handler

import (
	"net/http"

	"student-management-system/internal/domain"
	"student-management-system/internal/service"
	"student-management-system/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 使用用户名和密码进行管理员登录，返回JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body domain.LoginRequest true "登录信息"
// @Success 200 {object} domain.LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户名或密码错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// 执行登录
	response, err := h.authService.Login(&req)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication failed",
				Message: "Invalid username or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile 获取当前管理员信息
// @Summary 获取当前管理员信息
// @Description 根据JWT token获取当前登录的管理员信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.AdminInfo "管理员信息"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取JWT claims（由中间件设置）
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid or missing token",
		})
		return
	}

	jwtClaims, ok := claims.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal error",
			Message: "Invalid token claims",
		})
		return
	}

	adminInfo := h.authService.GetAdminInfo(jwtClaims)
	c.JSON(http.StatusOK, adminInfo)
}

// RefreshToken 刷新JWT token
// @Summary 刷新JWT token
// @Description 使用当前有效的JWT token获取新的token
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.LoginResponse "新的token信息"
// @Failure 400 {object} ErrorResponse "token不需要刷新"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从上下文中获取JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid or missing token",
		})
		return
	}

	jwtClaims, ok := claims.(*domain.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal error",
			Message: "Invalid token claims",
		})
		return
	}

	// 刷新token
	response, err := h.authService.RefreshToken(jwtClaims)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Refresh failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateToken 验证token（用于其他服务调用）
// @Summary 验证JWT token
// @Description 验证JWT token的有效性
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} domain.AdminInfo "token有效，返回管理员信息"
// @Failure 401 {object} ErrorResponse "token无效"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/validate [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// 从header中提取token
	authHeader := c.GetHeader("Authorization")
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid authorization header",
			Message: err.Error(),
		})
		return
	}

	// 验证token
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid token",
			Message: err.Error(),
		})
		return
	}

	adminInfo := h.authService.GetAdminInfo(claims)
	c.JSON(http.StatusOK, adminInfo)
}
