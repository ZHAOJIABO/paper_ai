package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "paper_ai/pkg/errors"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	traceID := getTraceID(c)
	c.JSON(http.StatusOK, Response{
		Code:    apperrors.CodeSuccess,
		Message: "success",
		Data:    data,
		TraceID: traceID,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	traceID := getTraceID(c)

	// 判断是否是自定义错误
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.HTTPStatus, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			TraceID: traceID,
		})
		return
	}

	// 未知错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    apperrors.CodeInternalError,
		Message: "internal server error",
		TraceID: traceID,
	})
}

// getTraceID 获取或生成TraceID
func getTraceID(c *gin.Context) string {
	traceID := c.GetString("trace_id")
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return traceID
}
