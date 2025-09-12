package handlers

import (
	"net/http"
	"strconv"

	"student-management-system/models"

	"github.com/gin-gonic/gin"
)

// GradeHandler 成绩处理器
type GradeHandler struct {
	gradeService *models.GradeService
}

// NewGradeHandler 创建成绩处理器实例
func NewGradeHandler(gradeService *models.GradeService) *GradeHandler {
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
// @Param grade body models.CreateGradeRequest true "成绩信息"
// @Success 201 {object} map[string]interface{} "成绩创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/grades [post]
func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var req models.CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	grade, err := h.gradeService.CreateGrade(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建成绩失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "成绩创建成功",
		"data":    grade,
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
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/grades [get]
func (h *GradeHandler) GetGrades(c *gin.Context) {
	var params models.GradeQueryParams

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取成绩列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    grades,
		"total":   total,
		"page":    params.Page,
		"size":    params.Size,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的成绩ID",
			"data":    nil,
		})
		return
	}

	grade, err := h.gradeService.GetGradeByID(id)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "成绩不存在",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取成绩失败: " + err.Error(),
				"data":    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    grade,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的成绩ID",
			"data":    nil,
		})
		return
	}

	var req models.UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	grade, err := h.gradeService.UpdateGrade(id, req)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "成绩不存在",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新成绩失败: " + err.Error(),
				"data":    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    grade,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的成绩ID",
			"data":    nil,
		})
		return
	}

	err = h.gradeService.DeleteGrade(id)
	if err != nil {
		if err.Error() == "成绩不存在" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "成绩不存在",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除成绩失败: " + err.Error(),
				"data":    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}