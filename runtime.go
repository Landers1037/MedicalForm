package main

import (
	"os"
	"path/filepath"
)

// 获取运行时的数据
func GetAppRuntimePath() string {
	// 根据Wails文档，如果Options中App的DataDirectory为空，
	// 将使用 %APPDATA%\[BinaryName.exe] 作为数据目录
	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		// 如果获取不到APPDATA环境变量，使用当前目录作为备选
		appDataDir = "."
	}

	// 获取可执行文件名称作为应用程序目录名
	execPath, err := os.Executable()
	var appName string
	if err != nil {
		// 如果获取可执行文件路径失败，使用默认名称
		appName = "medical-form"
	} else {
		// 提取可执行文件名（不含扩展名）
		appName = filepath.Base(execPath)
		if ext := filepath.Ext(appName); ext != "" {
			appName = appName[:len(appName)-len(ext)]
		}
	}

	// 组合应用程序数据目录
	appDataPath := filepath.Join(appDataDir, appName)
	return appDataPath
}
