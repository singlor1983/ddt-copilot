package utils

import (
	"fmt"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

// 定义 Windows API 函数
var (
	procEnumChildWindows         = user32.NewProc("EnumChildWindows")
	procIsWindow                 = user32.NewProc("IsWindow")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

// GetProcessID 根据进程名获取进程ID
func GetProcessID(processName string) ([]uint32, error) {
	var processID []uint32

	// 创建进程快照
	processSnap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(processSnap)

	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))

	// 遍历进程
	for {
		if err = windows.Process32Next(processSnap, &processEntry); err != nil {
			break
		}

		exeFileName := windows.UTF16ToString(processEntry.ExeFile[:])
		if strings.EqualFold(exeFileName, processName) {
			processID = append(processID, processEntry.ProcessID)
		}
	}

	if len(processID) == 0 {
		return nil, fmt.Errorf("process:%s not found", processName)
	}
	return processID, nil
}

func GetWindowsByPID(pid uint32) ([]windows.Handle, error) {
	var wds []windows.Handle

	callback := func(hwnd windows.Handle, lParam uintptr) uintptr {
		var curPid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&curPid)))

		if curPid == uint32(lParam) {
			wds = append(wds, hwnd)
		}
		return 1 // 继续枚举
	}

	// 调用 EnumWindows
	_, _, _ = procEnumWindows.Call(syscall.NewCallback(callback), uintptr(pid))

	return wds, nil
}

func GetFirstWindowByPID(pid uint32) (windows.Handle, error) {
	wds, err := GetWindowsByPID(pid)
	if err != nil {
		return 0, err
	}
	if len(wds) == 0 {
		return 0, fmt.Errorf("not found windows. pid:%d", pid)
	}
	return wds[0], nil
}

func GetAllChildWindows(parentHwnd windows.Handle) ([]windows.Handle, error) {
	// 检查父窗口有效性
	isValid, _, _ := procIsWindow.Call(uintptr(parentHwnd))
	if isValid == 0 {
		return nil, fmt.Errorf("parent Window Handle %v is not valid", parentHwnd)
	}

	var childWindows []windows.Handle

	callback := func(hwnd windows.Handle, lParam uintptr) uintptr {
		childWindows = append(childWindows, hwnd)

		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		return 1 // 继续枚举
	}

	procEnumChildWindows.Call(uintptr(parentHwnd), syscall.NewCallback(callback), 0)

	return childWindows, nil
}
