package handler

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/service"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"
)

// PolishMultiVersionHandler 多版本润色处理器
type PolishMultiVersionHandler struct {
	multiVersionService *service.PolishMultiVersionService
}

// NewPolishMultiVersionHandler 创建多版本润色处理器
func NewPolishMultiVersionHandler(multiVersionService *service.PolishMultiVersionService) *PolishMultiVersionHandler {
	return &PolishMultiVersionHandler{
		multiVersionService: multiVersionService,
	}
}

// PolishMultiVersion 处理多版本润色请求
// @Summary 多版本润色
// @Description 生成多个版本的润色结果（conservative、balanced、aggressive）
// @Tags polish
// @Accept json
// @Produce json
// @Param request body model.PolishMultiVersionRequest true "润色请求"
// @Success 200 {object} response.Response{data=model.PolishMultiVersionResponse}
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/polish/multi-version [post]
func (h *PolishMultiVersionHandler) PolishMultiVersion(c *gin.Context) {
	var req model.PolishMultiVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("请先登录"))
		return
	}

	// 调用服务
	resp, err := h.multiVersionService.PolishMultiVersion(c.Request.Context(), &req, userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// SelectVersion 选择一个版本
// @Summary 选择多版本润色中的一个版本
// @Description 用户从多版本润色结果中选择一个版本，系统会将该版本的内容复制到主记录
// @Tags polish
// @Accept json
// @Produce json
// @Param trace_id path string true "润色记录的 trace_id"
// @Param version query string true "版本类型：conservative/balanced/aggressive"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/polish/select-version/{trace_id} [post]
func (h *PolishMultiVersionHandler) SelectVersion(c *gin.Context) {
	traceID := c.Param("trace_id")
	versionType := c.Query("version")

	if versionType == "" {
		response.Error(c, apperrors.NewInvalidParameterError("version 参数不能为空"))
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, apperrors.NewUnauthorizedError("请先登录"))
		return
	}

	// 调用服务
	err := h.multiVersionService.SelectVersion(c.Request.Context(), traceID, userID.(int64), versionType)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, map[string]string{
		"message": "版本选择成功",
		"trace_id": traceID,
		"selected_version": versionType,
	})
}
