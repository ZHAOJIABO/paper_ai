package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "paper_ai/pkg/errors"
)

// Response 统一响应格式
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"` // 请求追踪ID（用于日志关联）
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	requestID := getRequestID(c)
	c.JSON(http.StatusOK, Response{
		Code:      apperrors.CodeSuccess,
		Message:   "success",
		Data:      data,
		RequestID: requestID,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	requestID := getRequestID(c)

	// 判断是否是自定义错误
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.HTTPStatus, Response{
			Code:      appErr.Code,
			Message:   appErr.Message,
			RequestID: requestID,
		})
		return
	}

	// 未知错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:      apperrors.CodeInternalError,
		Message:   "internal server error",
		RequestID: requestID,
	})
}

// getRequestID 获取或生成RequestID
func getRequestID(c *gin.Context) string {
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}
