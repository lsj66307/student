package validator

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator 验证器接口
type Validator interface {
	ValidateStruct(s interface{}) error
	ValidateVar(field interface{}, tag string) error
	Middleware() gin.HandlerFunc
}

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator 创建新的验证器
func NewValidator() *CustomValidator {
	v := validator.New()

	// 注册自定义验证规则
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("idcard", validateIDCard)
	v.RegisterValidation("studentid", validateStudentID)
	v.RegisterValidation("nohtml", validateNoHTML)
	v.RegisterValidation("nosql", validateNoSQL)
	v.RegisterValidation("safename", validateSafeName)

	return &CustomValidator{
		validator: v,
	}
}

// ValidateStruct 验证结构体
func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	return cv.validator.Struct(s)
}

// ValidateVar 验证单个变量
func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return cv.validator.Var(field, tag)
}

// Middleware 返回验证中间件
func (cv *CustomValidator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证请求参数
		if err := cv.validateRequest(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"code":    "VALIDATION_ERROR",
				"details": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// validateRequest 验证请求
func (cv *CustomValidator) validateRequest(c *gin.Context) error {
	// 清理查询参数
	for key, values := range c.Request.URL.Query() {
		for i, value := range values {
			c.Request.URL.Query()[key][i] = SanitizeInput(value)
		}
	}

	// 验证路径参数
	for _, param := range c.Params {
		if err := cv.validatePathParam(param.Key, param.Value); err != nil {
			return fmt.Errorf("invalid path parameter %s: %w", param.Key, err)
		}
	}

	return nil
}

// validatePathParam 验证路径参数
func (cv *CustomValidator) validatePathParam(key, value string) error {
	switch key {
	case "id":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("must be a valid integer")
		}
	case "student_id":
		if err := cv.ValidateVar(value, "studentid"); err != nil {
			return err
		}
	}
	return nil
}

// SanitizeInput 清理输入数据
func SanitizeInput(input string) string {
	// 移除前后空格
	input = strings.TrimSpace(input)

	// 移除潜在的HTML标签
	input = removeHTMLTags(input)

	// 移除潜在的SQL注入字符
	input = removeSQLInjection(input)

	// 移除控制字符
	input = removeControlChars(input)

	return input
}

// removeHTMLTags 移除HTML标签
func removeHTMLTags(input string) string {
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	return htmlTagRegex.ReplaceAllString(input, "")
}

// removeSQLInjection 移除SQL注入字符
func removeSQLInjection(input string) string {
	// 移除常见的SQL注入模式
	sqlPatterns := []string{
		`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`,
		`(?i)(script|javascript|vbscript|onload|onerror|onclick)`,
		`['";\-\-]`,
	}

	for _, pattern := range sqlPatterns {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, "")
	}

	return input
}

// removeControlChars 移除控制字符
func removeControlChars(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)
}

// 自定义验证函数

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	// 简化的身份证验证（18位数字，最后一位可能是X）
	idCardRegex := regexp.MustCompile(`^\d{17}[\dXx]$`)
	return idCardRegex.MatchString(idCard)
}

// validateStudentID 验证学号
func validateStudentID(fl validator.FieldLevel) bool {
	studentID := fl.Field().String()
	// 学号格式：年份(4位) + 专业代码(2位) + 序号(4位)
	studentIDRegex := regexp.MustCompile(`^20\d{2}\d{2}\d{4}$`)
	return studentIDRegex.MatchString(studentID)
}

// validateNoHTML 验证不包含HTML标签
func validateNoHTML(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	return !htmlRegex.MatchString(value)
}

// validateNoSQL 验证不包含SQL注入字符
func validateNoSQL(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	sqlRegex := regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|['";\-\-])`)
	return !sqlRegex.MatchString(value)
}

// validateSafeName 验证安全的姓名格式
func validateSafeName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	// 只允许中文、英文字母和空格
	nameRegex := regexp.MustCompile(`^[\p{Han}a-zA-Z\s]+$`)
	return nameRegex.MatchString(name) && len(name) >= 2 && len(name) <= 50
}

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMsg := getErrorMessage(fieldError)
			errors = append(errors, ValidationError{
				Field:   fieldError.Field(),
				Tag:     fieldError.Tag(),
				Value:   fmt.Sprintf("%v", fieldError.Value()),
				Message: errorMsg,
			})
		}
	}

	return errors
}

// getErrorMessage 获取错误消息
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "此字段为必填项"
	case "email":
		return "请输入有效的邮箱地址"
	case "phone":
		return "请输入有效的手机号码"
	case "idcard":
		return "请输入有效的身份证号码"
	case "studentid":
		return "请输入有效的学号"
	case "min":
		return fmt.Sprintf("最小长度为 %s", fe.Param())
	case "max":
		return fmt.Sprintf("最大长度为 %s", fe.Param())
	case "len":
		return fmt.Sprintf("长度必须为 %s", fe.Param())
	case "nohtml":
		return "不允许包含HTML标签"
	case "nosql":
		return "输入包含非法字符"
	case "safename":
		return "姓名只能包含中文、英文字母和空格，长度2-50字符"
	default:
		return "输入格式不正确"
	}
}

// BindAndValidate 绑定并验证请求数据
func BindAndValidate(c *gin.Context, obj interface{}) error {
	// 绑定请求数据
	if err := c.ShouldBindJSON(obj); err != nil {
		return fmt.Errorf("数据绑定失败: %w", err)
	}

	// 清理输入数据
	SanitizeStruct(obj)

	// 验证数据
	v := NewValidator()
	if err := v.ValidateStruct(obj); err != nil {
		return err
	}

	return nil
}

// SanitizeStruct 清理结构体中的字符串字段
func SanitizeStruct(obj interface{}) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			cleanValue := SanitizeInput(field.String())
			field.SetString(cleanValue)
		}
	}
}
