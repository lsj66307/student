package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"student-management-system/pkg/logger"
)

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// 通用错误
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"

	// 业务错误
	ErrCodeStudentNotFound    ErrorCode = "STUDENT_NOT_FOUND"
	ErrCodeTeacherNotFound    ErrorCode = "TEACHER_NOT_FOUND"
	ErrCodeDuplicateStudent   ErrorCode = "DUPLICATE_STUDENT"
	ErrCodeDuplicateTeacher   ErrorCode = "DUPLICATE_TEACHER"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidToken       ErrorCode = "INVALID_TOKEN"

	// 数据库错误
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrCodeConnectionError ErrorCode = "CONNECTION_ERROR"
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ErrorResponse API错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
}

// Newf 创建格式化的应用错误
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		HTTPStatus: getHTTPStatus(code),
	}
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
		Cause:      err,
	}
}

// Wrapf 包装现有错误并格式化消息
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		HTTPStatus: getHTTPStatus(code),
		Cause:      err,
	}
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithDetailsf 添加格式化的详细信息
func (e *AppError) WithDetailsf(format string, args ...interface{}) *AppError {
	e.Details = fmt.Sprintf(format, args...)
	return e
}

// getHTTPStatus 根据错误代码获取HTTP状态码
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeInvalidToken:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeStudentNotFound, ErrCodeTeacherNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeDuplicateStudent, ErrCodeDuplicateTeacher:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// HandleError 统一错误处理中间件
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 记录错误日志
	logger.WithError(err).WithFields(map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Error("Request failed")

	// 检查是否是应用错误
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Error:   string(appErr.Code),
			Message: appErr.Message,
			Code:    string(appErr.Code),
			Details: appErr.Details,
		})
		return
	}

	// 处理其他类型的错误
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   string(ErrCodeInternal),
		Message: "Internal server error",
		Code:    string(ErrCodeInternal),
	})
}

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithFields(map[string]interface{}{
			"panic":  recovered,
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   string(ErrCodeInternal),
			Message: "Internal server error",
			Code:    string(ErrCodeInternal),
		})
	})
}

// LoggingMiddleware 请求日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(map[string]interface{}{
			"status":     param.StatusCode,
			"method":     param.Method,
			"path":       param.Path,
			"ip":         param.ClientIP,
			"user_agent": param.Request.UserAgent(),
			"latency":    param.Latency,
		}).Info("Request processed")
		return ""
	})
}

// 预定义的常用错误
var (
	ErrInternal       = New(ErrCodeInternal, "Internal server error")
	ErrInvalidRequest = New(ErrCodeInvalidRequest, "Invalid request")
	ErrUnauthorized   = New(ErrCodeUnauthorized, "Unauthorized")
	ErrForbidden      = New(ErrCodeForbidden, "Forbidden")
	ErrNotFound       = New(ErrCodeNotFound, "Resource not found")
	ErrConflict       = New(ErrCodeConflict, "Resource conflict")
	ErrValidation     = New(ErrCodeValidation, "Validation failed")

	ErrStudentNotFound    = New(ErrCodeStudentNotFound, "Student not found")
	ErrTeacherNotFound    = New(ErrCodeTeacherNotFound, "Teacher not found")
	ErrDuplicateStudent   = New(ErrCodeDuplicateStudent, "Student already exists")
	ErrDuplicateTeacher   = New(ErrCodeDuplicateTeacher, "Teacher already exists")
	ErrInvalidCredentials = New(ErrCodeInvalidCredentials, "Invalid username or password")
	ErrTokenExpired       = New(ErrCodeTokenExpired, "Token has expired")
	ErrInvalidToken       = New(ErrCodeInvalidToken, "Invalid token")

	ErrDatabaseError   = New(ErrCodeDatabaseError, "Database operation failed")
	ErrConnectionError = New(ErrCodeConnectionError, "Database connection failed")
)
