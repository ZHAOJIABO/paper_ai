package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"paper_ai/internal/config"
	"paper_ai/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// 运行数据库迁移
	if cfg.AutoMigrate {
		if err := runMigrations(sqlDB); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		logger.Info("database migrations completed")
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

// runMigrations 运行数据库迁移
func runMigrations(sqlDB *sql.DB) error {
	// 获取项目根目录的 migrations 文件夹路径
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// 创建 postgres 驱动实例
	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// 创建迁移实例
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// 获取当前版本
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		logger.Info("no migrations applied yet")
	} else {
		logger.Info("current migration version",
			zap.Uint("version", version),
			zap.Bool("dirty", dirty),
		)
	}

	return nil
}

// GetGormDB 获取原生GORM DB实例（用于仓储实现）
func (db *DB) GetGormDB() *gorm.DB {
	return db.DB
}
