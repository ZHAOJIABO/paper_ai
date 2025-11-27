package database

import (
	"context"
	"fmt"
	"time"

	"paper_ai/internal/config"
	"paper_ai/internal/infrastructure/persistence"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DB 数据库管理器
type DB struct {
	*gorm.DB
}

var instance *DB

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	// 构建DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)

	// GORM日志配置
	var logLevel gormLogger.LogLevel
	switch cfg.LogMode {
	case "silent":
		logLevel = gormLogger.Silent
	case "error":
		logLevel = gormLogger.Error
	case "warn":
		logLevel = gormLogger.Warn
	case "info":
		logLevel = gormLogger.Info
	default:
		logLevel = gormLogger.Info
	}

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	instance = &DB{DB: db}

	logger.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	// 自动迁移
	if cfg.AutoMigrate {
		if err := migrate(db); err != nil {
			return fmt.Errorf("failed to auto migrate: %w", err)
		}
		logger.Info("database auto migration completed")
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *DB {
	return instance
}

// Close 关闭数据库连接
func Close() error {
	if instance == nil {
		return nil
	}

	sqlDB, err := instance.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	logger.Info("database connection closed")
	return nil
}

// Health 健康检查
func Health(ctx context.Context) error {
	if instance == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := instance.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// migrate 自动迁移表结构
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&persistence.PolishRecordPO{},
		&persistence.UserPO{},
		&persistence.RefreshTokenPO{},
	)
}

// GetGormDB 获取原生GORM DB实例（用于仓储实现）
func (db *DB) GetGormDB() *gorm.DB {
	return db.DB
}
