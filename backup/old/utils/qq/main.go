package main

import (
	"fmt"
	"github.com/lxn/win"
	"syscall"
)

var (
	hwndTarget      win.HWND
	originalWndProc uintptr
)

// 新的窗口过程
func newWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == win.WM_PAINT {
		fmt.Println("Target window is being painted.")
	}
	return win.CallWindowProc(originalWndProc, hwnd, msg, wParam, lParam)
}

// 替换窗口过程
func replaceWndProc() {
	originalWndProc = win.GetWindowLongPtr(hwndTarget, win.GWLP_WNDPROC)
	// 检查是否成功获取原始窗口过程
	if originalWndProc == 0 {
		fmt.Println("Failed to get original window procedure.")
		return
	}
	win.SetWindowLongPtr(hwndTarget, win.GWLP_WNDPROC, syscall.NewCallback(newWndProc))
}

// 取消替换
func restoreWndProc() {
	win.SetWindowLongPtr(hwndTarget, win.GWLP_WNDPROC, originalWndProc)
}

func main() {
	hwndTarget = 20388674

	// 替换窗口过程
	replaceWndProc()
	defer restoreWndProc()

	// 进入消息循环
	var msg win.MSG
	for {
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
