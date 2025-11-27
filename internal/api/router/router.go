package router

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/api/handler"
	"paper_ai/internal/api/middleware"
)

// Setup 设置路由
func Setup(polishHandler *handler.PolishHandler, queryHandler *handler.PolishQueryHandler) *gin.Engine {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 段落润色
		v1.POST("/polish", polishHandler.Polish)

		// 查询记录（新增）
		v1.GET("/polish/records", queryHandler.ListRecords)
		v1.GET("/polish/records/:trace_id", queryHandler.GetRecordByTraceID)

		// 统计信息（新增）
		v1.GET("/polish/statistics", queryHandler.GetStatistics)
	}

	return r
}
