package handler

import (
	"context"
	"time"

	"oppty-discovery-backend/internal/constant"
	"oppty-discovery-backend/internal/repository"
	"oppty-discovery-backend/internal/service"
	"oppty-discovery-backend/internal/types"
	"oppty-discovery-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 系统健康自检与监控控制器
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type HealthHandler struct {
	minioService service.MinioService
}

func NewHealthHandler(minioService service.MinioService) *HealthHandler {
	return &HealthHandler{minioService: minioService}
}

// CheckHealth 服务健康状态检查
// @param c *gin.Context
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "connected"
	if repository.DB != nil {
		sqlDB, err := repository.DB.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			dbStatus = "disconnected"
		}
	} else {
		dbStatus = "uninitialized"
	}

	minioStatus := "disconnected"
	if h.minioService != nil && h.minioService.IsHealthy(ctx) {
		minioStatus = "connected"
	}

	resp := types.HealthRespVO{
		Service:   constant.AppNameZh,
		Status:    "healthy",
		Database:  dbStatus,
		Minio:     minioStatus,
		Timestamp: time.Now().UnixMilli(),
	}

	response.Success(c, resp)
}
