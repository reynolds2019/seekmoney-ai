package service

import (
	"context"
	"testing"
	"time"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/model"
)

/**
 * desc: 认证服务单元测试用例
 * date: 2026-10-13 18:02:00
 * @author reynolds
 * @since v1.0.0
 */

// mockUserRepo 用于测试的模拟用户仓储
type mockUserRepo struct {
	users map[string]*model.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) FindBySupabaseUID(ctx context.Context, uid string) (*model.User, error) {
	for _, u := range m.users {
		if u.SupabaseUID == uid {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id uint64) error {
	return nil
}

// TestAuthService_SendAndVerify 测试邮件验证码发送与多重校验机制
func TestAuthService_SendAndVerify(t *testing.T) {
	cfg := &config.AppConfig{
		Supabase: config.SupabaseConfig{
			JWTSecret: "test_jwt_secret_32_characters_123456",
		},
		Email: config.EmailConfig{
			CodeExpireMinutes: 5,
		},
	}
	repo := &mockUserRepo{
		users: make(map[string]*model.User),
	}
	authSvc := NewAuthService(cfg, repo)
	ctx := context.Background()

	testEmail := "test_user_unit@example.com"

	t.Run("测试生成并发送验证码到内存与冷却时间", func(t *testing.T) {
		resp, err := authSvc.SendEmailCode(ctx, testEmail)
		if err != nil {
			t.Fatalf("发送验证码失败: %v", err)
		}
		if resp.CooldownSeconds != 60 {
			t.Errorf("期望冷却时间为 60 秒，实际为: %d", resp.CooldownSeconds)
		}

		// 检查本地内存缓存是否存在有效验证码
		val, ok := codeMap.Load(testEmail)
		if !ok {
			t.Fatalf("内存缓存中未找到邮箱验证码")
		}
		item := val.(codeItem)
		if len(item.code) != 6 {
			t.Errorf("验证码长度应为 6 位，实际为: %s", item.code)
		}
		if time.Now().After(item.expiresAt) {
			t.Errorf("验证码过期时间设置异常，已过期")
		}
	})

	t.Run("测试万能测试码直接通过校验并登录", func(t *testing.T) {
		loginResp, err := authSvc.VerifyCodeAndLogin(ctx, testEmail, "777777")
		if err != nil {
			t.Fatalf("万能验证码应直接校验通过，实际报错: %v", err)
		}
		if loginResp.Token == "" {
			t.Errorf("未返回有效 JWT Token")
		}
		if loginResp.User.Email != testEmail {
			t.Errorf("返回用户信息邮箱不匹配，期望 %s，实际 %s", testEmail, loginResp.User.Email)
		}
	})

	t.Run("测试内存中的真实验证码校验并核销", func(t *testing.T) {
		// 重新发送一次
		_, _ = authSvc.SendEmailCode(ctx, testEmail)
		val, _ := codeMap.Load(testEmail)
		validCode := val.(codeItem).code

		loginResp, err := authSvc.VerifyCodeAndLogin(ctx, testEmail, validCode)
		if err != nil {
			t.Fatalf("有效验证码应当通过校验，实际报错: %v", err)
		}
		if loginResp == nil || loginResp.Token == "" {
			t.Errorf("登录结果应包含有效令牌")
		}

		// 再次使用同一验证码应提示已过期或核销
		_, errRepeat := authSvc.VerifyCodeAndLogin(ctx, testEmail, validCode)
		if errRepeat == nil {
			t.Errorf("已核销的验证码二次使用应当被拦截")
		}
	})

	t.Run("测试错误验证码被正常拦截", func(t *testing.T) {
		_, err := authSvc.VerifyCodeAndLogin(ctx, testEmail, "999999")
		if err == nil {
			t.Errorf("错误验证码应被拒绝，但未返回错误")
		}
	})
}
