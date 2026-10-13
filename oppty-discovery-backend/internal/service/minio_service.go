package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"oppty-discovery-backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/**
 * desc: MinIO 对象存储服务 (S3兼容)
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type MinioService interface {
	UploadReport(ctx context.Context, jobID, format, content string) (string, error)
	IsHealthy(ctx context.Context) bool
}

type minioServiceImpl struct {
	client     *minio.Client
	bucketName string
	publicBase string
}

func NewMinioService(cfg *config.MinioConfig) MinioService {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("[MinIO 客户端初始化警告] %v", err)
	} else {
		// 检查或自动创建专属文件桶 oppty-discovery
		ctx := context.Background()
		exists, err := client.BucketExists(ctx, cfg.BucketName)
		if err == nil && !exists {
			log.Printf("[MinIO 自动建桶] 正在创建存储桶: %s...", cfg.BucketName)
			err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Printf("[MinIO 自动建桶警告] 创建存储桶失败: %v", err)
			} else {
				log.Printf("[MinIO 自动建桶] 存储桶 %s 创建成功", cfg.BucketName)
			}
		}
	}

	return &minioServiceImpl{
		client:     client,
		bucketName: cfg.BucketName,
		publicBase: cfg.PublicBaseURL,
	}
}

// UploadReport 上传商机报告 (Markdown 或 HTML)
func (s *minioServiceImpl) UploadReport(ctx context.Context, jobID, format, content string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("MinIO 客户端未初始化")
	}

	contentType := "text/markdown; charset=utf-8"
	ext := "md"
	if strings.ToLower(format) == "html" {
		contentType = "text/html; charset=utf-8"
		ext = "html"
	}

	objectName := fmt.Sprintf("reports/%s.%s", jobID, ext)
	contentBytes := []byte(content)
	reader := bytes.NewReader(contentBytes)

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, int64(len(contentBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("[MinIO 上传失败] 文件 %s 上传失败: %v", objectName, err)
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s", strings.TrimRight(s.publicBase, "/"), objectName)
	log.Printf("[MinIO 上传成功] 文件已就绪: %s", fileURL)
	return fileURL, nil
}

// IsHealthy 检查 MinIO 服务连通性
func (s *minioServiceImpl) IsHealthy(ctx context.Context) bool {
	if s.client == nil {
		return false
	}
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	return err == nil && exists
}
