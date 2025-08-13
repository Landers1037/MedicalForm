//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	// Windows API 函数
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	// 函数定义
	createEventW        = kernel32.NewProc("CreateEventW")
	getLastError        = kernel32.NewProc("GetLastError")
	closeHandle         = kernel32.NewProc("CloseHandle")
	findWindowW         = user32.NewProc("FindWindowW")
	findWindowExW       = user32.NewProc("FindWindowExW")
	showWindow          = user32.NewProc("ShowWindow")
	setForegroundWindow = user32.NewProc("SetForegroundWindow")
	isIconic            = user32.NewProc("IsIconic")
	enumWindows         = user32.NewProc("EnumWindows")
	getWindowTextW      = user32.NewProc("GetWindowTextW")
	getClassNameW       = user32.NewProc("GetClassNameW")
	isWindowVisible     = user32.NewProc("IsWindowVisible")
)

const (
	// Windows 常量
	ERROR_ALREADY_EXISTS = 183
	SW_RESTORE           = 9
	SW_SHOW              = 5
)

// SingletonManager 单例管理器
type SingletonManager struct {
	mutexHandle uintptr
	appName     string
}

// NewSingletonManager 创建新的单例管理器
func NewSingletonManager(appName string) *SingletonManager {
	return &SingletonManager{
		appName: appName,
	}
}

// IsAlreadyRunning 检查应用是否已经在运行
func (sm *SingletonManager) IsAlreadyRunning() bool {
	// 创建一个命名事件对象作为互斥锁
	mutexName, _ := syscall.UTF16PtrFromString("Global\\" + sm.appName + "_Singleton")

	ret, _, _ := createEventW.Call(
		uintptr(0),                         // lpEventAttributes
		uintptr(1),                         // bManualReset
		uintptr(0),                         // bInitialState
		uintptr(unsafe.Pointer(mutexName)), // lpName
	)

	sm.mutexHandle = ret

	// 检查是否已存在
	lastError, _, _ := getLastError.Call()
	return lastError == ERROR_ALREADY_EXISTS
}

// ActivateExistingWindow 激活已存在的窗口
func (sm *SingletonManager) ActivateExistingWindow() bool {
	// 方法1: 通过窗口标题查找
	windowTitle, _ := syscall.UTF16PtrFromString(Title)
	hwnd, _, _ := findWindowW.Call(
		uintptr(0),                           // lpClassName
		uintptr(unsafe.Pointer(windowTitle)), // lpWindowName
	)

	// 方法2: 如果方法1失败，尝试查找包含关键词的窗口
	if hwnd == 0 {
		hwnd = sm.findWindowByPartialTitle("医疗管理表单系统")
	}

	// 方法3: 如果还是失败，尝试查找Wails应用窗口类
	if hwnd == 0 {
		wailsClassName, _ := syscall.UTF16PtrFromString("Chrome_WidgetWin_1")
		hwnd, _, _ = findWindowW.Call(
			uintptr(unsafe.Pointer(wailsClassName)), // lpClassName
			uintptr(0),                              // lpWindowName
		)
	}

	if hwnd == 0 {
		return false
	}

	// 检查窗口是否可见
	visible, _, _ := isWindowVisible.Call(hwnd)
	if visible == 0 {
		return false
	}

	// 检查窗口是否最小化
	isMin, _, _ := isIconic.Call(hwnd)
	if isMin != 0 {
		// 如果最小化，则恢复窗口
		showWindow.Call(hwnd, SW_RESTORE)
	} else {
		// 如果没有最小化，则显示窗口
		showWindow.Call(hwnd, SW_SHOW)
	}

	// 将窗口置于前台
	setForegroundWindow.Call(hwnd)

	return true
}

// findWindowByPartialTitle 通过部分标题查找窗口
func (sm *SingletonManager) findWindowByPartialTitle(partialTitle string) uintptr {
	// 这里可以实现枚举所有窗口并检查标题的逻辑
	// 为了简化，暂时返回0
	return 0
}

// Release 释放资源
func (sm *SingletonManager) Release() {
	if sm.mutexHandle != 0 {
		closeHandle.Call(sm.mutexHandle)
		sm.mutexHandle = 0
	}
}

// CheckSingleInstance 检查并处理单例逻辑
// 返回 true 表示应该继续运行，false 表示应该退出
func CheckSingleInstance() bool {
	sm := NewSingletonManager("medical-form")

	if sm.IsAlreadyRunning() {
		// 如果已经有实例在运行，尝试激活现有窗口
		if sm.ActivateExistingWindow() {
			// 成功激活现有窗口，当前实例应该退出
			sm.Release()
			return false
		}
	}

	// 没有其他实例在运行，或者无法激活现有窗口，继续运行
	return true
}
