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
	"paper_ai/internal/infrastructure/database"
	"paper_ai/internal/infrastructure/persistence"
	"paper_ai/internal/infrastructure/security"
	"paper_ai/internal/service"
	"paper_ai/pkg/idgen"
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

	// 初始化数据库（新增）
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("failed to init database", zap.Error(err))
	}
	defer database.Close()
	logger.Info("database initialized successfully")

	// 初始化ID生成器
	if err := idgen.Init(cfg.IDGen.WorkerID); err != nil {
		logger.Fatal("failed to init ID generator", zap.Error(err))
	}
	logger.Info("ID generator initialized", zap.Int64("worker_id", cfg.IDGen.WorkerID))

	// 初始化AI提供商工厂
	factory := ai.GetFactory()
	if err := factory.InitProviders(cfg); err != nil {
		logger.Fatal("failed to init AI providers", zap.Error(err))
	}
	logger.Info("AI providers initialized", zap.Strings("providers", factory.ListProviders()))

	// 创建仓储实现（新增）
	polishRepo := persistence.NewPolishRepository(database.GetDB().GetGormDB())
	userRepo := persistence.NewUserRepository(database.GetDB().GetGormDB())
	tokenRepo := persistence.NewRefreshTokenRepository(database.GetDB().GetGormDB())

	// 初始化JWT管理器
	jwtManager := security.NewJWTManager(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.AccessTokenExpiry)*time.Second,
		time.Duration(cfg.JWT.RefreshTokenExpiry)*time.Second,
	)
	logger.Info("JWT manager initialized")

	// 初始化服务层（注入仓储）
	polishService := service.NewPolishService(factory, polishRepo)
	comparisonService := service.NewComparisonService(polishRepo)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	// 初始化处理器
	polishHandler := handler.NewPolishHandler(polishService)
	queryHandler := handler.NewPolishQueryHandler(polishService)
	comparisonHandler := handler.NewComparisonHandler(comparisonService)
	authHandler := handler.NewAuthHandler(authService)

	// 设置路由（传入所有handler和jwtManager）
	r := router.Setup(polishHandler, queryHandler, comparisonHandler, authHandler, jwtManager)

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
