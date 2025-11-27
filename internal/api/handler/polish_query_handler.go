package handler

import (
	"strconv"
	"time"

	"paper_ai/internal/domain/repository"
	"paper_ai/internal/service"
	"paper_ai/pkg/response"

	"github.com/gin-gonic/gin"
)

// PolishQueryHandler 润色记录查询处理器
type PolishQueryHandler struct {
	polishService *service.PolishService
}

// NewPolishQueryHandler 创建查询处理器
func NewPolishQueryHandler(service *service.PolishService) *PolishQueryHandler {
	return &PolishQueryHandler{polishService: service}
}

// GetRecordByTraceID 根据TraceID查询记录
// GET /api/v1/polish/records/:trace_id
func (h *PolishQueryHandler) GetRecordByTraceID(c *gin.Context) {
	traceID := c.Param("trace_id")
	if traceID == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "trace_id is required",
		})
		return
	}

	// 从JWT上下文获取用户ID，确保只能查询自己的记录
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"code":    401,
			"message": "unauthorized",
		})
		return
	}

	record, err := h.polishService.GetRecordByTraceID(c.Request.Context(), traceID, userID.(int64))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// ListRecords 查询记录列表
// GET /api/v1/polish/records?page=1&page_size=20&provider=doubao&status=success&language=zh
func (h *PolishQueryHandler) ListRecords(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 解析过滤参数
	provider := c.Query("provider")
	status := c.Query("status")
	language := c.Query("language")
	style := c.Query("style")

	// 解析时间范围
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// 解析是否排除大文本
	excludeText := c.Query("exclude_text") == "true"

	// 构建查询选项
	builder := repository.NewQueryOptions().
		Page(page, pageSize).
		OrderBy("created_at", true) // 默认按创建时间降序

	// 从JWT上下文获取用户ID，确保只能查询自己的记录
	userID, exists := c.Get("user_id")
	if exists {
		builder.WithUserID(userID.(int64))
	}

	if provider != "" {
		builder.WithProvider(provider)
	}
	if status != "" {
		builder.WithStatus(status)
	}
	if language != "" {
		builder.WithLanguage(language)
	}
	if style != "" {
		builder.WithStyle(style)
	}

	if startTimeStr != "" && endTimeStr != "" {
		startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
		endTime, err2 := time.Parse(time.RFC3339, endTimeStr)
		if err1 == nil && err2 == nil {
			builder.WithTimeRange(startTime, endTime)
		}
	}

	if excludeText {
		builder.ExcludeText()
	}

	opts := builder.Build()

	// 查询记录
	records, total, err := h.polishService.ListRecords(c.Request.Context(), opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"records":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetStatistics 获取统计信息
// GET /api/v1/polish/statistics?start_time=2024-01-01T00:00:00Z&end_time=2024-12-31T23:59:59Z
func (h *PolishQueryHandler) GetStatistics(c *gin.Context) {
	// 解析时间范围
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	opts := repository.StatisticsOptions{}

	// 从JWT上下文获取用户ID，确保只能查看自己的统计信息
	userID, exists := c.Get("user_id")
	if exists {
		uid := userID.(int64)
		opts.UserID = &uid
	}

	if startTimeStr != "" && endTimeStr != "" {
		startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
		endTime, err2 := time.Parse(time.RFC3339, endTimeStr)
		if err1 == nil && err2 == nil {
			opts.TimeRange = &repository.TimeRange{
				Start: startTime,
				End:   endTime,
			}
		}
	}

	stats, err := h.polishService.GetStatistics(c.Request.Context(), opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stats)
}
