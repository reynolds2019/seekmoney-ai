package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 统一接口响应规范结构体
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type Result struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// 成功响应 (HTTP 200)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code:      200,
		Message:   "操作成功",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 异步任务接收响应 (HTTP 202)
func Accepted(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusAccepted, Result{
		Code:      202,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 失败响应
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Result{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}
