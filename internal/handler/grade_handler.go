package handler

import (
	"net/http"
	"strconv"

	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

// GradeHandler 成绩处理器
type GradeHandler struct {
	gradeService *service.GradeService
}

// NewGradeHandler 创建新的成绩处理器
func NewGradeHandler(gradeService *service.GradeService) *GradeHandler {
	return &GradeHandler{
		gradeService: gradeService,
	}
}

// CreateGrade 创建成绩
// @Summary 创建成绩
// @Description 创建新的学生成绩记录
// @Tags 成绩管理
// @Accept json
// @Produce json
// @Param grade body domain.CreateGradeRequest true "成绩信息"
// @Success 201 {object} Response "成绩创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/grades [post]
func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var req domain.CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	grade, err := h.gradeService.CreateGrade(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Create grade failed",
			Message: "创建成绩失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "成绩创建成功",
		Data:    grade,
	})
}

// GetGrades 获取成绩列表
// @Summary 获取成绩列表
// @Description 获取成绩列表，支持按学生筛选和分页
// @Tags 成绩管理
// @Accept json
// @Produce json
// @Param student_id query int false "学生ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} PaginatedResponse "获取成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/grades [get]
func (h *GradeHandler) GetGrades(c *gin.Context) {
	var params domain.GradeQueryParams

	// 解析查询参数
	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		if studentID, err := strconv.Atoi(studentIDStr); err == nil {
			params.StudentID = &studentID
		}
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}
	if params.Page == 0 {
		params.Page = 1
	}

	if sizeStr := c.Query("size"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			params.Size = size
		}
	}
	if params.Size == 0 {
		params.Size = 10
	}

	grades, total, err := h.gradeService.GetAllGrades(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取成绩列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    grades,
		Total:   total,
		Page:    params.Page,
		Size:    params.Size,
	})
}

// GetGrade 获取单个成绩
// @Summary 获取单个成绩
// @Description 根据ID获取成绩详情
// @Tags 成绩管理
// @Accept json
// @Produce json
// @Param id path int true "成绩ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "成绩不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/grades/{id} [get]
func (h *GradeHandler) GetGrade(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的成绩ID",
		})
		return
	}

	grade, err := h.gradeService.GetGradeByID(id)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "成绩不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "获取成绩失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    grade,
	})
}

// UpdateGrade 更新成绩
// @Summary 更新成绩
// @Description 更新成绩信息
// @Tags 成绩管理
// @Accept json
// @Produce json
// @Param id path int true "成绩ID"
// @Param grade body models.UpdateGradeRequest true "更新的成绩信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "成绩不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/grades/{id} [put]
func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的成绩ID",
		})
		return
	}

	var req domain.UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	grade, err := h.gradeService.UpdateGrade(id, req)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "成绩不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "更新成绩失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新成功",
		Data:    grade,
	})
}

// DeleteGrade 删除成绩
// @Summary 删除成绩
// @Description 删除成绩记录
// @Tags 成绩管理
// @Accept json
// @Produce json
// @Param id path int true "成绩ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "成绩不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/grades/{id} [delete]
func (h *GradeHandler) DeleteGrade(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的成绩ID",
		})
		return
	}

	err = h.gradeService.DeleteGrade(id)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "成绩不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "删除成绩失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除成功",
	})
}
