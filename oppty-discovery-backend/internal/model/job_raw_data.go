package model

import (
	"time"
)

/**
 * desc: 原始抓取数据与原声快照持久化模型
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type JobRawData struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	JobID        string    `gorm:"uniqueIndex:uk_raw_job_id;type:varchar(64);not null;column:job_id" json:"jobId"`
	VideoCount   int       `gorm:"not null;default:0;column:video_count" json:"videoCount"`
	CommentCount int       `gorm:"not null;default:0;column:comment_count" json:"commentCount"`
	TextCount    int       `gorm:"not null;default:0;column:text_count" json:"textCount"`
	RawVideos    string    `gorm:"type:text;not null;column:raw_videos" json:"rawVideos"`
	RawComments  string    `gorm:"type:text;not null;column:raw_comments" json:"rawComments"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (JobRawData) TableName() string {
	return "t_job_raw_data"
}
