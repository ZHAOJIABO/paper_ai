package router

import (
	"github.com/gin-gonic/gin"
	"paper_ai/internal/api/handler"
	"paper_ai/internal/api/middleware"
	"paper_ai/internal/infrastructure/security"
)

// Setup 设置路由
func Setup(
	polishHandler *handler.PolishHandler,
	multiVersionHandler *handler.PolishMultiVersionHandler,
	queryHandler *handler.PolishQueryHandler,
	comparisonHandler *handler.ComparisonHandler,
	authHandler *handler.AuthHandler,
	jwtManager *security.JWTManager,
) *gin.Engine {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 根路径 - 欢迎页面
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Paper AI API",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/v1",
		})
	})

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthRequired(jwtManager))
		{
			// 用户相关
			authenticated.GET("/auth/me", authHandler.GetCurrentUser)
			authenticated.POST("/auth/logout", authHandler.Logout)

			// 段落润色（需要认证）
			authenticated.POST("/polish", polishHandler.Polish)
			// 多版本润色（需要认证）
			authenticated.POST("/polish/multi", multiVersionHandler.PolishMultiVersion)

			// 查询记录（需要认证）
			authenticated.GET("/polish/records", queryHandler.ListRecords)
			authenticated.GET("/polish/records/:trace_id", queryHandler.GetRecordByTraceID)

			// 对比功能（需要认证）
			authenticated.GET("/polish/compare/:trace_id", comparisonHandler.GetComparison)
			authenticated.POST("/polish/compare/:trace_id/action", comparisonHandler.ApplyAction)
			authenticated.POST("/polish/compare/:trace_id/batch-action", comparisonHandler.BatchApplyAction)

			// 统计信息（需要认证）
			authenticated.GET("/polish/statistics", queryHandler.GetStatistics)
		}
	}

	return r
}
