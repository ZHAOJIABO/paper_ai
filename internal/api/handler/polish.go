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

	resp, err := h.polishService.Polish(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}
