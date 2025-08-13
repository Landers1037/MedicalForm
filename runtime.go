package main

import (
	"os"
	"path/filepath"
)

// 获取运行时的数据
// 单程序运行时可能存在进程名称不一致，使用固定的程序签名作为目录
func GetAppRuntimePath() string {
	// 根据Wails文档，如果Options中App的DataDirectory为空，
	// 将使用 %APPDATA%\[BinaryName.exe] 作为数据目录
	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		// 如果获取不到APPDATA环境变量，使用当前目录作为备选
		appDataDir = "."
	}

	// 获取可执行文件名称作为应用程序目录名
	appName := "medical-form"

	// 组合应用程序数据目录
	appDataPath := filepath.Join(appDataDir, appName)
	return appDataPath
}
