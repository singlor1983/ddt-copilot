package utils

import (
	"github.com/lxn/win"
	"time"
)

func LeftClick(hWnd win.HWND, x int, y int, duration time.Duration) {
	lParam := uintptr(y<<16 | x)
	win.PostMessage(hWnd, win.WM_LBUTTONDOWN, win.MK_LBUTTON, lParam)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_LBUTTONUP, 0, lParam)
}

func KeyBoard(hWnd win.HWND, key uintptr, duration time.Duration) {
	win.PostMessage(hWnd, win.WM_KEYDOWN, key, 0)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_KEYUP, key, 0)
}
