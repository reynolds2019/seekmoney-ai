package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

/**
 * desc: 系统全局配置结构体与加载逻辑
 * date: 2026-10-13 18:25:00
 * @author reynolds
 * @since v1.0.0
 */

// ServerConfig 服务运行配置
type ServerConfig struct {
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
	AppName string `yaml:"appName"`
}

// DatabaseConfig PostgreSQL 17 数据库连接配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

// DSN 获取 PostgreSQL 数据库连接串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// SupabaseConfig Supabase BaaS 统一鉴权配置
type SupabaseConfig struct {
	APIUrl     string `yaml:"apiUrl"`
	JWTSecret  string `yaml:"jwtSecret"`
	ServiceKey string `yaml:"serviceKey"`
}

// MinioConfig MinIO S3 对象存储配置
type MinioConfig struct {
	Endpoint      string `yaml:"endpoint"`
	AccessKey     string `yaml:"accessKey"`
	SecretKey     string `yaml:"secretKey"`
	UseSSL        bool   `yaml:"useSSL"`
	BucketName    string `yaml:"bucketName"`
	PublicBaseURL string `yaml:"publicBaseUrl"`
}

// LLMConfig 大模型分析引擎配置
type LLMConfig struct {
	BaseURL            string  `yaml:"baseUrl"`
	APIKey             string  `yaml:"apiKey"`
	ModelName          string  `yaml:"modelName"`
	Temperature        float32 `yaml:"temperature"`
	TimeoutSeconds     int     `yaml:"timeoutSeconds"`
	EmbeddingModel     string  `yaml:"embeddingModel"`     // 文本向量模型名称，默认 text-embedding-v4
	EmbeddingBatchSize int     `yaml:"embeddingBatchSize"` // 批处理文本上限，默认 10
	EmbeddingURL       string  `yaml:"embeddingUrl"`       // 向量模型接口地址，若为空则默认使用 baseUrl + "/embeddings"
}

// TikHubConfig 社交媒体数据源抓取配置
type TikHubConfig struct {
	BaseURL        string `yaml:"baseUrl"`
	APIToken       string `yaml:"apiToken"`
	RetryCount     int    `yaml:"retryCount"`
	TimeoutSeconds int    `yaml:"timeoutSeconds"`
}

// EmailConfig 邮件服务配置
type EmailConfig struct {
	UseInbucket       bool   `yaml:"use_inbucket"`
	InbucketURL       string `yaml:"inbucket_url"`
	SupabaseURL       string `yaml:"supabase_url"`
	AnonKey           string `yaml:"anon_key"`
	CodeExpireMinutes int    `yaml:"code_expire_minutes"`
}

// AppConfig 应用聚合总配置
type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Supabase SupabaseConfig `yaml:"supabase"`
	Email    EmailConfig    `yaml:"email"`
	Minio    MinioConfig    `yaml:"minio"`
	LLM      LLMConfig      `yaml:"llm"`
	TikHub   TikHubConfig   `yaml:"tikhub"`
}

var GlobalConfig *AppConfig

// loadDotEnv 自动搜索并加载本地 .env 文件至环境变量
func loadDotEnv(configPath string) {
	candidates := []string{
		os.Getenv("ENV_FILE_PATH"),
		".env",
		"../.env",
		"configs/.env",
		filepath.Join(filepath.Dir(configPath), ".env"),
		filepath.Join(filepath.Dir(configPath), "../.env"),
	}

	for _, path := range candidates {
		if path == "" {
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		log.Printf("[配置中心] 成功检测并加载本地环境文件: %s", path)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, "\"'")
				// 环境变量优先级高于 .env 文件：仅在系统环境变量未设置时注入
				if os.Getenv(k) == "" {
					_ = os.Setenv(k, v)
				}
			}
		}
		return
	}
}

// LoadConfig 从指定 YAML 文件加载配置
func LoadConfig(configPath string) (*AppConfig, error) {
	// 优先尝试加载本地 .env 配置文件
	loadDotEnv(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 [%s]: %w", configPath, err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件内容失败: %w", err)
	}

	// 敏感凭证环境变量优先覆盖逻辑
	if envToken := os.Getenv("TIKHUB_API_TOKEN"); envToken != "" {
		cfg.TikHub.APIToken = envToken
		log.Println("[配置中心] 已从环境变量/本地.env动态注入 TikHub API Token")
	} else if strings.HasPrefix(cfg.TikHub.APIToken, "${") {
		log.Println("[配置中心警告] 未检测到 TIKHUB_API_TOKEN 环境变量，且配置文件中使用占位符")
	}

	// 补全大模型 Embedding 默认配置项
	if cfg.LLM.EmbeddingModel == "" {
		cfg.LLM.EmbeddingModel = "text-embedding-v4"
	}
	if cfg.LLM.EmbeddingBatchSize <= 0 {
		cfg.LLM.EmbeddingBatchSize = 10
	}
	if cfg.LLM.EmbeddingURL == "" && cfg.LLM.BaseURL != "" {
		cfg.LLM.EmbeddingURL = fmt.Sprintf("%s/embeddings", strings.TrimRight(cfg.LLM.BaseURL, "/"))
	}

	GlobalConfig = &cfg
	return &cfg, nil
}
