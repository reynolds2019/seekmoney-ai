package config

import (
	"os"
	"testing"
)

/**
 * desc: 配置加载与环境变量安全注入单元测试
 * date: 2026-10-13 18:26:00
 * @author reynolds
 * @since v1.0.0
 */
func TestLoadConfig_EnvInjection(t *testing.T) {
	// 1. 设置模拟环境变量
	testToken := "test_token_from_env_123456"
	_ = os.Setenv("TIKHUB_API_TOKEN", testToken)
	defer os.Unsetenv("TIKHUB_API_TOKEN")

	// 2. 加载实际配置文件
	cfg, err := LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 3. 校验 TikHub API Token 是否被环境变量成功覆盖
	if cfg.TikHub.APIToken != testToken {
		t.Fatalf("预期 TikHub APIToken 为 [%s]，实际获取为 [%s]", testToken, cfg.TikHub.APIToken)
	}
	t.Logf("配置加载与环境变量安全注入测试通过，注入的Token: %s", cfg.TikHub.APIToken)
}

/**
 * desc: 验证无系统环境变量时自动从本地.env文件兜底注入Token
 * date: 2026-10-13 18:28:00
 * @author reynolds
 * @since v1.0.0
 */
func TestLoadConfig_DotEnvFallback(t *testing.T) {
	// 确保没有系统环境变量
	_ = os.Unsetenv("TIKHUB_API_TOKEN")

	cfg, err := LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 此时应成功从本地 .env 中注入真实凭证
	if cfg.TikHub.APIToken == "" {
		t.Fatalf("预期从本地 .env 注入 Token，但获取到空值")
	}
	t.Logf("成功从本地 .env 文件回退注入 Token，长度: %d", len(cfg.TikHub.APIToken))
}
