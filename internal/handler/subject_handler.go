package handler

import (
	"net/http"
	"strconv"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"
	"student-management-system/pkg/validator"

	"github.com/gin-gonic/gin"
)

// SubjectHandler 科目处理器
type SubjectHandler struct {
	subjectService *service.SubjectService
	validator      *validator.CustomValidator
}

// NewSubjectHandler 创建新的科目处理器
func NewSubjectHandler(subjectService *service.SubjectService, validator *validator.CustomValidator) *SubjectHandler {
	return &SubjectHandler{
		subjectService: subjectService,
		validator:      validator,
	}
}

// CreateSubject 创建科目
// @Summary 创建新科目
// @Description 创建一个新的科目记录
// @Tags subjects
// @Accept json
// @Produce json
// @Param subject body domain.CreateSubjectRequest true "科目信息"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req domain.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证和清理输入数据
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	subject, err := h.subjectService.CreateSubject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create subject",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "科目创建成功",
		Data:    subject,
	})
}

// GetSubject 获取单个科目
// @Summary 获取科目详情
// @Description 根据ID获取科目详细信息
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/{id} [get]
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "科目ID格式错误",
		})
		return
	}

	subject, err := h.subjectService.GetSubjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Subject not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取科目成功",
		Data:    subject,
	})
}

// GetSubjects 获取科目列表
// @Summary 获取科目列表
// @Description 获取科目列表，支持分页和筛选
// @Tags subjects
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "科目名称"
// @Param code query string false "科目代码"
// @Param status query string false "状态"
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects [get]
func (h *SubjectHandler) GetSubjects(c *gin.Context) {
	var req domain.SubjectListRequest

	// 解析查询参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}
	if sizeStr := c.Query("size"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil {
			req.Size = size
		}
	}
	req.Name = c.Query("name")
	req.Code = c.Query("code")
	req.Status = c.Query("status")

	// 验证请求参数
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "查询参数验证失败: " + err.Error(),
		})
		return
	}

	subjects, total, err := h.subjectService.GetAllSubjects(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get subjects",
			Message: err.Error(),
		})
		return
	}

	// 转换指针切片为值切片
	subjectList := make([]domain.Subject, len(subjects))
	for i, subject := range subjects {
		if subject != nil {
			subjectList[i] = *subject
		}
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取科目列表成功",
		Data:    subjectList,
		Total:   int(total),
		Page:    req.Page,
		Size:    req.Size,
	})
}

// UpdateSubject 更新科目
// @Summary 更新科目信息
// @Description 根据ID更新科目信息
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Param subject body domain.UpdateSubjectRequest true "更新的科目信息"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/{id} [put]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "科目ID格式错误",
		})
		return
	}

	var req domain.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证和清理输入数据
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	subject, err := h.subjectService.UpdateSubject(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update subject",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "科目更新成功",
		Data:    subject,
	})
}

// DeleteSubject 删除科目
// @Summary 删除科目
// @Description 根据ID删除科目
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "科目ID格式错误",
		})
		return
	}

	err = h.subjectService.DeleteSubject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete subject",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "科目删除成功",
	})
}

// GetSubjectByCode 根据科目代码获取科目
// @Summary 根据科目代码获取科目
// @Description 根据科目代码获取科目详细信息
// @Tags subjects
// @Accept json
// @Produce json
// @Param code path string true "科目代码"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/code/{code} [get]
func (h *SubjectHandler) GetSubjectByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid code",
			Message: "科目代码不能为空",
		})
		return
	}

	subject, err := h.subjectService.GetSubjectByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Subject not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取科目成功",
		Data:    subject,
	})
}

// GetActiveSubjects 获取活跃科目列表
// @Summary 获取活跃科目列表
// @Description 获取状态为活跃的科目列表
// @Tags subjects
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/active [get]
func (h *SubjectHandler) GetActiveSubjects(c *gin.Context) {
	subjects, err := h.subjectService.GetActiveSubjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get active subjects",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取活跃科目列表成功",
		Data:    subjects,
	})
}
