package handler

import (
	"net/http"
	"strconv"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"
	"student-management-system/pkg/validator"

	"github.com/gin-gonic/gin"
)

// StudentHandler 学生处理器
type StudentHandler struct {
	studentService *service.StudentService
	validator      *validator.CustomValidator
}

// NewStudentHandler 创建新的学生处理器
func NewStudentHandler(studentService *service.StudentService, validator *validator.CustomValidator) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
		validator:      validator,
	}
}

// CreateStudent 创建学生
// @Summary 创建新学生
// @Description 创建一个新的学生记录
// @Tags students
// @Accept json
// @Produce json
// @Param student body domain.CreateStudentRequest true "学生信息"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/students [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req domain.CreateStudentRequest
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

	// 清理输入数据
	req.Name = validator.SanitizeInput(req.Name)
	req.Email = validator.SanitizeInput(req.Email)
	req.Major = validator.SanitizeInput(req.Major)
	req.StudentID = validator.SanitizeInput(req.StudentID)
	req.Address = validator.SanitizeInput(req.Address)

	student, err := h.studentService.CreateStudent(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Create student failed",
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
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/students/{id} [get]
func (h *StudentHandler) GetStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid student ID",
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

	// 验证和清理输入数据
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	// 清理输入数据
	if req.Name != "" {
		req.Name = validator.SanitizeInput(req.Name)
	}
	if req.Email != "" {
		req.Email = validator.SanitizeInput(req.Email)
	}
	if req.Major != "" {
		req.Major = validator.SanitizeInput(req.Major)
	}
	if req.StudentID != "" {
		req.StudentID = validator.SanitizeInput(req.StudentID)
	}
	if req.Address != "" {
		req.Address = validator.SanitizeInput(req.Address)
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

// BatchCreateStudents 批量创建学生
// @Summary 批量创建学生
// @Description 批量创建多个学生记录
// @Tags students
// @Accept json
// @Produce json
// @Param students body []domain.CreateStudentRequest true "学生信息列表"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/batch [post]
func (h *StudentHandler) BatchCreateStudents(c *gin.Context) {
	var req domain.BatchCreateStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证批量请求
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	// 清理每个学生的输入数据
	for i := range req.Students {
		req.Students[i].Name = validator.SanitizeInput(req.Students[i].Name)
		req.Students[i].Email = validator.SanitizeInput(req.Students[i].Email)
		req.Students[i].Major = validator.SanitizeInput(req.Students[i].Major)
		req.Students[i].StudentID = validator.SanitizeInput(req.Students[i].StudentID)
		req.Students[i].Address = validator.SanitizeInput(req.Students[i].Address)
	}

	if len(req.Students) == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "学生列表不能为空",
		})
		return
	}

	if len(req.Students) > 100 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "批量创建学生数量不能超过100个",
		})
		return
	}

	students, err := h.studentService.BatchCreateStudents(req.Students)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "批量创建学生失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "批量创建成功",
		Data:    students,
	})
}

// BatchDeleteStudents 批量删除学生
// @Summary 批量删除学生
// @Description 根据ID列表批量删除学生
// @Tags students
// @Accept json
// @Produce json
// @Param ids body []int true "学生ID列表"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/batch [delete]
func (h *StudentHandler) BatchDeleteStudents(c *gin.Context) {
	var req domain.BatchDeleteStudentsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证批量删除请求
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "学生ID列表不能为空",
		})
		return
	}

	if len(req.IDs) > 100 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "批量删除学生数量不能超过100个",
		})
		return
	}

	err := h.studentService.BatchDeleteStudents(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "批量删除学生失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "批量删除成功",
	})
}

// TransferStudentMajor 学生转专业
// @Summary 学生转专业
// @Description 为学生办理转专业手续
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "学生ID"
// @Param transfer body domain.TransferMajorRequest true "转专业信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/students/{id}/transfer-major [post]
func (h *StudentHandler) TransferStudentMajor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的学生ID",
		})
		return
	}

	var req domain.TransferStudentMajorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证转专业请求
	if err := h.validator.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "数据验证失败: " + err.Error(),
		})
		return
	}

	// 清理输入数据
	req.NewMajor = validator.SanitizeInput(req.NewMajor)
	req.Reason = validator.SanitizeInput(req.Reason)

	err = h.studentService.TransferStudentMajor(id, req.NewMajor, req.Reason)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "学生不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "转专业失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "转专业成功",
	})
}
