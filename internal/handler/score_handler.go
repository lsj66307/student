package handler

import (
	"net/http"
	"strconv"
	"student-management-system/internal/domain"
	"student-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

// ScoreHandler 成绩处理器
type ScoreHandler struct {
	scoreService service.ScoreService
}

// NewScoreHandler 创建新的成绩处理器
func NewScoreHandler(scoreService service.ScoreService) *ScoreHandler {
	return &ScoreHandler{
		scoreService: scoreService,
	}
}

// CreateScore 创建成绩
func (h *ScoreHandler) CreateScore(c *gin.Context) {
	var req domain.CreateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, err := h.scoreService.CreateScore(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": score})
}

// GetScore 获取成绩详情
func (h *ScoreHandler) GetScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score ID"})
		return
	}

	score, err := h.scoreService.GetScoreByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": score})
}

// GetScores 获取成绩列表
func (h *ScoreHandler) GetScores(c *gin.Context) {
	var req domain.ScoreListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	scores, total, err := h.scoreService.ListScores(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换指针切片为值切片
	scoreList := make([]domain.Score, len(scores))
	for i, score := range scores {
		scoreList[i] = *score
	}

	response := domain.ScoreListResponse{
		Scores: scoreList,
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateScore 更新成绩
func (h *ScoreHandler) UpdateScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score ID"})
		return
	}

	var req domain.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, err := h.scoreService.UpdateScore(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": score})
}

// DeleteScore 删除成绩
func (h *ScoreHandler) DeleteScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score ID"})
		return
	}

	err = h.scoreService.DeleteScore(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Score deleted successfully"})
}