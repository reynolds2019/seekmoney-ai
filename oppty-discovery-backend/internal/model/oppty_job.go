package model

import (
	"time"

	"gorm.io/gorm"
)

/**
 * desc: 商机分析挖掘任务持久化模型
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type OpptyJob struct {
	ID                  string         `gorm:"primaryKey;type:varchar(64);column:id" json:"id"`
	UserID              uint64         `gorm:"index:idx_user_created;not null;column:user_id" json:"userId"`
	Status              string         `gorm:"index;type:varchar(32);not null;default:'processing';column:status" json:"status"`
	ProgressStage       string         `gorm:"type:varchar(32);not null;default:'init';column:progress_stage" json:"progressStage"`
	ProgressPercent     int            `gorm:"not null;default:0;column:progress_percent" json:"progressPercent"`
	ProgressMessage     string         `gorm:"type:varchar(255);not null;column:progress_message" json:"progressMessage"`
	Keywords            string         `gorm:"type:text;not null;column:keywords" json:"keywords"`
	DataSource          string         `gorm:"type:varchar(32);not null;default:'tiktok';column:data_source" json:"dataSource"`
	DeepCrawl           bool           `gorm:"not null;default:true;column:deep_crawl" json:"deepCrawl"`
	MaxVideos           int            `gorm:"not null;default:10;column:max_videos" json:"maxVideos"`
	MaxCommentsPerVideo int            `gorm:"not null;default:20;column:max_comments_per_video" json:"maxCommentsPerVideo"`
	EstimatedCost       float64        `gorm:"type:decimal(10,4);not null;default:0.0000;column:estimated_cost" json:"estimatedCost"`
	ErrorMessage        string         `gorm:"type:text;column:error_message" json:"errorMessage"`
	CreatedAt           time.Time      `gorm:"index:idx_user_created,sort:desc;column:created_at" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (OpptyJob) TableName() string {
	return "t_oppty_job"
}
