package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"paper_ai/pkg/logger"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 生成TraceID
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		logger.Info("request completed",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
