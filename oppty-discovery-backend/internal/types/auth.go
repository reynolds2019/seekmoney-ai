package types

/**
 * desc: 用户认证与账号管理请求与响应 DTO/VO
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */

// 发送邮箱验证码请求入参
type SendCodeReq struct {
	Email string `json:"email" binding:"required,email"`
}

// 发送验证码响应出参
type SendCodeResp struct {
	CooldownSeconds int `json:"cooldownSeconds"`
}

// 校验验证码与无感登录请求入参
type VerifyCodeReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// 用户简要信息出参
type UserProfileVO struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	IsNewUser bool   `json:"isNewUser"`
}

// 认证登录响应出参
type AuthRespVO struct {
	Token string        `json:"token"`
	User  UserProfileVO `json:"user"`
}
