package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// MetaConfig 软件配置结构体
type MetaConfig struct {
	Title   string `json:"title"`   // 程序标题
	Version string `json:"version"` // 程序版本
}

// 全局变量
var (
	Title   string = "医疗管理表单系统" // 程序标题
	Version string = "1.0.0"    // 程序版本
)

// LoadMetaConfig 从meta.json文件中加载软件配置
func LoadMetaConfig() error {
	// 获取程序运行目录
	runDir := GetAppRuntimePath()
	metaPath := filepath.Join(runDir, "meta.json")

	// 检查meta.json是否存在
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		// 文件不存在，创建默认配置文件
		return createDefaultMetaConfig(metaPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return err
	}

	// 解析JSON
	var config MetaConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// 更新全局变量
	if config.Title != "" {
		Title = config.Title
	}
	if config.Version != "" {
		Version = config.Version
	}

	return nil
}

// createDefaultMetaConfig 创建默认的meta.json配置文件
func createDefaultMetaConfig(metaPath string) error {
	// 创建默认配置
	defaultConfig := MetaConfig{
		Title:   Title,
		Version: Version,
	}

	// 转换为JSON
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(metaPath, data, 0644)
}
