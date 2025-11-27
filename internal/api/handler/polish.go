package handler

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/service"
	"paper_ai/pkg/response"
)

// PolishHandler 润色处理器
type PolishHandler struct {
	polishService *service.PolishService
}

// NewPolishHandler 创建润色处理器
func NewPolishHandler(polishService *service.PolishService) *PolishHandler {
	return &PolishHandler{
		polishService: polishService,
	}
}

// Polish 处理段落润色请求
func (h *PolishHandler) Polish(c *gin.Context) {
	var req model.PolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// 从上下文获取用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		userID = int64(0) // 如果没有登录，使用0（后续可以改为返回错误）
	}

	resp, err := h.polishService.Polish(c.Request.Context(), &req, userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}
