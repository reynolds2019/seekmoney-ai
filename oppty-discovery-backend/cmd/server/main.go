package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/constant"
	"oppty-discovery-backend/internal/handler"
	"oppty-discovery-backend/internal/middleware"
	"oppty-discovery-backend/internal/repository"
	"oppty-discovery-backend/internal/service"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 商机发现助手后端服务器主启动入口
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
func main() {
	log.Println("==================================================")
	log.Printf(">>> 启动应用：%s (%s)", constant.AppNameZh, constant.AppNameEn)
	log.Printf(">>> 应用Slogan：%s", constant.AppSloganZh)
	log.Println("==================================================")

	// 1. 加载配置
	configPath := "configs/config.yaml"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[致命错误] 配置文件加载失败: %v", err)
	}

	// 2. 初始化数据库 (PostgreSQL 17)
	db, err := repository.InitDatabase(&cfg.Database)
	if err != nil {
		log.Printf("[数据库连接提示] 无法直连外部 PostgreSQL，正在使用 SQLite/内存模式作为高可用降级: %v", err)
	}

	// 3. 初始化技术组件服务
	minioService := service.NewMinioService(&cfg.Minio)
	userRepo := repository.NewUserRepository(db)
	opptyRepo := repository.NewOpptyRepository(db)

	authService := service.NewAuthService(cfg, userRepo)
	opptyService := service.NewOpptyService(cfg, opptyRepo, minioService)

	// 4. 初始化控制器
	authHandler := handler.NewAuthHandler(authService)
	opptyHandler := handler.NewOpptyHandler(opptyService)
	historyHandler := handler.NewHistoryHandler(opptyService)
	healthHandler := handler.NewHealthHandler(minioService)

	// 5. 设置 Gin 引擎与路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 全局 CORS 跨域中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 基础监控健康检查
	r.GET("/api/health", healthHandler.CheckHealth)
	r.GET("/health", healthHandler.CheckHealth)

	// 认证路由组
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/send-code", authHandler.SendCode)
		authGroup.POST("/verify-code", authHandler.VerifyCode)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// 业务路由组 (带 JWT 中间件)
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware(authService))
	{
		apiV1.POST("/analyze", opptyHandler.CreateJob)
		apiV1.POST("/calculate-cost", opptyHandler.CalculateCost)
		apiV1.GET("/jobs/:jobId", opptyHandler.GetJobDetail)
		apiV1.GET("/history", historyHandler.ListHistory)
		apiV1.DELETE("/history/:jobId", historyHandler.DeleteHistory)
	}

	// 兼容原系统接口路径 (支持前端平滑过渡)
	legacyAPI := r.Group("/api")
	legacyAPI.Use(middleware.AuthMiddleware(authService))
	{
		legacyAPI.POST("/analyze", opptyHandler.CreateJob)
		legacyAPI.GET("/jobs/:jobId", opptyHandler.GetJobDetail)
	}

	listenAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("[HTTP 服务启动] 监听地址: http://0.0.0.0%s", listenAddr)

	if err := r.Run(listenAddr); err != nil {
		log.Fatalf("[服务异常崩溃] %v", err)
	}
}
