package handler

import (
	"net/http"
	"strconv"

	"oppty-discovery-backend/internal/constant"
	"oppty-discovery-backend/internal/service"
	"oppty-discovery-backend/pkg/errors"
	"oppty-discovery-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 商机历史中心与归档控制器
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type HistoryHandler struct {
	opptyService service.OpptyService
}

func NewHistoryHandler(opptyService service.OpptyService) *HistoryHandler {
	return &HistoryHandler{opptyService: opptyService}
}

// ListHistory 分页获取用户商机历史
// @param c *gin.Context
func (h *HistoryHandler) ListHistory(c *gin.Context) {
	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	userID := c.GetUint64(constant.ContextKeyUserID)
	if userID == 0 {
		userID = 1
	}

	result, err := h.opptyService.ListUserHistory(c.Request.Context(), userID, pageNo, pageSize, keyword)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteHistory 删除单条商机历史
// @param c *gin.Context
func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "任务标识不能为空")
		return
	}

	userID := c.GetUint64(constant.ContextKeyUserID)
	if userID == 0 {
		userID = 1
	}

	if err := h.opptyService.DeleteHistory(c.Request.Context(), jobID, userID); err != nil {
		response.Fail(c, http.StatusInternalServerError, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "历史记录删除成功"})
}
