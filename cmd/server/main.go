package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"paper_ai/internal/api/handler"
	"paper_ai/internal/api/router"
	"paper_ai/internal/config"
	"paper_ai/internal/infrastructure/ai"
	"paper_ai/internal/service"
	"paper_ai/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting paper_ai service...")

	// 加载配置
	configPath := getConfigPath()
	if err := config.Load(configPath); err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	cfg := config.Get()
	logger.Info("config loaded successfully", zap.String("path", configPath))

	// 初始化AI提供商工厂
	factory := ai.GetFactory()
	if err := factory.InitProviders(cfg); err != nil {
		logger.Fatal("failed to init AI providers", zap.Error(err))
	}
	logger.Info("AI providers initialized", zap.Strings("providers", factory.ListProviders()))

	// 初始化服务层
	polishService := service.NewPolishService(factory)

	// 初始化处理器
	polishHandler := handler.NewPolishHandler(polishService)

	// 设置路由
	r := router.Setup(polishHandler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info("server started", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server startup failed", zap.Error(err))
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// 优雅关闭，最多等待30秒
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "./config/config.yaml"
}
