package handler

import (
	"net/http"
	"strconv"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

// StudentHandler 学生处理器
type StudentHandler struct {
	studentService *service.StudentService
}

// NewStudentHandler 创建新的学生处理器
func NewStudentHandler() *StudentHandler {
	return &StudentHandler{
		studentService: service.NewStudentService(),
	}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// CreateStudent 创建学生
// @Summary 创建新学生
// @Description 创建一个新的学生记录
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.CreateStudentRequest true "学生信息"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/students [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req domain.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	student, err := h.studentService.CreateStudent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建学生失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "学生创建成功",
		Data:    student,
	})
}

// GetStudent 获取单个学生信息
// @Summary 获取学生信息
// @Description 根据学生ID获取学生详细信息
// @Tags students
// @Produce json
// @Param id path int true "学生ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/{id} [get]
func (h *StudentHandler) GetStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的学生ID",
		})
		return
	}

	student, err := h.studentService.GetStudentByID(id)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "学生不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "获取学生信息失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    student,
	})
}

// GetStudents 获取学生列表
// @Summary 获取学生列表
// @Description 分页获取学生列表
// @Tags students
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/students [get]
func (h *StudentHandler) GetStudents(c *gin.Context) {
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

	students, total, err := h.studentService.GetAllStudents(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取学生列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取成功",
		Data:    students,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// UpdateStudent 更新学生信息
// @Summary 更新学生信息
// @Description 根据学生ID更新学生信息
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "学生ID"
// @Param student body models.UpdateStudentRequest true "更新的学生信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/{id} [put]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的学生ID",
		})
		return
	}

	var req domain.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	student, err := h.studentService.UpdateStudent(id, req)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "学生不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "更新学生信息失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新成功",
		Data:    student,
	})
}

// DeleteStudent 删除学生
// @Summary 删除学生
// @Description 根据学生ID删除学生
// @Tags students
// @Produce json
// @Param id path int true "学生ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的学生ID",
		})
		return
	}

	err = h.studentService.DeleteStudent(id)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "学生不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    500,
				Message: "删除学生失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除成功",
	})
}
