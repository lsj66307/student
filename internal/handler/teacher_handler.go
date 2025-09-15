package handler

import (
	"net/http"
	"strconv"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

// TeacherHandler 老师处理器
type TeacherHandler struct {
	teacherService *service.TeacherService
}

// NewTeacherHandler 创建新的老师处理器
func NewTeacherHandler(teacherService *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

// CreateTeacher 创建老师
// @Summary 创建新老师
// @Description 创建一个新的老师记录
// @Tags teachers
// @Accept json
// @Produce json
// @Param teacher body models.CreateTeacherRequest true "老师信息"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/teachers [post]
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req domain.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	teacher, err := h.teacherService.CreateTeacher(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建老师失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "老师创建成功",
		Data:    teacher,
	})
}

// GetTeacher 获取单个老师信息
// @Summary 获取老师信息
// @Description 根据老师ID获取老师详细信息
// @Tags teachers
// @Produce json
// @Param id path int true "老师ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的老师ID",
		})
		return
	}

	teacher, err := h.teacherService.GetTeacherByID(id)
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "老师不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "获取老师信息失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    teacher,
	})
}

// GetTeachers 获取老师列表
// @Summary 获取老师列表
// @Description 分页获取老师列表
// @Tags teachers
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/teachers [get]
func (h *TeacherHandler) GetTeachers(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	teachers, total, err := h.teacherService.GetAllTeachers(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取老师列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    teachers,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// UpdateTeacher 更新老师信息
// @Summary 更新老师信息
// @Description 根据老师ID更新老师信息
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "老师ID"
// @Param teacher body models.UpdateTeacherRequest true "更新的老师信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/teachers/{id} [put]
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的老师ID",
		})
		return
	}

	var req domain.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	teacher, err := h.teacherService.UpdateTeacher(id, req)
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "老师不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "更新老师信息失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新成功",
		Data:    teacher,
	})
}

// DeleteTeacher 删除老师
// @Summary 删除老师
// @Description 根据老师ID删除老师
// @Tags teachers
// @Produce json
// @Param id path int true "老师ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/teachers/{id} [delete]
func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的老师ID",
		})
		return
	}

	err = h.teacherService.DeleteTeacher(id)
	if err != nil {
		if err.Error() == "teacher not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "老师不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "删除老师失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除成功",
	})
}
