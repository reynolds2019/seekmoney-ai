package model

import (
	"time"

	"gorm.io/gorm"
)

/**
 * desc: 用户基础数据实体
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type User struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SupabaseUID string         `gorm:"uniqueIndex:uk_supabase_uid;type:varchar(64);not null;column:supabase_uid" json:"supabaseUid"`
	Email       string         `gorm:"uniqueIndex:uk_email;type:varchar(128);not null;column:email" json:"email"`
	Nickname    string         `gorm:"type:varchar(64);column:nickname" json:"nickname"`
	AvatarURL   string         `gorm:"type:varchar(255);column:avatar_url" json:"avatarUrl"`
	Status      string         `gorm:"type:varchar(32);not null;default:'ACTIVE';column:status" json:"status"`
	LastLoginAt time.Time      `gorm:"column:last_login_at" json:"lastLoginAt"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`
}

func (User) TableName() string {
	return "t_user"
}
