package handler

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/service"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"
)

// ComparisonHandler 对比处理器
type ComparisonHandler struct {
	comparisonService *service.ComparisonService
}

// NewComparisonHandler 创建对比处理器
func NewComparisonHandler(comparisonService *service.ComparisonService) *ComparisonHandler {
	return &ComparisonHandler{
		comparisonService: comparisonService,
	}
}

// GetComparison 获取对比详情
// @Summary 获取润色对比详情
// @Description 根据 trace_id 获取原文和润色后文本的详细对比信息
// @Tags 对比
// @Accept json
// @Produce json
// @Param trace_id path string true "润色记录的 trace_id"
// @Success 200 {object} model.ComparisonResult
// @Failure 404 {object} response.ErrorResponse "记录不存在"
// @Failure 403 {object} response.ErrorResponse "无权访问"
// @Router /api/v1/polish/compare/{trace_id} [get]
func (h *ComparisonHandler) GetComparison(c *gin.Context) {
	traceID := c.Param("trace_id")

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("未登录"))
		return
	}

	result, err := h.comparisonService.GetComparison(c.Request.Context(), traceID, userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// ApplyAction 应用修改操作
// @Summary 接受或拒绝修改
// @Description 对单个修改执行接受或拒绝操作
// @Tags 对比
// @Accept json
// @Produce json
// @Param trace_id path string true "润色记录的 trace_id"
// @Param request body model.ChangeActionRequest true "操作请求"
// @Success 200 {object} model.ChangeActionResponse
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "记录不存在"
// @Failure 403 {object} response.ErrorResponse "无权访问"
// @Router /api/v1/polish/compare/{trace_id}/action [post]
func (h *ComparisonHandler) ApplyAction(c *gin.Context) {
	traceID := c.Param("trace_id")

	var req model.ChangeActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("未登录"))
		return
	}

	result, err := h.comparisonService.ApplyAction(c.Request.Context(), traceID, userID.(int64), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// BatchApplyAction 批量应用修改操作
// @Summary 批量接受或拒绝修改
// @Description 批量执行接受或拒绝操作
// @Tags 对比
// @Accept json
// @Produce json
// @Param trace_id path string true "润色记录的 trace_id"
// @Param request body model.BatchActionRequest true "批量操作请求"
// @Success 200 {object} model.BatchActionResponse
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "记录不存在"
// @Router /api/v1/polish/compare/{trace_id}/batch-action [post]
func (h *ComparisonHandler) BatchApplyAction(c *gin.Context) {
	traceID := c.Param("trace_id")

	var req model.BatchActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("未登录"))
		return
	}

	// 目前只支持 accept_all
	if req.Action != "accept_all" {
		response.Error(c, apperrors.NewInvalidParameterError("目前只支持 accept_all 操作"))
		return
	}

	result, err := h.comparisonService.BatchAcceptAll(c.Request.Context(), traceID, userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
