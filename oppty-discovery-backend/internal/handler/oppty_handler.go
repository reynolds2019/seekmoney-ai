package handler

import (
	"net/http"

	"oppty-discovery-backend/internal/constant"
	"oppty-discovery-backend/internal/service"
	"oppty-discovery-backend/internal/types"
	"oppty-discovery-backend/pkg/errors"
	"oppty-discovery-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

/**
 * desc: 商机挖掘与智能分析接口控制器
 * date: 2026-10-13 11:25:00
 * @author reynolds
 * @since v1.0.0
 */
type OpptyHandler struct {
	opptyService service.OpptyService
}

func NewOpptyHandler(opptyService service.OpptyService) *OpptyHandler {
	return &OpptyHandler{opptyService: opptyService}
}

// CreateJob 创建商机分析任务
// @param c *gin.Context
func (h *OpptyHandler) CreateJob(c *gin.Context) {
	var req types.AnalyzeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "请求参数校验失败：关键词必填且视频数(5~20)、评论数(10~100)")
		return
	}

	userID := c.GetUint64(constant.ContextKeyUserID)
	if userID == 0 {
		userID = 1
	}

	resp, err := h.opptyService.CreateAnalyzeJob(c.Request.Context(), userID, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errors.CodeInternalError, err.Error())
		return
	}

	response.Accepted(c, "商机挖掘分析任务已创建，正在后台进行全网深度洞察", resp)
}

// GetJobDetail 查询任务进度与深度报告
// @param c *gin.Context
func (h *OpptyHandler) GetJobDetail(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "任务标识不能为空")
		return
	}

	userID := c.GetUint64(constant.ContextKeyUserID)
	jobDetail, err := h.opptyService.GetJobDetail(c.Request.Context(), jobID, userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errors.CodeInternalError, err.Error())
		return
	}

	if jobDetail == nil {
		response.Fail(c, http.StatusNotFound, errors.CodeJobNotFound, "未查询到该商机任务")
		return
	}

	response.Success(c, jobDetail)
}

// CalculateCost 实时计算预估成本
// @param c *gin.Context
func (h *OpptyHandler) CalculateCost(c *gin.Context) {
	var req struct {
		MaxVideos   int  `json:"maxVideos"`
		MaxComments int  `json:"maxComments"`
		DeepCrawl   bool `json:"deepCrawl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errors.CodeParamInvalid, "参数错误")
		return
	}

	cost := h.opptyService.CalculateEstimatedCost(req.MaxVideos, req.MaxComments, req.DeepCrawl)
	response.Success(c, gin.H{"estimatedCost": cost})
}
