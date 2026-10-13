package model

import (
	"time"

	"gorm.io/gorm"
)

/**
 * desc: 商机分析报告明细持久化模型
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type OpptyReport struct {
	ID                  uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	JobID               string         `gorm:"index;type:varchar(64);not null;column:job_id" json:"jobId"`
	UserID              uint64         `gorm:"index;not null;column:user_id" json:"userId"`
	ClusterCode         string         `gorm:"type:varchar(32);not null;column:cluster_code" json:"clusterCode"`
	SignalSize          int            `gorm:"not null;default:0;column:signal_size" json:"signalSize"`
	OneLinePain         string         `gorm:"type:text;not null;column:one_line_pain" json:"oneLinePain"`
	PaidInterest        string         `gorm:"type:varchar(16);not null;default:'Medium';column:paid_interest" json:"paidInterest"`
	Rationale           string         `gorm:"type:text;not null;column:rationale" json:"rationale"`
	PotentialProduct    string         `gorm:"type:text;not null;column:potential_product" json:"potentialProduct"`
	PainDepth           string         `gorm:"type:text;not null;column:pain_depth" json:"painDepth"`
	MarketLandscape     string         `gorm:"type:text;not null;column:market_landscape" json:"marketLandscape"`
	MVPPlan             string         `gorm:"type:text;not null;column:mvp_plan" json:"mvpPlan"`
	KeywordRelevance    int            `gorm:"not null;default:0;column:keyword_relevance" json:"keywordRelevance"`
	PriorityScore       int            `gorm:"not null;default:0;column:priority_score" json:"priorityScore"`
	RepresentativeTexts string         `gorm:"type:text;not null;column:representative_texts" json:"representativeTexts"`
	MinioMarkdownPath   string         `gorm:"type:varchar(255);column:minio_markdown_path" json:"minioMarkdownPath"`
	MinioHTMLPath       string         `gorm:"type:varchar(255);column:minio_html_path" json:"minioHtmlPath"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (OpptyReport) TableName() string {
	return "t_oppty_report"
}
