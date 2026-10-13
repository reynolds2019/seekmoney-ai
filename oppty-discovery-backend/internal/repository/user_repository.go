package repository

import (
	"context"
	"errors"
	"time"

	"oppty-discovery-backend/internal/model"

	"gorm.io/gorm"
)

/**
 * desc: 用户数据访问层接口与实现
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindBySupabaseUID(ctx context.Context, uid string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	UpdateLastLogin(ctx context.Context, id uint64) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindBySupabaseUID(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("supabase_uid = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) UpdateLastLogin(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		}).Error
}
