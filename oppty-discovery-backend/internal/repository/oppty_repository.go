package repository

import (
	"context"
	"errors"
	"time"

	"oppty-discovery-backend/internal/model"

	"gorm.io/gorm"
)

/**
 * desc: 商机分析任务与报告数据访问接口与实现
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type OpptyRepository interface {
	CreateJob(ctx context.Context, job *model.OpptyJob) error
	UpdateJobProgress(ctx context.Context, jobID, stage, message string, percent int, status string) error
	UpdateJobError(ctx context.Context, jobID, errorMsg string) error
	GetJobByID(ctx context.Context, jobID string) (*model.OpptyJob, error)
	CreateReports(ctx context.Context, reports []model.OpptyReport) error
	GetReportsByJobID(ctx context.Context, jobID string) ([]model.OpptyReport, error)
	SaveRawData(ctx context.Context, raw *model.JobRawData) error
	GetRawDataByJobID(ctx context.Context, jobID string) (*model.JobRawData, error)
	ListJobsByUserID(ctx context.Context, userID uint64, pageNo, pageSize int, keyword string) ([]model.OpptyJob, int64, error)
	DeleteJob(ctx context.Context, jobID string, userID uint64) error
}

type opptyRepositoryImpl struct {
	db *gorm.DB
}

func NewOpptyRepository(db *gorm.DB) OpptyRepository {
	return &opptyRepositoryImpl{db: db}
}

func (r *opptyRepositoryImpl) CreateJob(ctx context.Context, job *model.OpptyJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *opptyRepositoryImpl) UpdateJobProgress(ctx context.Context, jobID, stage, message string, percent int, status string) error {
	updates := map[string]interface{}{
		"progress_stage":   stage,
		"progress_message": message,
		"progress_percent": percent,
		"updated_at":       time.Now(),
	}
	if status != "" {
		updates["status"] = status
	}
	return r.db.WithContext(ctx).Model(&model.OpptyJob{}).Where("id = ?", jobID).Updates(updates).Error
}

func (r *opptyRepositoryImpl) UpdateJobError(ctx context.Context, jobID, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&model.OpptyJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":           "failed",
		"progress_stage":   "failed",
		"progress_message": "任务执行失败",
		"error_message":    errorMsg,
		"updated_at":       time.Now(),
	}).Error
}

func (r *opptyRepositoryImpl) GetJobByID(ctx context.Context, jobID string) (*model.OpptyJob, error) {
	var job model.OpptyJob
	err := r.db.WithContext(ctx).Where("id = ?", jobID).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *opptyRepositoryImpl) CreateReports(ctx context.Context, reports []model.OpptyReport) error {
	if len(reports) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&reports).Error
}

func (r *opptyRepositoryImpl) GetReportsByJobID(ctx context.Context, jobID string) ([]model.OpptyReport, error) {
	var reports []model.OpptyReport
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Order("priority_score DESC").Find(&reports).Error
	return reports, err
}

func (r *opptyRepositoryImpl) SaveRawData(ctx context.Context, raw *model.JobRawData) error {
	return r.db.WithContext(ctx).Create(raw).Error
}

func (r *opptyRepositoryImpl) GetRawDataByJobID(ctx context.Context, jobID string) (*model.JobRawData, error) {
	var raw model.JobRawData
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).First(&raw).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &raw, nil
}

func (r *opptyRepositoryImpl) ListJobsByUserID(ctx context.Context, userID uint64, pageNo, pageSize int, keyword string) ([]model.OpptyJob, int64, error) {
	var jobs []model.OpptyJob
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OpptyJob{}).Where("user_id = ?", userID)
	if keyword != "" {
		query = query.Where("keywords LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pageNo - 1) * pageSize
	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&jobs).Error
	return jobs, total, err
}

func (r *opptyRepositoryImpl) DeleteJob(ctx context.Context, jobID string, userID uint64) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", jobID, userID).Delete(&model.OpptyJob{}).Error
}
