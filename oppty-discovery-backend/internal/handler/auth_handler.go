package handler

import (
	"net/http"

	"oppty-discovery-backend/internal/service"
	"oppty-discovery-backend/internal/types"
	"oppty-discovery-backend/pkg/errors"
	"oppty-discovery-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 用户认证与无感注册控制器
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SendCode 发送邮箱验证码
// @param c *gin.Context
func (h *AuthHandler) SendCode(c *gin.Context) {
	var req types.SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "请输入有效的邮箱地址")
		return
	}

	resp, err := h.authService.SendEmailCode(c.Request.Context(), req.Email)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, resp)
}

// VerifyCode 校验验证码与无感注册建档
// @param c *gin.Context
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req types.VerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "验证码必须为6位数字")
		return
	}

	resp, err := h.authService.VerifyCodeAndLogin(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, resp)
}

// Logout 退出登录
// @param c *gin.Context
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"message": "已安全退出登录"})
}
