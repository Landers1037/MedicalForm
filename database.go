package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseService 数据库服务
type DatabaseService struct {
	db     *gorm.DB
	dbPath string
}

// NewDatabaseService 创建新的数据库服务实例
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// InitDB 初始化数据库连接
func (ds *DatabaseService) InitDB() error {
	// 获取应用程序数据目录

	// 创建应用程序专用目录
	appDir := GetAppRuntimePath()
	err := os.MkdirAll(appDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create app directory: %v", err)
	}

	// 数据库文件路径
	dbPath := filepath.Join(appDir, DBName)

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent, // 生产环境使用Silent
			Colorful:      false,
		},
	)

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	ds.db = db

	// 自动迁移数据库表
	err = ds.AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

// AutoMigrate 自动迁移所有表
func (ds *DatabaseService) AutoMigrate() error {
	return ds.db.AutoMigrate(&Patient{}, &Config{}, &Doctor{}, &Template{}, &Medicine{}, &ColumnConfig{})
}

// GetDB 获取数据库实例
func (ds *DatabaseService) GetDB() *gorm.DB {
	return ds.db
}

// Close 关闭数据库连接
func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// connectDB 连接指定路径的数据库
func (ds *DatabaseService) connectDB() error {
	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(ds.dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	ds.db = db
	return nil
}

// CloseDB 关闭数据库连接（别名方法）
func (ds *DatabaseService) CloseDB() error {
	return ds.Close()
}
