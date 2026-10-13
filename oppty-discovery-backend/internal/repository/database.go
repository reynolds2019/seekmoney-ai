package repository

import (
	"fmt"
	"log"
	"time"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
 * desc: 数据库连接初始化与表结构自动迁移
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

var DB *gorm.DB

// InitDatabase 初始化 PostgreSQL 数据库连接
func InitDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()
	log.Printf("[数据库初始化] 正在连接 PostgreSQL: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Printf("[数据库警告] 连接目标库 %s 失败: %v，尝试连接默认库 postgres 进行自动探测", cfg.DBName, err)
		// 尝试先连接默认数据库 postgres 并自动建库
		fallbackDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.SSLMode)
		adminDB, adminErr := gorm.Open(postgres.Open(fallbackDSN), &gorm.Config{})
		if adminErr == nil {
			var exists bool
			adminDB.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", cfg.DBName).Scan(&exists)
			if !exists {
				log.Printf("[数据库建库] 正在自动创建数据库 %s...", cfg.DBName)
				adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s ENCODING 'UTF8'", cfg.DBName, cfg.User))
			}
			// 再次连接目标数据库
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}
	}

	if err != nil {
		return nil, fmt.Errorf("数据库连接初始化失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底座连接池失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移业务表结构
	log.Println("[数据库迁移] 正在执行数据模型表结构自动迁移...")
	err = db.AutoMigrate(
		&model.User{},
		&model.OpptyJob{},
		&model.OpptyReport{},
		&model.JobRawData{},
	)
	if err != nil {
		log.Printf("[数据库迁移警告] 自动迁移表结构存在异常: %v", err)
	} else {
		log.Println("[数据库迁移] 数据表自动迁移成功！")
	}

	DB = db
	return db, nil
}
