package middleware

import (
	"net/http"
	"strings"

	"oppty-discovery-backend/internal/constant"
	"oppty-discovery-backend/internal/service"
	"oppty-discovery-backend/pkg/errors"
	"oppty-discovery-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

/**
 * desc: JWT 鉴权与当前用户会话注入中间件
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constant.HeaderAuthorization)

		// 支持从 Header 或 Query 或 Cookie 中获取凭据
		tokenString := ""
		if authHeader != "" && strings.HasPrefix(authHeader, constant.BearerPrefix) {
			tokenString = strings.TrimPrefix(authHeader, constant.BearerPrefix)
		} else {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			// 如果没有任何 Token，给一个默认的匿名体验用户（ID=1），避免阻断无感初次体验
			c.Set(constant.ContextKeyUserID, uint64(1))
			c.Set(constant.ContextKeyUserEmail, "demo@oppty.ai")
			c.Next()
			return
		}

		userID, email, err := authService.ParseToken(tokenString)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, errors.CodeUnauthorized, "登录状态已失效，请重新登录")
			c.Abort()
			return
		}

		c.Set(constant.ContextKeyUserID, userID)
		c.Set(constant.ContextKeyUserEmail, email)
		c.Next()
	}
}
