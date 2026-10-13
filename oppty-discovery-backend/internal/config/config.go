package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

/**
 * desc: 系统全局配置结构体与加载逻辑
 * date: 2026-10-13 11:25:00
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
	BaseURL        string  `yaml:"baseUrl"`
	APIKey         string  `yaml:"apiKey"`
	ModelName      string  `yaml:"modelName"`
	Temperature    float32 `yaml:"temperature"`
	TimeoutSeconds int     `yaml:"timeoutSeconds"`
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

// LoadConfig 从指定 YAML 文件加载配置
func LoadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 [%s]: %w", configPath, err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件内容失败: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}
