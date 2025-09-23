package service

import (
	"student-management-system/internal/domain"
	"student-management-system/internal/repository"
	"student-management-system/pkg/logger"
)

// ScoreService 成绩服务接口
type ScoreService interface {
	CreateScore(req *domain.CreateScoreRequest) (*domain.Score, error)
	GetScoreByID(id int) (*domain.Score, error)
	UpdateScore(id int, req *domain.UpdateScoreRequest) (*domain.Score, error)
	DeleteScore(id int) error
	ListScores(req *domain.ScoreListRequest) ([]*domain.Score, int64, error)
}

// scoreService 成绩服务实现
type scoreService struct {
	scoreRepo repository.ScoreRepository
}

// NewScoreService 创建成绩服务实例
func NewScoreService(scoreRepo repository.ScoreRepository) ScoreService {
	return &scoreService{
		scoreRepo: scoreRepo,
	}
}

// CreateScore 创建成绩
func (s *scoreService) CreateScore(req *domain.CreateScoreRequest) (*domain.Score, error) {
	logger.Info("Creating score", "student_id", req.StudentID, "subject_id", req.SubjectID)

	score := &domain.Score{
		StudentID: req.StudentID,
		SubjectID: req.SubjectID,
		Score:     req.Score,
		Semester:  req.Semester,
		ExamType:  req.ExamType,
	}

	err := s.scoreRepo.Create(score)
	if err != nil {
		logger.Error("Failed to create score", "error", err)
		return nil, err
	}

	logger.Info("Score created successfully", "score_id", score.ID)
	return score, nil
}

// GetScoreByID 根据ID获取成绩
func (s *scoreService) GetScoreByID(id int) (*domain.Score, error) {
	logger.Info("Getting score by ID", "score_id", id)

	score, err := s.scoreRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get score", "score_id", id, "error", err)
		return nil, err
	}

	return score, nil
}

// UpdateScore 更新成绩
func (s *scoreService) UpdateScore(id int, req *domain.UpdateScoreRequest) (*domain.Score, error) {
	logger.Info("Updating score", "score_id", id)

	// 先获取现有成绩
	score, err := s.scoreRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get score for update", "score_id", id, "error", err)
		return nil, err
	}

	// 更新字段
	if req.Score > 0 {
		score.Score = req.Score
	}
	if req.Semester != "" {
		score.Semester = req.Semester
	}
	if req.ExamType != "" {
		score.ExamType = req.ExamType
	}

	err = s.scoreRepo.Update(score)
	if err != nil {
		logger.Error("Failed to update score", "score_id", id, "error", err)
		return nil, err
	}

	logger.Info("Score updated successfully", "score_id", id)
	return score, nil
}

// DeleteScore 删除成绩
func (s *scoreService) DeleteScore(id int) error {
	logger.Info("Deleting score", "score_id", id)

	err := s.scoreRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete score", "score_id", id, "error", err)
		return err
	}

	logger.Info("Score deleted successfully", "score_id", id)
	return nil
}

// ListScores 获取成绩列表
func (s *scoreService) ListScores(req *domain.ScoreListRequest) ([]*domain.Score, int64, error) {
	logger.Info("Listing scores", "page", req.Page, "size", req.Size)

	scores, total, err := s.scoreRepo.List(req)
	if err != nil {
		logger.Error("Failed to list scores", "error", err)
		return nil, 0, err
	}

	logger.Info("Scores listed successfully", "total", total, "returned", len(scores))
	return scores, total, nil
}