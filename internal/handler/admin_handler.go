package handler

import (
	"net/http"
	"strconv"

	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AdminHandler struct {
	adminService *service.AdminService
	logger       *logrus.Logger
}

func NewAdminHandler(adminService *service.AdminService, logger *logrus.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       logger,
	}
}

// CreateAdmin 创建管理员
// @Summary 创建管理员
// @Description 创建新的管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param admin body domain.CreateAdminRequest true "管理员信息"
// @Success 200 {object} Response{data=domain.AdminInfo} "创建成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/admin [post]
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req domain.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind create admin request")
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	adminInfo, err := h.adminService.CreateAdmin(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create admin")
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    adminInfo,
	})
}

// GetAdmin 获取管理员信息
// @Summary 获取管理员信息
// @Description 根据ID获取管理员详细信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID"
// @Success 200 {object} Response{data=domain.AdminInfo} "获取成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 404 {object} Response "管理员不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/admin/{id} [get]
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的管理员ID",
			Data:    nil,
		})
		return
	}

	adminInfo, err := h.adminService.GetAdminByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get admin")
		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    adminInfo,
	})
}

// UpdateAdmin 更新管理员信息
// @Summary 更新管理员信息
// @Description 更新指定管理员的信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID"
// @Param admin body domain.UpdateAdminRequest true "更新的管理员信息"
// @Success 200 {object} Response{data=domain.AdminInfo} "更新成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 404 {object} Response "管理员不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/admin/{id} [put]
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的管理员ID",
			Data:    nil,
		})
		return
	}

	var req domain.UpdateAdminRequest
	if err2 := c.ShouldBindJSON(&req); err2 != nil {
		h.logger.WithError(err2).Error("Failed to bind update admin request")
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	adminInfo, err := h.adminService.UpdateAdmin(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update admin")
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    adminInfo,
	})
}

// DeleteAdmin 删除管理员
// @Summary 删除管理员
// @Description 删除指定的管理员
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id path int true "管理员ID"
// @Success 200 {object} Response "删除成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 404 {object} Response "管理员不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/admin/{id} [delete]
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的管理员ID",
			Data:    nil,
		})
		return
	}

	err = h.adminService.DeleteAdmin(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete admin")
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "删除成功",
		Data:    nil,
	})
}

// ListAdmins 获取管理员列表
// @Summary 获取管理员列表
// @Description 分页获取管理员列表
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param account query string false "账号筛选"
// @Param name query string false "姓名筛选"
// @Success 200 {object} Response{data=domain.AdminListResponse} "获取成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/admin [get]
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	var req domain.AdminListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind list admin request")
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    nil,
		})
		return
	}

	response, err := h.adminService.ListAdmins(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list admins")
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    response,
	})
}
