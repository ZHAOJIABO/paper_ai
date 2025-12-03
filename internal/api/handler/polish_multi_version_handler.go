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
