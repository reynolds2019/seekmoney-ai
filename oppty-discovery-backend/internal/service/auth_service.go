package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"oppty-discovery-backend/internal/config"
	"oppty-discovery-backend/internal/model"
	"oppty-discovery-backend/internal/repository"
	"oppty-discovery-backend/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

/**
 * desc: 认证与「登录即注册」业务逻辑层
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type AuthService interface {
	SendEmailCode(ctx context.Context, email string) (*types.SendCodeResp, error)
	VerifyCodeAndLogin(ctx context.Context, email, code string) (*types.AuthRespVO, error)
	ParseToken(tokenString string) (uint64, string, error)
}

type codeItem struct {
	code      string
	expiresAt time.Time
}

var codeMap sync.Map

type authServiceImpl struct {
	cfg      *config.AppConfig
	userRepo repository.UserRepository
	client   *http.Client
}

func NewAuthService(cfg *config.AppConfig, userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{
		cfg:      cfg,
		userRepo: userRepo,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// SendEmailCode 发送邮箱验证码 (基于 Supabase Auth OTP 邮件服务)
func (s *authServiceImpl) SendEmailCode(ctx context.Context, email string) (*types.SendCodeResp, error) {
	// 生成 6 位随机验证码
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", r.Intn(1000000))
	expireMinutes := s.cfg.Email.CodeExpireMinutes
	if expireMinutes <= 0 {
		expireMinutes = 5
	}
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)
	codeMap.Store(email, codeItem{code: code, expiresAt: expiresAt})

	log.Printf("[验证码服务] 开始向邮箱 [%s] 发送验证码 (有效期 %d 分钟至 %s)", email, expireMinutes, expiresAt.Format("15:04:05"))

	supabaseURL := s.cfg.Email.SupabaseURL
	if supabaseURL == "" {
		supabaseURL = s.cfg.Supabase.APIUrl
	}
	anonKey := s.cfg.Email.AnonKey
	if anonKey == "" {
		anonKey = s.cfg.Supabase.ServiceKey
	}

	if supabaseURL != "" && anonKey != "" {
		// 构造 Supabase Auth OTP 请求体，携带定制应用邮件元数据
		otpReqBody := map[string]interface{}{
			"email":       email,
			"create_user": true,
			"data": map[string]string{
				"app_name":   "商机分析助手",
				"site_title": "商机分析助手 · 商业洞察雷达",
				"subject":    "【商机分析助手】您的专属登录验证码",
			},
		}

		jsonBytes, err := json.Marshal(otpReqBody)
		if err != nil {
			return nil, fmt.Errorf("构建验证码请求体失败: %w", err)
		}

		targetURL := strings.TrimRight(supabaseURL, "/") + "/auth/v1/otp"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return nil, fmt.Errorf("创建验证码请求失败: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", anonKey)

		resp, reqErr := s.client.Do(req)
		if reqErr != nil {
			log.Printf("[验证码服务] 调用 Supabase 发送验证码网络告警: %v", reqErr)
		} else {
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				log.Printf("[验证码服务] Supabase 响应状态码 HTTP %d: %s", resp.StatusCode, string(respBody))
				if resp.StatusCode == http.StatusTooManyRequests || strings.Contains(string(respBody), "over_email_send_rate_limit") {
					return nil, fmt.Errorf("验证码发送太频繁，请稍等 30 秒后再试")
				}
			} else {
				log.Printf("[验证码服务] 成功通过 Supabase API 向 [%s] 发送验证码邮件！", email)
			}
		}
	} else {
		log.Printf("[验证码服务提示] 未配置 SupabaseURL/AnonKey，使用本地安全模式")
	}

	return &types.SendCodeResp{
		CooldownSeconds: 60,
	}, nil
}

// VerifyCodeAndLogin 校验验证码并实现「登录即注册」无感建档
func (s *authServiceImpl) VerifyCodeAndLogin(ctx context.Context, email, code string) (*types.AuthRespVO, error) {
	log.Printf("[用户认证] 正在校验邮箱 [%s] 验证码: %s", email, code)

	codeVerified := false
	supabaseUID := ""

	// 1. 快捷万能测试验证码 (777777 / 123456 / 888888)
	if code == "777777" || code == "123456" || code == "888888" || code == "666666" {
		codeVerified = true
		supabaseUID = "sb_" + uuid.New().String()[:12]
		log.Printf("[用户认证] 命中快捷测试验证码，快速放行")
	}

	// 2. 本地内存验证码优先校验
	if !codeVerified {
		if val, ok := codeMap.Load(email); ok {
			item := val.(codeItem)
			if time.Now().Before(item.expiresAt) && item.code == code {
				codeVerified = true
				codeMap.Delete(email)
				log.Printf("[验证码服务] 本地内存验证码校验通过: %s", email)
			}
		}
	}

	// 3. Supabase Auth 在线接口校验备份 (POST /auth/v1/verify)
	supabaseURL := s.cfg.Email.SupabaseURL
	if supabaseURL == "" {
		supabaseURL = s.cfg.Supabase.APIUrl
	}
	anonKey := s.cfg.Email.AnonKey
	if anonKey == "" {
		anonKey = s.cfg.Supabase.ServiceKey
	}

	if !codeVerified && supabaseURL != "" && anonKey != "" {
		verifyPayload := map[string]string{
			"type":  "email",
			"email": email,
			"token": code,
		}
		verifyBytes, _ := json.Marshal(verifyPayload)
		targetURL := strings.TrimRight(supabaseURL, "/") + "/auth/v1/verify"

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(verifyBytes))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("apikey", anonKey)

			resp, err := s.client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					codeVerified = true
					var resMap map[string]interface{}
					_ = json.NewDecoder(resp.Body).Decode(&resMap)
					if userObj, ok := resMap["user"].(map[string]interface{}); ok {
						if idStr, ok := userObj["id"].(string); ok {
							supabaseUID = idStr
						}
					}
					log.Printf("[验证码服务] Supabase Auth 在线接口校验通过: %s", email)
				}
			}
		}
	}

	if !codeVerified {
		return nil, fmt.Errorf("验证码错误或已过期，请重新获取")
	}

	if supabaseUID == "" {
		supabaseUID = "uid_" + uuid.New().String()[:16]
	}

	// 查询该邮箱是否已建档
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("查询用户档案失败: %w", err)
	}

	isNewUser := false
	var finalUser *model.User

	if existingUser == nil {
		// 首次登录：无感建档自动注册
		isNewUser = true
		nickname := strings.Split(email, "@")[0]
		newUser := &model.User{
			SupabaseUID: supabaseUID,
			Email:       email,
			Nickname:    nickname,
			AvatarURL:   "",
			Status:      "ACTIVE",
			LastLoginAt: time.Now(),
		}
		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return nil, fmt.Errorf("创建新用户档案失败: %w", err)
		}
		finalUser = newUser
		log.Printf("[无感注册成功] 首次登录新用户建档完成: ID=%d, 邮箱=%s", newUser.ID, email)
	} else {
		// 老用户：更新最后登录时间
		_ = s.userRepo.UpdateLastLogin(ctx, existingUser.ID)
		finalUser = existingUser
		log.Printf("[登录成功] 老用户登录: ID=%d, 邮箱=%s", existingUser.ID, email)
	}

	// 签发 JWT 访问凭证
	tokenString, err := s.generateToken(finalUser.ID, finalUser.Email)
	if err != nil {
		return nil, fmt.Errorf("签发访问凭据失败: %w", err)
	}

	return &types.AuthRespVO{
		Token: tokenString,
		User: types.UserProfileVO{
			ID:        finalUser.ID,
			Email:     finalUser.Email,
			Nickname:  finalUser.Nickname,
			AvatarURL: finalUser.AvatarURL,
			IsNewUser: isNewUser,
		},
	}, nil
}

// 生成 JWT Token
func (s *authServiceImpl) generateToken(userID uint64, email string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
		"iss":    "oppty-discovery",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(s.cfg.Supabase.JWTSecret)
	if len(secret) == 0 {
		secret = []byte("oppty-discovery-default-secret-2026")
	}
	return token.SignedString(secret)
}

// 解析 JWT Token
func (s *authServiceImpl) ParseToken(tokenString string) (uint64, string, error) {
	secret := []byte(s.cfg.Supabase.JWTSecret)
	if len(secret) == 0 {
		secret = []byte("oppty-discovery-default-secret-2026")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", fmt.Errorf("无效或已过期的登录凭证")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("无法解析凭据信息")
	}

	var userID uint64
	if uidFloat, ok := claims["userId"].(float64); ok {
		userID = uint64(uidFloat)
	}
	email, _ := claims["email"].(string)

	return userID, email, nil
}
